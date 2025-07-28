package metrics

type noopAuthMetrics struct{}

func NewNoopAuthMetrics() AuthMetrics {
	return &noopAuthMetrics{}
}

func (m *noopAuthMetrics) IncLogin(status string)                              {}
func (m *noopAuthMetrics) ObserveLoginDuration(status string, seconds float64) {}
