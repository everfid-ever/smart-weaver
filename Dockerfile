# 构建阶段
FROM golang:1.23.0-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o smart-weaver-app ./cmd/app

# 运行阶段
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/smart-weaver-app .
COPY --from=builder /app/configs ./configs

# 暴露端口
EXPOSE 8091

# 运行应用
CMD ["./smart-weaver-app"]