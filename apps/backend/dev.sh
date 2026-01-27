#!/bin/bash
# 进入脚本所在目录
cd "$(dirname "$0")"

# 检查 air 是否安装
if ! command -v air &> /dev/null; then
    echo "air command not found. Installing..."
    go install github.com/air-verse/air@latest
    
    # 尝试将 GOPATH/bin 加入 PATH (如果只在当前 shell 有效)
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# 确认安装成功
if ! command -v air &> /dev/null; then
    echo "Error: Failed to install air. Please install it manually: go install github.com/air-verse/air@latest"
    exit 1
fi

echo "Starting Air for hot reload..."
# 执行 air，默认使用当前目录下的 .air.toml
air
