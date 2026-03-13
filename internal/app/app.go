package app

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"time"

	"minfo/internal/config"
	"minfo/internal/httpapi"
	"minfo/internal/media"
)

func NewServer(staticFS fs.FS) (*http.Server, error) {
	port := config.Getenv("PORT", config.DefaultPort)

	preloadCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := media.LoadUDFModule(preloadCtx); err != nil {
		log.Printf("udf auto-load skipped: %v", err)
	}
	cancel()

	assets, err := fs.Sub(staticFS, "webui/dist")
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: httpapi.NewHandler(assets),
	}, nil
}
