package app

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"sso/internal/adapter/controller/grpc/server"
	"sso/internal/config"
	"sso/internal/service/auth"
	"sso/internal/storage/postgres"
)

// App todo
type App struct {
	gRPCServer *grpc.Server
	postgresDB *postgres.Storage
	log        *slog.Logger
	port       int
}

func New(log *slog.Logger, cfg *config.Config, db *postgres.Storage) *App {
	authService := auth.New(log, db, cfg.TokenTTL)
	gRPCServer := grpc.NewServer()
	server.RegisterGRPC(gRPCServer, authService)
	return &App{
		gRPCServer: gRPCServer,
		postgresDB: db,
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
	const op = "grpcapp.Run"
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
