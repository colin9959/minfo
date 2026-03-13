## 项目介绍

`minfo` 是一个本地媒体信息检测 Web 工具，主要功能：
- 输出 MediaInfo 信息
- 输出 BDInfo 信息
- 使用 Seedbox 截图脚本生成 4 张截图压缩包
- 直接输出 Pixhost 图床链接

## 项目结构

- `cmd/minfo`：可执行程序入口
- `assets.go`：嵌入前端静态资源并向应用层提供文件系统
- `internal/app`：HTTP 服务组装与启动
- `internal/httpapi`：API 路由入口
- `internal/httpapi/handlers`：接口处理器
- `internal/httpapi/middleware`：认证与日志中间件
- `internal/httpapi/transport`：请求解析、响应输出、DTO
- `internal/media`：媒体路径解析、ISO 挂载、根目录与候选文件发现
- `internal/screenshot`：截图脚本调度、随机时间点生成、打包压缩
- `internal/system`：外部命令执行与进程组回收
- `webui`：前端界面

## Docker 运行

示例 `docker-compose.yml`：

```yaml
services:
  minfo:
    image: ghcr.io/mirrorb/minfo:latest
    container_name: minfo
    privileged: true
    ports:
      - "28081:8080"
    environment:
      PORT: "8080"
      WEB_PASSWORD: "adminadmin"
      REQUEST_TIMEOUT: "20m"
    volumes:
      - /lib/modules:/lib/modules:ro # 程序会自动尝试加载 `udf` 内核模块用于挂载ISO
      - /your/media/path1:/media_path1:ro
      - /your/media/path2:/media_path2:ro
      - /your/media/path3:/media_path3:ro
      - /your/media/path4:/media_path4:ro
    restart: unless-stopped
```

启动：

```bash
docker compose up -d
```

访问：
- `http://localhost:28081`
