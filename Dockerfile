# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 交叉编译（适配amd64/arm64）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sd-tv-live main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o sd-tv-live-arm64 main.go

# 运行阶段（轻量Alpine）
FROM alpine:3.19
WORKDIR /app
# 复制对应架构的可执行文件
COPY --from=builder /app/sd-tv-live /app/sd-tv-live-amd64
COPY --from=builder /app/sd-tv-live-arm64 /app/sd-tv-live-arm64
# 自动适配架构
RUN echo '#!/bin/sh\n\
ARCH=$(uname -m)\n\
if [ "$ARCH" = "aarch64" ]; then\n\
    exec /app/sd-tv-live-arm64\n\
else\n\
    exec /app/sd-tv-live-amd64\n\
fi' > /app/start.sh && chmod +x /app/start.sh
# 暴露端口
EXPOSE 9003
# 启动服务
CMD ["/app/start.sh"]