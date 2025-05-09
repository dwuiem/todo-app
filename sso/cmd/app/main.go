package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/storage/postgres"
	"syscall"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Setup logger
	log := setupLogger(cfg.Env)
	log.Info("Start application ...",
		slog.String("env", cfg.Env),
		slog.Any("config", cfg),
		slog.Int("port", cfg.GRPCServer.Port),
	)

	// Storage
	storage, err := postgres.New(cfg)
	if err != nil {
		slog.Error("Failed to connect to postgres", "error", err)
		os.Exit(1)
	}

	// Run application
	application := app.New(log, cfg, storage)
	go application.GRPCApp.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("Got signal. Stopping application ...", slog.String("signal", sign.String()))
	application.GRPCApp.Stop()
	log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case "local":
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return logger
}
