package appGRPC

import (
	server "crypto/internal/grpc"
	"crypto/internal/metrics"
	"crypto/internal/service"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *zap.Logger, service service.Service, port int) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.Interceptor()),
	)

	server.Register(gRPCServer, service)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.log.Info("start gRPC server at port " + fmt.Sprint(a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.log.Error("failed to listen", zap.Error(err))
		return fmt.Errorf("failed to listen: %w", err)
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		a.log.Error("failed to serve", zap.Error(err))
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping gRPC server...")
	a.gRPCServer.GracefulStop()
}
