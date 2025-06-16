package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"ginhub.com/Aller101/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService auth.Auth, port int) *GRPCApp {
	gRPCServer := grpc.NewServer()

	//TODO: еще будет сервисный слой auth - мб будет auth.Register... -> authgrpcRegister...
	auth.Register(gRPCServer, authService)

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

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *GRPCApp) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
	log.Info("Stopping server grpc")
}
