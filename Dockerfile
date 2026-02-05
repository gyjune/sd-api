# 构建阶段（多架构编译）
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder
WORKDIR /app
# 安装git并清理缓存
RUN apk add --no-cache git && rm -rf /var/cache/apk/*
# 复制代码
COPY . .
# 修正编译命令：GOOS拼写正确+引号格式规范
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o sd-tv-live main.go

# 运行阶段（精简镜像）
FROM scratch
COPY --from=builder /app/sd-tv-live /sd-tv-live
EXPOSE 9003
CMD ["/sd-tv-live"]