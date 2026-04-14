#!/bin/bash
set -e

# 编译 BDInfoCLI 适配多架构
build_bdinfo() {
    local arch=$1
    local output_dir="./tools/bin/${arch}"
    mkdir -p ${output_dir}

    # 切换到 BDInfo 源码目录（需确认源码路径，此处假设为 ./bdinfo-src）
    cd ./bdinfo-src
    if [ "${arch}" = "amd64" ]; then
        GOARCH=amd64 GOOS=linux go build -o ${output_dir}/bdinfo ./cmd/bdinfo
    elif [ "${arch}" = "arm64" ]; then
        GOARCH=arm64 GOOS=linux go build -o ${output_dir}/bdinfo ./cmd/bdinfo
    fi
    cd -
}

# 编译 amd64 + arm64 版本
build_bdinfo amd64
build_bdinfo arm64

# 复制对应架构二进制到容器默认路径（Dockerfile 中根据架构选择）
mkdir -p /usr/local/bin
if [ "$(uname -m)" = "aarch64" ]; then
    cp ./tools/bin/arm64/bdinfo /usr/local/bin/bdinfo
else
    cp ./tools/bin/amd64/bdinfo /usr/local/bin/bdinfo
fi

# 赋权
chmod +x /usr/local/bin/bdinfo