package httpx

import (
	"encoding/json"
	"net/http"
)

type responseJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(payload)
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	writeJSON(w, status, responseJSON{
		Message: message,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{
		Message: message,
	})
}
