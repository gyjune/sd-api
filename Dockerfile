# 构建阶段（仅用于编译，最终丢弃）
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder
WORKDIR /app
# 仅安装必要依赖（git用于拉取go mod依赖，编译后删除）
RUN apk add --no-cache git && rm -rf /var/cache/apk/*
COPY . .
# 编译优化：静态编译+剥离调试信息（大幅减小二进制体积）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o sd-tv-live main.go

# 运行阶段（使用最小基础镜像，仅3.3兆）
FROM scratch  # 空镜像，仅包含二进制文件和必要配置
# 复制编译好的二进制文件（仅几兆）
COPY --from=builder /app/sd-tv-live /sd-tv-live
# 暴露端口（scratch镜像无需额外依赖）
EXPOSE 9003
# 启动命令
CMD ["/sd-tv-live"]
