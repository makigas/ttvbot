package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/fx"
	"gopkg.makigas.es/ttvbot/domain/config"
)

var Module = fx.Module(
	"HttpServer",
	fx.Provide(func(lc fx.Lifecycle, cfg *config.Config) *HttpServer {
		srv := NewServer()
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				addr := fmt.Sprintf("%s:%d", cfg.Httpd.ServerBind, cfg.Httpd.ServerPort)
				ln, err := net.Listen("tcp", addr)
				if err != nil {
					return err
				}
				go srv.server.Serve(ln)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.server.Shutdown(ctx)
			},
		})
		return srv
	}),
)

type HttpServer struct {
	server *http.Server
	router *chi.Mux
}

func NewServer() *HttpServer {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "ttvbot is always listening")
	})
	server := &http.Server{Handler: router}
	return &HttpServer{router: router, server: server}
}

func (httpd *HttpServer) AddHandler(closure func(router *chi.Mux)) {
	closure(httpd.router)
}
