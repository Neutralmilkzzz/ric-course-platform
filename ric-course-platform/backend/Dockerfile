# 使用 Go 1.23 构建阶段
FROM golang:1.23 AS build
WORKDIR /app

# 先拷贝 go.mod 和 go.sum，下载依赖（利用缓存）
COPY go.mod ./
RUN go mod download

# 再拷贝源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/ric-backend

# 使用精简镜像运行
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=build /out/ric-backend /ric-backend

ENV PORT=8080
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/ric-backend"]
