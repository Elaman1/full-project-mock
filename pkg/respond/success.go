package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

type SuccessResponse struct {
	Success string `json:"success"`
}

func WithSuccess(w http.ResponseWriter, code int, success string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(SuccessResponse{Success: success})
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}

func WithSuccessJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}
