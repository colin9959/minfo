package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func screenshotsHandler(w http.ResponseWriter, r *http.Request) {
	if !ensurePost(w, r) {
		return
	}
	if err := parseForm(w, r); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer cleanupMultipart(r)

	path, cleanup, err := inputPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer cleanup()

	mode := requestedScreenshotMode(r.FormValue("mode"))
	variant := requestedScreenshotVariant(r.FormValue("variant"))

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	tempDir, err := os.MkdirTemp("", "minfo-shots-*")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer os.RemoveAll(tempDir)

	if mode == screenshotModeLinks {
		output, err := runScreenshotUpload(ctx, path, tempDir, variant)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, infoResponse{OK: true, Output: output})
		return
	}

	files, err := runScreenshotScript(ctx, path, tempDir, variant)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	zipBytes, err := zipFiles(files)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"screenshots.zip\"")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(zipBytes); err != nil {
		log.Printf("write response: %v", err)
	}
}
