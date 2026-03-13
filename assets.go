package minfo

import (
	"embed"
	"io/fs"
)

//go:embed webui/dist/*
var embeddedWebUI embed.FS

func EmbeddedWebUI() fs.FS {
	return embeddedWebUI
}
