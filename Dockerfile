# 构建阶段（多架构编译）
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git && rm -rf /var/cache/apk/*
COPY . .
# 静态编译+剥离调试信息
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o sd-tv-live main.go

# 运行阶段（使用scratch空镜像，语法修正）
FROM scratch
COPY --from=builder /app/sd-tv-live /sd-tv-live
EXPOSE 9003
CMD ["/sd-tv-live"]
