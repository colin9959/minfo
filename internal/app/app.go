// Package app 负责组装 HTTP 服务、嵌入式资源和运行时配置。

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

// NewServer 会根据当前配置创建 HTTP Server，并在启动前预加载截图流程依赖的 loop / UDF 模块。
func NewServer(staticFS fs.FS) (*http.Server, error) {
	port := config.Getenv("PORT", config.DefaultPort)

	preloadCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := media.LoadLoopModule(preloadCtx); err != nil {
		log.Printf("loop auto-load skipped: %v", err)
	}
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
