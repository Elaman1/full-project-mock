package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetricsCollector struct {
	requestsTotal       *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
	loginSuccessCounter prometheus.Counter
	loginFailureCounter *prometheus.CounterVec
}

func NewPrometheusMetricsCollector(reg prometheus.Registerer) *PrometheusMetricsCollector {
	collector := &PrometheusMetricsCollector{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		loginSuccessCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "http_login_success_total",
				Help: "Total number of successful logins",
			},
		),
		loginFailureCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_login_failure_total",
				Help: "Total number of failed logins",
			},
			[]string{"reason"},
		),
	}

	// Регистрация метрик в переданный реестр
	reg.MustRegister(
		collector.requestsTotal,
		collector.requestDuration,
		collector.loginSuccessCounter,
		collector.loginFailureCounter,
	)

	return collector
}

func (c *PrometheusMetricsCollector) IncHttpRequest(method, path, status string) {
	c.requestsTotal.WithLabelValues(method, path, status).Inc()
}

func (c *PrometheusMetricsCollector) ObserveRequestDuration(method, path string, durationSeconds float64) {
	c.requestDuration.WithLabelValues(method, path).Observe(durationSeconds)
}

func (c *PrometheusMetricsCollector) LoginSuccessCounter() {
	c.loginSuccessCounter.Inc()
}

func (c *PrometheusMetricsCollector) LoginFailureCounter(reason string) {
	c.loginFailureCounter.WithLabelValues(reason).Inc()
}
