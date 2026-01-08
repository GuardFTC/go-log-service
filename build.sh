#!/bin/bash

echo "开始编译 logging-mon-service..."

# 创建输出目录
mkdir -p ../../bin

# 设置项目名称
PROJECT_NAME="logging-mon-service"

# 编译Linux版本 (amd64)
echo "编译Linux版本..."
GOOS=linux GOARCH=amd64 go build -o "../../bin/${PROJECT_NAME}-linux-amd64" .
if [ $? -ne 0 ]; then
    echo "Linux编译失败!"
    exit 1
fi

# 编译Windows版本 (amd64)
echo "编译Windows版本..."
GOOS=windows GOARCH=amd64 go build -o "../../bin/${PROJECT_NAME}-windows.exe" .
if [ $? -ne 0 ]; then
    echo "Windows编译失败!"
    exit 1
fi

# 编译Linux ARM64版本
echo "编译Linux ARM64版本..."
GOOS=linux GOARCH=arm64 go build -o "../../bin/${PROJECT_NAME}-linux-arm64" .
if [ $? -ne 0 ]; then
    echo "Linux ARM64编译失败!"
    exit 1
fi

echo ""
echo "编译完成! 可执行文件已生成到 ../../bin/ 目录:"
echo "- ${PROJECT_NAME}-linux-amd64 (Linux x64)"
echo "- ${PROJECT_NAME}-linux-arm64 (Linux ARM64)"
echo "- ${PROJECT_NAME}-windows.exe (Windows x64)"