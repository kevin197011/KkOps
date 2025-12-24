# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（不使用 vendor）
RUN go mod download && go mod verify

# 复制源代码（vendor 目录已被 .dockerignore 排除）
COPY . .

# 构建应用（确保在 Linux 环境下编译，使用 -mod=mod 忽略 vendor）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -a -installsuffix cgo -ldflags="-w -s" -o /app/api ./cmd/api

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/api .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./api"]

