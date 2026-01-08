# 多阶段构建 - 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用 - 静态链接
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o logging-mon-service .

# 运行阶段 - 使用scratch最小镜像
FROM scratch

# 从builder阶段复制时区信息
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 从builder阶段复制CA证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 复制编译好的二进制文件
COPY --from=builder /app/logging-mon-service /logging-mon-service

# 创建非root用户目录结构
COPY --from=builder /etc/passwd /etc/passwd

# 暴露端口 (默认39801，可通过环境变量SERVICE_PORT覆盖)
EXPOSE 39801

# 设置环境变量
ENV TZ=Asia/Shanghai

# 运行应用
ENTRYPOINT ["/logging-mon-service"]

# 默认参数 - 生产环境
CMD ["-env", "prod"]