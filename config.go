package main

import (
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultPort           = "8080"
	defaultRoot           = "/media"
	maxUploadBytes        = int64(8 << 30)
	maxMemoryBytes        = int64(32 << 20)
	maxSuggestions        = 200
	mountTimeout          = 30 * time.Second
	umountTimeout         = 30 * time.Second
	defaultRequestTimeout = 10 * time.Minute
)

var requestTimeout = durationFromEnv("REQUEST_TIMEOUT", defaultRequestTimeout)

func noop() {}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		log.Printf("invalid %s=%q; fallback to %s", key, value, fallback)
		return fallback
	}
	return duration
}
