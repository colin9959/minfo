package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, payload infoResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, infoResponse{OK: false, Error: msg})
}

func writePathJSON(w http.ResponseWriter, status int, payload pathResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writePathError(w http.ResponseWriter, status int, msg string) {
	writePathJSON(w, status, pathResponse{OK: false, Error: msg})
}

func bestErrorMessage(err error, stderr, stdout string) string {
	msg := strings.TrimSpace(stderr)
	if msg == "" {
		msg = err.Error()
	}
	if strings.TrimSpace(stdout) != "" {
		msg += "\n\n" + strings.TrimSpace(stdout)
	}
	return msg
}

func combineCommandOutput(stdout, stderr string) string {
	output := strings.TrimSpace(stdout)
	if strings.TrimSpace(stderr) != "" {
		if output != "" {
			output += "\n\n"
		}
		output += strings.TrimSpace(stderr)
	}
	return output
}
