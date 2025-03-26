package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitTracer инициализирует OpenTelemetry с экспортом данных в Jaeger.
func InitTracer(ctx context.Context, serviceName, jaegerURL string) (func(context.Context) error, error) {
	// Создаем экспортёр для Jaeger
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return nil, err
	}

	// Создаем провайдер трассировки
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Устанавливаем глобальный провайдер трассировки
	otel.SetTracerProvider(tp)

	// Возвращаем функцию для корректного завершения работы
	return tp.Shutdown, nil
}
