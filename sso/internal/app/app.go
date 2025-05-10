package app

import (
	"google.golang.org/grpc"
	"log/slog"
	"sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/service/auth"
	"sso/internal/storage/postgres"
)

// todo
type App struct {
	gRPCServer *grpc.Server
	PostgresDB *postgres.Storage
}

func New(log *slog.Logger, cfg *config.Config, db *postgres.Storage) *App {
	authService := auth.New(log, db, cfg.TokenTTL)
	grpcApp := grpcapp.New(log, authService, cfg.GRPCServer.Port)
	return &App{
		//GRPCApp:    grpcApp,
		PostgresDB: db,
	}
}
