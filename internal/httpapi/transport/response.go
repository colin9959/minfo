package transport

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, payload InfoResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, InfoResponse{OK: false, Error: msg})
}

func WritePathJSON(w http.ResponseWriter, status int, payload PathResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WritePathError(w http.ResponseWriter, status int, msg string) {
	WritePathJSON(w, status, PathResponse{OK: false, Error: msg})
}
