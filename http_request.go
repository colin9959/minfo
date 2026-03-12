package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ensurePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return false
	}
	return true
}

func parseForm(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	return r.ParseMultipartForm(maxMemoryBytes)
}

func cleanupMultipart(r *http.Request) {
	if r.MultipartForm != nil {
		_ = r.MultipartForm.RemoveAll()
	}
}

func inputPath(r *http.Request) (string, func(), error) {
	path := strings.TrimSpace(r.FormValue("path"))
	path = strings.Trim(path, "\"")
	if path != "" {
		path = filepath.Clean(path)
		if _, err := os.Stat(path); err != nil {
			return "", noop, fmt.Errorf("path not found: %v", err)
		}
		return path, noop, nil
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return "", noop, errors.New("missing file or path")
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	tempFile, err := os.CreateTemp("", "minfo-*"+ext)
	if err != nil {
		return "", noop, err
	}

	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
		return "", noop, err
	}
	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempFile.Name())
		return "", noop, err
	}

	cleanup := func() {
		_ = os.Remove(tempFile.Name())
	}
	return tempFile.Name(), cleanup, nil
}
