FROM golang:1.21-alpine
WORKDIR /app
# 安装git（go mod可能需要拉取依赖）
RUN apk add --no-cache git
# 复制代码
COPY . .
# 强制下载依赖，生成go.sum
RUN go mod tidy
# 编译（默认amd64，若需多架构可保留原编译逻辑）
RUN CGO_ENABLED=0 GOOS=linux go build -o sd-tv-live main.go
EXPOSE 9003
CMD ["./sd-tv-live"]
