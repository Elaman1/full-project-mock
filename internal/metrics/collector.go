package metrics

type MetricsCollector interface {
	IncHttpRequest(method, path, status string)
	ObserveRequestDuration(method, path string, durationSeconds float64)
	LoginSuccessCounter()
	LoginFailureCounter(reason string)
}
