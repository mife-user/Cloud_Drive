# 使用轻量级基础镜像
FROM alpine:latest

WORKDIR /app

# 复制预先构建好的 Linux 可执行文件
COPY cloud_drive .

# 复制配置文件
COPY configs ./configs

# 创建必要的目录
RUN mkdir -p logs storage header

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./cloud_drive"]
