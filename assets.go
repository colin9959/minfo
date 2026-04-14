// Package minfo 提供嵌入式 WebUI 静态资源入口。

package minfo

import (
	"embed"
	"io/fs"
)

//go:embed webui/dist/*
var embeddedWebUI embed.FS

// EmbeddedWebUI 返回嵌入到二进制中的 WebUI 静态文件系统。
func EmbeddedWebUI() fs.FS {
	return embeddedWebUI
}
