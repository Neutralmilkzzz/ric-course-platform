# Step 1: Build React 前端
FROM node:18 AS frontend-builder
WORKDIR /app/frontend
COPY frontend/ .
RUN npm install && npm run build

# Step 2: Build Go 后端
FROM golang:1.21 AS backend-builder
WORKDIR /app/backend
COPY backend/ .
# 把 React build 结果复制到 Go 的 static 文件夹
COPY --from=frontend-builder /app/frontend/build ./static
RUN go mod download
RUN go build -o server main.go

# Step 3: 运行容器
FROM alpine:latest
WORKDIR /root/
COPY --from=backend-builder /app/backend/server .
CMD ["./server"]
