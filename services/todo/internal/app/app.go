package app

import (
	"context"
	"log/slog"
	"net/http"
	"todo/internal/config"
	"todo/internal/service"
	"todo/internal/storage"
	"todo/internal/storage/postgres"
	ssogrpc "todo/internal/transport/client/sso/grpc"
	"todo/internal/transport/http-server/handler"
)

type App struct {
	httpServer *http.Server
	log        *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) *App {
	client, err := ssogrpc.New(
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		panic(err)
	}

	storage, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}

	taskStorage := postgres.NewTaskStorage(storage)
	listStorage := postgres.NewListStorage(storage)

	taskService := service.NewTaskService(taskStorage)
	listService := service.NewListService(listStorage)
	handlers := handler.New(log, serv, client, cfg.AppSecret)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      handlers.InitRoutes(),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &App{
		httpServer: srv,
		log:        log,
	}
}

func (app *App) MustRun() error {
	err := app.run()
	return err
}

func (app *App) Stop(ctx context.Context) error {
	app.log.Info("Stopping HTTP server")
	return app.httpServer.Shutdown(ctx)
}

func (app *App) run() error {
	app.log.Info("Listening on ", app.httpServer.Addr)
	return app.httpServer.ListenAndServe()
}
