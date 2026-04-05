// Package transport 提供统一的 JSON 响应封装。

package transport

import (
	"encoding/json"
	"net/http"
)

// WriteJSON 将通用信息响应编码为 JSON 并写回指定状态码。
func WriteJSON(w http.ResponseWriter, status int, payload InfoResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError 将错误消息包装成通用 JSON 响应。
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, InfoResponse{OK: false, Error: msg})
}

// WritePathJSON 将路径联想响应编码为 JSON 并写回指定状态码。
func WritePathJSON(w http.ResponseWriter, status int, payload PathResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WritePathError 将路径联想错误包装成统一的 JSON 响应。
func WritePathError(w http.ResponseWriter, status int, msg string) {
	WritePathJSON(w, status, PathResponse{OK: false, Error: msg})
}
