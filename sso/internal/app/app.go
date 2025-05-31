package app

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	server "sso/internal/adapter/controller/grpc"
	"sso/internal/adapter/repository/postgres"
	"sso/internal/app/config"
	"sso/internal/domain/usecase"
)

type App struct {
	gRPCServer *grpc.Server
	log        *slog.Logger
	port       int
}

func New(log *slog.Logger, cfg *config.Config, repository *postgres.Auth) *App {
	uc := usecase.New(log, repository, cfg.TokenTTL)
	gRPCServer := grpc.NewServer()
	server.RegisterGRPC(gRPCServer, uc)
	return &App{
		gRPCServer: gRPCServer,
		log:        log,
		port:       cfg.GRPCServer.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"
	log := a.log.With(slog.String("op", op))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Running GRPC Server", slog.String("addr", lis.Addr().String()))
	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("Stopping GRPC Server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
