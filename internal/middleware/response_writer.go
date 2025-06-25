package middleware

import "net/http"

type StatusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *StatusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
