package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Interceptor возвращает interceptor для трассировки и метрик.
func Interceptor() grpc.UnaryServerInterceptor {
	tracer := otel.Tracer("grpc-server")
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Начало трассировки
		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		// Замер времени выполнения
		start := time.Now()

		// Вызов основного обработчика
		resp, err := handler(ctx, req)

		// Вычисление времени выполнения
		duration := time.Since(start)

		// Обработка метрик
		ProcessRequest(info.FullMethod, duration)

		// Обработка ошибок
		if err != nil {
			s, _ := status.FromError(err)
			span.SetAttributes(attribute.String("error", s.Message()))
			span.RecordError(err)
		}

		return resp, err
	}
}
