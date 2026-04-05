// Package transport 提供 HTTP 请求解析和输入路径处理。

package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"minfo/internal/config"
	"minfo/internal/media"
)

// EnsurePost 确认请求方法为 POST；不满足时会直接写回 405 响应。
func EnsurePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return false
	}
	return true
}

// ParseForm 会解析Form，并把原始输入转换成结构化结果。
func ParseForm(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadBytes)
	return r.ParseMultipartForm(config.MaxMemoryBytes)
}

// CleanupMultipart 释放 ParseMultipartForm 创建的临时文件。
func CleanupMultipart(r *http.Request) {
	if r.MultipartForm != nil {
		_ = r.MultipartForm.RemoveAll()
	}
}

// InputPath 从表单里的 path 或上传文件中解析输入路径，并返回对应的清理函数。
func InputPath(r *http.Request) (string, func(), error) {
	path := strings.TrimSpace(r.FormValue("path"))
	path = strings.Trim(path, "\"")
	if path != "" {
		ctx, cancel := context.WithTimeout(r.Context(), config.RequestTimeout)
		defer cancel()
		return media.ResolveInputPath(ctx, path)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return "", func() {}, errors.New("missing file or path")
	}
	defer file.Close()

	tempDir, err := os.MkdirTemp("", "minfo-upload-*")
	if err != nil {
		return "", func() {}, err
	}
	tempPath := filepath.Join(tempDir, uploadFileName(header.Filename))
	tempFile, err := os.Create(tempPath)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return "", func() {}, err
	}

	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close()
		_ = os.RemoveAll(tempDir)
		return "", func() {}, err
	}
	if err := tempFile.Close(); err != nil {
		_ = os.RemoveAll(tempDir)
		return "", func() {}, err
	}

	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}
	return tempFile.Name(), cleanup, nil
}

// uploadFileName 清理上传文件名，避免路径穿越并为无效名称提供稳定兜底值。
func uploadFileName(name string) string {
	cleaned := strings.TrimSpace(name)
	cleaned = strings.ReplaceAll(cleaned, "\\", "/")
	cleaned = pathpkg.Base(cleaned)
	if cleaned == "" || cleaned == "." || cleaned == "/" {
		return "upload.bin"
	}
	return cleaned
}
