package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"net"
	"net/http"
	"os"
)

type HelloWorldHandler struct {
}

func NewHelloWorldHandler() *HelloWorldHandler {
	return &HelloWorldHandler{}
}

func (*HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Printf("Body: %v", body)
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
	}
}

func NewServeMux(handler *HelloWorldHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hello", handler)

	return mux
}
func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("starting http server")
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
func main() {

	//setup containers
	fx.New(
		fx.Provide(NewHttpServer),
		fx.Provide(NewHelloWorldHandler),
		fx.Provide(NewServeMux),
		fx.Invoke(func(srv *http.Server) {}),
	).Run()

}
