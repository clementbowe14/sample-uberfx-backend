package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type HelloWorldHandler struct {
	log *zap.SugaredLogger
}

type BankAccountHandler struct {
	log *zap.SugaredLogger
}

type Route interface {
	http.Handler
	Pattern() string
}

func NewBankAccountHandler(log *zap.SugaredLogger) *BankAccountHandler {
	return &BankAccountHandler{log: log}
}

func (h *BankAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorw("failed to read request body", "error", err)
	}

	fmt.Fprintf(w, "This is a reminder that you're still broke! %s", string(body))

}

func (h *BankAccountHandler) Pattern() string {
	return "/bank"
}

func NewHelloWorldHandler(log *zap.SugaredLogger) *HelloWorldHandler {
	return &HelloWorldHandler{log: log}
}

func (*HelloWorldHandler) Pattern() string {
	return "/hello"
}

func (h *HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	h.log.Debugf("Someone has came to say hello to us... fuck %v", body)
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		h.log.Errorf("Thank god we failed our response... i mean whoops Error writing response: %v %v", err, os.Stderr)
	}
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}
func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.SugaredLogger) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				log.Errorf("Failed to start http server: %v", zap.Error(err))
				return err
			}
			log.Debugf("starting http server %v", ln.Addr())

			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Errorf("failed to initialize zap logger: %v", err)
	}

	sugar := logger.Sugar()

	//setup containers
	fx.New(
		fx.Supply(sugar),
		fx.Provide(
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			NewHttpServer,
		),
		fx.Provide(
			AsRoute(NewHelloWorldHandler),
			AsRoute(NewBankAccountHandler),
		),
		fx.Invoke(func(srv *http.Server) {}),
	).Run()

}
