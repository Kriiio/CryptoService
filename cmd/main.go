package main

import (
	"context"
	"crypto/config"
	"crypto/internal/app"
	zaplogger "crypto/internal/logger"
	"crypto/internal/metrics"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func main() {
	zaplogger.BuildLogger("INFO")
	logger = zaplogger.Logger().Named("main")
	cfg := config.MustLoad(logger)

	logger.Info("starting app")

	// Инициализация OpenTelemetry
	shutdown, err := metrics.InitTracer(context.Background(), "your-service", "http://jaeger:14268/api/traces")
	if err != nil {
		logger.Error("failed to initialize OpenTelemetry", zap.Error(err))
	}

	application := app.New(logger, cfg.GRPC.Port, cfg.DB.DSN)

	if err := application.RunMetricsServer(); err != nil {
		logger.Error("failed to run metrics server", zap.Error(err))
	}

	go application.GrpcServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	shutdown(context.TODO())

	application.GrpcServer.Stop()
}
