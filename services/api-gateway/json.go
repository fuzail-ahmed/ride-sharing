package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)

func writeJSON(w http.ResponseWriter, statusCode int, data contracts.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
