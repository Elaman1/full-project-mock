package respond

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WithError(w http.ResponseWriter, code int, msg string, lgr *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
	lgr.Error(msg)
}

func WithErrorJSON(w http.ResponseWriter, code int, payload interface{}, lgr *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error encoding error response %v", err)
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling error response %v", err)
	}

	lgr.Error(string(data))
}
