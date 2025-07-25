# 多阶段构建 - 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git

# 复制go mod文件
COPY server/go.mod server/go.sum ./server/
COPY api/go.mod api/go.sum ./api/

# 下载依赖
WORKDIR /app/server
RUN go mod download

WORKDIR /app/api
RUN go mod download

# 复制源代码
WORKDIR /app
COPY . .

# 构建gRPC服务器
WORKDIR /app/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# 构建API服务器
WORKDIR /app/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api .

# 运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata wget

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 appgroup && \
    adduser -D -s /bin/sh -u 1001 -G appgroup appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server/myapp ./
COPY --from=builder /app/api/api ./

# 复制配置文件（如果有的话）
COPY --from=builder /app/server/config ./config

# 更改文件所有者
RUN chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 7777 8889

# 创建启动脚本
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "Starting gRPC server..."' >> /app/start.sh && \
    echo './myapp &' >> /app/start.sh && \
    echo 'GRPC_PID=$!' >> /app/start.sh && \
    echo 'echo "Waiting for gRPC server to start..."' >> /app/start.sh && \
    echo 'sleep 5' >> /app/start.sh && \
    echo 'echo "Starting API server..."' >> /app/start.sh && \
    echo './api &' >> /app/start.sh && \
    echo 'API_PID=$!' >> /app/start.sh && \
    echo 'echo "Both servers started. gRPC PID: $GRPC_PID, API PID: $API_PID"' >> /app/start.sh && \
    echo 'wait $GRPC_PID $API_PID' >> /app/start.sh && \
    chmod +x /app/start.sh

# 启动命令
CMD ["/app/start.sh"]