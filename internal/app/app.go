package app

import (
	appGRPC "crypto/internal/app/grpc"
	"crypto/internal/metrics"
	"crypto/internal/service"
	"crypto/internal/storage"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type App struct {
	GrpcServer *appGRPC.App
	log        *zap.Logger
}

func New(log *zap.Logger, port int, dbcfg string) *App {

	database, err := storage.New(dbcfg, log)
	if err != nil {
		log.Error("failed to create database", zap.Error(err))
	}

	service := service.New(log, database)

	grpcServer := appGRPC.New(log, service, port)

	return &App{
		GrpcServer: grpcServer,
		log:        log,
	}
}

func (a *App) RunMetricsServer() error {
	a.log.Info("starting metrics server at port 9090")

	// Создаем HTTP-сервер для метрик
	http.Handle("/metrics", metrics.MetricsHandler())

	// Запускаем сервер на указанном порту
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", 9090), nil); err != nil {
			a.log.Error("failed to start metrics server", zap.Error(err))
		}
	}()

	return nil
}
