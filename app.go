package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed webui/dist/*
var staticFS embed.FS

func newServer() (*http.Server, error) {
	port := getenv("PORT", defaultPort)

	preloadCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := loadUDFModule(preloadCtx); err != nil {
		log.Printf("udf auto-load skipped: %v", err)
	}
	cancel()

	assets, err := fs.Sub(staticFS, "webui/dist")
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: logging(authenticate(newMux(assets))),
	}, nil
}

func newMux(assets fs.FS) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("/api/mediainfo", mediainfoHandler("MEDIAINFO_BIN", "mediainfo"))
	mux.HandleFunc("/api/bdinfo", bdinfoHandler("BDINFO_BIN", "bdinfo"))
	mux.HandleFunc("/api/screenshots", screenshotsHandler)
	mux.HandleFunc("/api/path", pathSuggestHandler)
	return mux
}
