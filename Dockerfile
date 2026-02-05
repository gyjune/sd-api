# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
# 先复制go.mod和go.sum，利用缓存
COPY go.mod go.sum ./
RUN go mod download
# 再复制所有代码
COPY . .
# 编译amd64架构
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sd-tv-live main.go
# 编译arm64架构
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o sd-tv-live-arm64 main.go

# 运行阶段
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/sd-tv-live /app/sd-tv-live-amd64
COPY --from=builder /app/sd-tv-live-arm64 /app/sd-tv-live-arm64
RUN echo '#!/bin/sh\n\
ARCH=$(uname -m)\n\
if [ "$ARCH" = "aarch64" ]; then\n\
    exec /app/sd-tv-live-arm64\n\
else\n\
    exec /app/sd-tv-live-amd64\n\
fi' > /app/start.sh && chmod +x /app/start.sh
EXPOSE 9003
CMD ["/app/start.sh"]
