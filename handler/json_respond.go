package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, status int, msg string) {
	if status > 499 {
		log.Printf("responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, status, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshaling JSON: %s", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
