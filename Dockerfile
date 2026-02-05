# 构建阶段（多阶段构建，最小化体积）
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

ARG TARGETARCH
ARG TARGETOS=linux
ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 复制依赖文件并下载依赖（利用Docker缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 启用 Go modules 的最小模式，减少不必要的依赖
ENV GO111MODULE=on

# 编译时去掉调试信息和符号表，使用静态编译
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -ldflags="-s -w -extldflags '-static'" \
    -buildvcs=false \
    -o sd-tv-live \
    .

# 最终运行阶段（使用 scratch 最小镜像）
FROM scratch

# 添加时区信息（可选）
# COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /app/sd-tv-live /sd-tv-live

# 复制 SSL 证书（如果需要 HTTPS 请求）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 9003

# 健康检查（可选）
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD ["/sd-tv-live", "--health-check"]

CMD ["/sd-tv-live"]