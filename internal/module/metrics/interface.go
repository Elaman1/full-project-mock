package metrics

var (
	IncLoginSuccess = "success"
	IncLoginFail    = "fail"
)

type AuthMetrics interface {
	IncLogin(status string) // IncLoginSuccess | IncLoginFail
	ObserveLoginDuration(status string, seconds float64)
}
