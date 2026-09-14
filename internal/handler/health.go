package handler

import (
	"encoding/json"
	"net/http"
)

// HTTP Handler for GET /health
func Health(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(response)
}
