package server

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Route interface {
	http.Handler
	Pattern() string
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
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

			go func() {
				err := srv.Serve(ln)
				if err != nil {
					log.Errorf("Failed to start http server: %v", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
