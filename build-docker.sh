#!/bin/bash

# Docker镜像构建脚本

# 设置变量
IMAGE_NAME="logging-mon-service"
VERSION=${1:-"latest"}
REGISTRY=${2:-""}

# 如果指定了镜像仓库，添加前缀
if [ -n "$REGISTRY" ]; then
    FULL_IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}:${VERSION}"
else
    FULL_IMAGE_NAME="${IMAGE_NAME}:${VERSION}"
fi

echo "开始构建Docker镜像..."
echo "镜像名称: \"$FULL_IMAGE_NAME\""
echo ""

# 构建镜像
docker build -t "$FULL_IMAGE_NAME" .

if [ $? -eq 0 ]; then
    echo ""
    echo "Docker镜像构建成功!"
    echo "镜像名称: \"$FULL_IMAGE_NAME\""
    echo ""
    echo "使用方法:"
    echo "1. 运行容器 (开发环境):"
    echo "   docker run -p 39801:39801 \"$FULL_IMAGE_NAME\" -env dev"
    echo ""
    echo "2. 运行容器 (生产环境):"
    echo "   docker run -p 39801:39801 \"$FULL_IMAGE_NAME\""
    echo ""
    echo "3. 运行容器 (自定义端口):"
    echo "   docker run -p 8080:8080 -e SERVICE_PORT=8080 \"$FULL_IMAGE_NAME\""
    echo ""
    echo "4. 推送到镜像仓库:"
    echo "   docker push \"$FULL_IMAGE_NAME\""
    echo ""
    
    # 显示镜像信息
    echo "镜像信息:"
    docker images | grep "$IMAGE_NAME" | head -1
else
    echo ""
    echo "Docker镜像构建失败!"
    exit 1
fi