package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	processingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "processing_time",
			Help: "Время обработки запроса",
		},
		[]string{"method"},
	)

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Количество запросов",
		},
		[]string{"method"},
	)

	apiTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "api_processing_seconds",
		Help: "Время обращения во внешний API",
	}, []string{"method"})

	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "response_time",
			Help: "Время ответа",
		},
		[]string{"method"},
	)
)

// init регистрирует метрики Prometheus.
func init() {
	prometheus.MustRegister(processingTime)
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(apiTime)
	prometheus.MustRegister(responseTime)
}

func ProcessRequest(endpoint string, duration time.Duration) {
	requestCount.WithLabelValues(endpoint).Inc()
	processingTime.WithLabelValues(endpoint).Observe(duration.Seconds())
	apiTime.WithLabelValues(endpoint).Observe(duration.Seconds())
	responseTime.WithLabelValues(endpoint).Observe(duration.Seconds())
}

// MetricsHandler возвращает HTTP обработчик для Prometheus.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
