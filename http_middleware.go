package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func authenticate(next http.Handler) http.Handler {
	password := getenv("WEB_PASSWORD", "")
	if password == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pass, ok := parseBasicAuth(r.Header.Get("Authorization"))
		if !ok || pass != password {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"minfo\"")
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parseBasicAuth(header string) (string, string, bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(header, prefix) {
		return "", "", false
	}

	encoded := strings.TrimSpace(header[len(prefix):])
	if encoded == "" {
		return "", "", false
	}

	decoded, err := decodeBase64(encoded)
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(decoded, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func decodeBase64(value string) (string, error) {
	data, err := base64Decode(value)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(data) {
		return "", errors.New("invalid encoding")
	}
	return string(data), nil
}

func base64Decode(value string) ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(value)))
	n, err := base64.StdEncoding.Decode(buf, []byte(value))
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
