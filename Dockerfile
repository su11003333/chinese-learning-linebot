# 使用官方 Go 映像作為建構階段
FROM golang:1.21-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 安裝必要的套件
RUN apk add --no-cache git ca-certificates tzdata

# 複製 go mod 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 建構應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o linebot main.go

# 使用輕量級的 Alpine 映像作為運行階段
FROM alpine:latest

# 安裝 ca-certificates 用於 HTTPS 請求
RUN apk --no-cache add ca-certificates tzdata

# 設定工作目錄
WORKDIR /root/

# 從建構階段複製執行檔
COPY --from=builder /app/linebot .

# 設定時區
ENV TZ=Asia/Taipei

# 暴露端口 (Cloud Run 使用 PORT 環境變數)
EXPOSE 8080

# 執行應用程式
CMD ["./linebot"]