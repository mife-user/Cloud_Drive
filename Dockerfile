FROM golang:1.25-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o cloud_drive ./cmd/main.go

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/cloud_drive .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 创建必要的目录
RUN mkdir -p logs storage

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./cloud_drive"]