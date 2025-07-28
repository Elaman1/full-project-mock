package metrics

import "github.com/prometheus/client_golang/prometheus"

type promAuthMetrics struct {
	loginCounter   *prometheus.CounterVec
	loginHistogram *prometheus.HistogramVec
}

func NewPromAuthMetrics(reg prometheus.Registerer) AuthMetrics {
	m := &promAuthMetrics{
		loginCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "auth_login_total",
			Help: "Total number of login attempts",
		}, []string{"status"}),
		loginHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "auth_login_duration_seconds",
			Help:    "Duration of login attempts",
			Buckets: prometheus.DefBuckets,
		}, []string{"status"}),
	}

	reg.MustRegister(m.loginCounter, m.loginHistogram)
	return m
}

func (m *promAuthMetrics) IncLogin(status string) {
	m.loginCounter.WithLabelValues(status).Inc()
}

func (m *promAuthMetrics) ObserveLoginDuration(status string, seconds float64) {
	m.loginHistogram.WithLabelValues(status).Observe(seconds)
}
