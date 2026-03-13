package transport

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"minfo/internal/config"
)

func EnsurePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return false
	}
	return true
}

func ParseForm(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadBytes)
	return r.ParseMultipartForm(config.MaxMemoryBytes)
}

func CleanupMultipart(r *http.Request) {
	if r.MultipartForm != nil {
		_ = r.MultipartForm.RemoveAll()
	}
}

func InputPath(r *http.Request) (string, func(), error) {
	path := strings.TrimSpace(r.FormValue("path"))
	path = strings.Trim(path, "\"")
	if path != "" {
		path = filepath.Clean(path)
		if _, err := os.Stat(path); err != nil {
			return "", func() {}, fmt.Errorf("path not found: %v", err)
		}
		return path, func() {}, nil
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return "", func() {}, errors.New("missing file or path")
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	tempFile, err := os.CreateTemp("", "minfo-*"+ext)
	if err != nil {
		return "", func() {}, err
	}

	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
		return "", func() {}, err
	}
	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempFile.Name())
		return "", func() {}, err
	}

	cleanup := func() {
		_ = os.Remove(tempFile.Name())
	}
	return tempFile.Name(), cleanup, nil
}
