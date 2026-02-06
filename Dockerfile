# 阶段1：跨架构编译（ARM64/ARM32通用，含国内代理）
FROM golang:1.22-alpine AS builder
# 安装依赖（git拉取依赖、upx压缩二进制）
RUN apk add --no-cache git upx
# 设置工作目录
WORKDIR /app
# 复制依赖文件，利用Docker层缓存
COPY go.mod go.sum ./
# 配置国内Go代理（解决依赖拉取超时）
RUN go env -w GOPROXY=https://goproxy.cn,direct
# 拉取依赖并清理缓存
RUN go mod download && go clean -modcache
# 复制项目源码
COPY . .
# 静态编译Go程序（关闭CGO+剔除调试信息+UPX极致压缩）
RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo -ldflags="-s -w" -o live-tv . && \
    upx --best --lzma live-tv

# 阶段2：运行阶段（替换为Alpine，适配ARM且国内拉取稳定）
FROM alpine:3.19
# 安装时区包（解决时间不一致问题）
RUN apk add --no-cache tzdata
# 配置时区
ENV TZ=Asia/Shanghai
# 复制编译好的二进制文件（核心运行产物）
COPY --from=builder /app/live-tv /live-tv
# 暴露端口（与代码默认9003一致）
EXPOSE 9003
# 启动程序
ENTRYPOINT ["/live-tv"]