package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"sso/internal/adapter/controller/grpc/auth"
)

type GRPCApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService auth.Service, port int) *GRPCApp {
	gRPCServer := grpc.NewServer()
	auth.RegisterGRPC(gRPCServer, authService)
	return &GRPCApp{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *GRPCApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *GRPCApp) Run() error {
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

func (a *GRPCApp) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("Stopping GRPC Server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
