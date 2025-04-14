package main

import (
	"log"
	"log/slog"
	"net/http"
	"todo/internal/config"
	"todo/internal/repository"
	"todo/internal/repository/postgres"
	"todo/internal/service"
	ssogrpc "todo/internal/transport/client/sso/grpc"
	"todo/internal/transport/http-server/handler"
)

func main() {
	cfg := config.MustLoad()
	db := postgres.New(*cfg)
	repos := repository.New(db)
	services := service.New(repos)
	handlers := handler.New(services)

	_, err := ssogrpc.New(
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      handlers.InitRoutes(),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server", err.Error())
	}
	log.Fatal("Stopping server", slog.String("address", cfg.HTTPServer.Addr))
}
