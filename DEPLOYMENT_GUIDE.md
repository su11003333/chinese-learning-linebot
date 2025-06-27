# Google Cloud Run 部署指南

本指南詳細說明如何將 Chinese Learning LINE Bot 部署到 Google Cloud Run。

## 前置作業

### 1. 確認 Google Cloud 專案設定
```bash
# 設定專案 ID
gcloud config set project kid-mcp

# 檢查設定是否正確
gcloud config list
```

### 2. 啟用必要的 Google Cloud 服務
```bash
# 啟用 Cloud Run 和 Cloud Build 服務
gcloud services enable run.googleapis.com cloudbuild.googleapis.com
```

## 初次部署流程

### 1. 檢查 Dockerfile 版本
確保 Dockerfile 中的 Go 版本與 go.mod 中要求的版本一致：

```dockerfile
# 檢查 go.mod 中的 Go 版本，例如：go 1.23.4
# Dockerfile 應該使用相應版本
FROM golang:1.23-alpine AS builder
```

### 2. 執行部署指令
```bash
gcloud run deploy chinese-linebot \
  --source . \
  --region asia-east1 \
  --platform managed \
  --allow-unauthenticated
```

**這個指令會自動執行以下步驟：**

#### A. 上傳程式碼
- 將當前目錄的所有檔案（除了 .dockerignore 中指定的）打包上傳到 Google Cloud Build

#### B. 建置 Docker 映像檔
Cloud Build 會按照 Dockerfile 的指令執行：
1. **建置階段** (`FROM golang:1.23-alpine AS builder`)
   - 下載 Go 1.23 開發環境
   - 設定工作目錄為 `/app`
   - 安裝必要套件：git, ca-certificates, tzdata
   - 複製 `go.mod` 和 `go.sum` 檔案
   - 下載 Go 依賴包 (`go mod download`)
   - 複製所有程式碼
   - 編譯程式：`go build -o linebot main.go`

2. **運行階段** (`FROM alpine:latest`)
   - 使用輕量級 Alpine Linux 環境
   - 安裝 ca-certificates 和 tzdata
   - 從建置階段複製編譯好的執行檔
   - 設定時區為 Asia/Taipei
   - 開放 8080 port
   - 設定啟動指令

#### C. 推送映像檔
將建置完成的 Docker 映像檔推送到 Google Container Registry：
```
asia-east1-docker.pkg.dev/kid-mcp/cloud-run-source-deploy/chinese-linebot:latest
```

#### D. 部署服務
- 建立 Cloud Run 服務：`chinese-linebot`
- 部署地區：`asia-east1` (台灣)
- 允許公開存取：`--allow-unauthenticated`
- 自動分配服務 URL

### 3. 設定環境變數
```bash
gcloud run services update chinese-linebot \
  --set-env-vars "LINE_CHANNEL_SECRET=your_channel_secret,LINE_CHANNEL_ACCESS_TOKEN=your_access_token,FIREBASE_PROJECT_ID=kid-mcp" \
  --region asia-east1
```

**重要環境變數說明：**
- `LINE_CHANNEL_SECRET`：LINE Bot 的 Channel Secret
- `LINE_CHANNEL_ACCESS_TOKEN`：LINE Bot 的 Channel Access Token
- `FIREBASE_PROJECT_ID`：Firebase 專案 ID
- `PORT`：系統自動設定，不需要手動指定

## 部署完成後的設定

### 1. 服務資訊
- **服務 URL**：`https://chinese-linebot-x5h7hw374a-de.a.run.app`
- **Webhook URL**：`https://chinese-linebot-x5h7hw374a-de.a.run.app/webhook`

### 2. LINE Bot Webhook 設定
1. 前往 [LINE Developers Console](https://developers.line.biz/)
2. 選擇你的 Bot 專案
3. 進入 "Messaging API" 設定頁面
4. 設定 **Webhook URL**：
   ```
   https://chinese-linebot-x5h7hw374a-de.a.run.app/webhook
   ```
5. 啟用 **Use webhook** 選項
6. 點擊 **Verify** 按鈕測試連線

## 程式碼修改後的重新部署

### 簡單重新部署
```bash
# 修改程式碼後，執行重新部署
gcloud run deploy chinese-linebot \
  --source . \
  --region asia-east1 \
  --platform managed \
  --allow-unauthenticated
```

### 完整開發流程（建議）
```bash
# 1. 修改程式碼
# 2. 本地測試（可選）
go run main.go

# 3. 版本控制
git add .
git commit -m "描述你的修改內容"
git push origin main

# 4. 重新部署
gcloud run deploy chinese-linebot \
  --source . \
  --region asia-east1 \
  --platform managed \
  --allow-unauthenticated
```

## 重要注意事項

### 1. 不變的項目
- **服務 URL**：重新部署後 URL 不會改變
- **環境變數**：會自動保留，不需要重新設定
- **LINE Webhook 設定**：不需要更改

### 2. 部署時間
- **首次部署**：約 2-3 分鐘
- **後續部署**：約 1-2 分鐘（會利用 Docker 層快取加速）

### 3. 費用資訊
- **Google Cloud Run 免費額度**：
  - 每月 200 萬次請求
  - 每月 40 萬 GB-秒 CPU 時間
  - 每月 80 萬 GB-秒記憶體時間
  - 一般小型應用很難超過免費額度

### 4. 故障排除
如果部署失敗，可以檢查建置日誌：
```bash
# 查看最近的建置記錄
gcloud builds list --limit=5

# 查看特定建置的詳細日誌
gcloud builds log BUILD_ID
```

## 服務監控

### 檢查服務狀態
```bash
# 列出所有 Cloud Run 服務
gcloud run services list --platform managed --region asia-east1

# 查看服務詳細資訊
gcloud run services describe chinese-linebot --region asia-east1

# 查看服務日誌
gcloud logs read --service=chinese-linebot --region=asia-east1
```

### 測試服務連線
```bash
# 測試服務是否正常運行
curl https://chinese-linebot-x5h7hw374a-de.a.run.app/
```

## 進階設定

### 自訂服務設定
```bash
# 設定 CPU 和記憶體限制
gcloud run deploy chinese-linebot \
  --source . \
  --region asia-east1 \
  --platform managed \
  --allow-unauthenticated \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10
```

### 設定自訂域名（可選）
```bash
# 對應自訂域名到服務
gcloud run domain-mappings create \
  --service chinese-linebot \
  --domain your-domain.com \
  --region asia-east1
```

這份指南涵蓋了從初次部署到日常維護的完整流程，適合 Google Cloud 新手參考使用。