package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/internal/app"
	"todo/internal/config"
)

func main() {
	cfg := config.MustLoad()

	// Setup logger
	log := setupLogger(cfg.Env)
	log.Info("Start application ...",
		slog.String("env", cfg.Env),
		slog.Any("config", cfg),
	)

	application := app.New(log, cfg)

	serverError := make(chan error, 1)

	go func() {
		serverError <- application.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverError:
		log.Info("Application exited with error", slog.String("error", err.Error()))
	case sig := <-stop:
		log.Info("Got signal. Stopping application...", slog.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := application.Stop(ctx); err != nil {
			log.Info("Application exited with error", slog.String("error", err.Error()))
		} else {
			log.Info("Application exited gracefully", slog.String("error", err.Error()))
		}
	}
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
