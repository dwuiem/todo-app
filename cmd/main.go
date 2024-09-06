package main

import (
	"log"
	"log/slog"
	"net/http"
	"todo-app/internal/config"
	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/repository/postgres"
	"todo-app/internal/service"
)

func main() {
	cfg := config.MustLoad()
	db := postgres.New(*cfg)
	repos := repository.New(db)
	services := service.New(repos)
	handlers := handler.New(services)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      handlers.InitRoutes(),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server", err.Error())
	}
	log.Fatal("Stopping server", slog.String("address", cfg.HTTPServer.Addr))
}
