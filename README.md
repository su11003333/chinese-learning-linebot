# 中文學習 LINE Bot

這是一個專為家長設計的 LINE Bot，可以快速回答關於孩童中文字詞學習進度的問題。

## 功能特色

### 🔍 字詞查詢
- 直接輸入中文字即可查詢詳細資訊
- 顯示注音、筆畫、部首、字義和例句
- 顯示該字出現在哪些課程中

### 📚 課程進度查詢
- 按出版社、年級、學期查詢學習進度
- 顯示課程數量、新學字數、累積字數
- 計算學習進度百分比

### 🎮 互動練習
- 注音練習：選擇正確的注音符號
- 筆畫練習：猜測字的筆畫數
- 造句練習：完成句子填空

### 📊 學習統計
- 追蹤學習進度
- 顯示累積學習字數
- 提供個人化學習建議

## 技術架構

### 後端技術
- **Go 1.21+**: 主要開發語言
- **Gin**: Web 框架
- **LINE Bot SDK**: LINE 機器人開發
- **Firebase Firestore**: 資料庫
- **Firebase Auth**: 用戶認證

### 專案結構
```
chinese-learning-linebot/
├── main.go                 # 應用程式入口
├── config/                 # 配置文件
│   ├── firebase.go         # Firebase 初始化
│   └── line.go            # LINE Bot 初始化
├── handlers/              # 請求處理器
│   ├── webhook.go         # Webhook 處理
│   ├── message.go         # 訊息處理
│   └── postback.go        # 回調處理
├── services/              # 業務邏輯
│   ├── character.go       # 字詞服務
│   ├── lesson.go          # 課程服務
│   └── practice.go        # 練習服務
├── models/                # 資料模型
│   ├── character.go       # 字詞模型
│   ├── lesson.go          # 課程模型
│   └── practice.go        # 練習模型
├── utils/                 # 工具函數
│   ├── response.go        # 回應格式化
│   └── string.go          # 字串處理
├── go.mod                 # Go 模組定義
├── go.sum                 # 依賴版本鎖定
├── .env.example           # 環境變數範例
└── README.md              # 專案說明
```

## 安裝與設定

### 1. 環境需求
- Go 1.21 或更高版本
- Firebase 專案
- LINE Developers 帳號

### 2. 複製專案
```bash
git clone <repository-url>
cd chinese-learning-linebot
```

### 3. 安裝依賴
```bash
go mod tidy
```

### 4. 環境變數設定
複製 `.env.example` 為 `.env` 並填入相關設定：

```bash
cp .env.example .env
```

編輯 `.env` 文件：
```env
# LINE Bot Configuration
LINE_CHANNEL_SECRET=your_line_channel_secret_here
LINE_CHANNEL_ACCESS_TOKEN=your_line_channel_access_token_here

# Firebase Configuration
FIREBASE_PROJECT_ID=your_firebase_project_id
GOOGLE_APPLICATION_CREDENTIALS=path/to/your/firebase-service-account-key.json

# Server Configuration
PORT=8080
GIN_MODE=release
```

### 5. Firebase 設定
1. 在 [Firebase Console](https://console.firebase.google.com/) 建立新專案
2. 啟用 Firestore 資料庫
3. 下載服務帳號金鑰 JSON 文件
4. 將金鑰文件路徑設定到 `GOOGLE_APPLICATION_CREDENTIALS`

### 6. LINE Bot 設定
1. 在 [LINE Developers](https://developers.line.biz/) 建立新的 Provider 和 Channel
2. 取得 Channel Secret 和 Channel Access Token
3. 設定 Webhook URL: `https://your-domain.com/webhook`

## 執行應用程式

### 開發模式
```bash
go run main.go
```

### 生產模式
```bash
go build -o linebot main.go
./linebot
```

## API 端點

- `GET /health` - 健康檢查
- `POST /webhook` - LINE Bot Webhook

## 資料庫結構

### Characters Collection
```json
{
  "character": "學",
  "phonetic": "ㄒㄩㄝˊ",
  "strokeCount": 16,
  "radical": "子",
  "meaning": "學習、求知",
  "examples": ["學習", "學校", "學生"],
  "lessons": ["康軒三上第一課", "翰林三下第五課"]
}
```

### Lessons Collection
```json
{
  "publisher": "康軒",
  "grade": "3",
  "semester": "1",
  "lesson": "1",
  "title": "時間是什麼",
  "characters": ["時", "間", "什", "麼"]
}
```

### CumulativeCharacters Collection
```json
{
  "publisher": "康軒",
  "grade": "3",
  "semester": "1",
  "cumulativeCount": 450
}
```

## 使用方式

### 字詞查詢
直接在 LINE 中輸入中文字，例如：
- 輸入「學」→ 顯示「學」字的詳細資訊

### 課程查詢
1. 輸入「課程查詢」或點擊選單
2. 選擇出版社（康軒、翰林、南一）
3. 選擇年級（1-6年級）
4. 選擇學期（上學期、下學期、全年）
5. 查看學習進度統計

### 練習模式
1. 輸入「練習模式」或點擊選單
2. 選擇練習類型：
   - 注音練習
   - 筆畫練習
   - 造句練習
3. 回答問題並獲得即時反饋

## 部署

### 使用 Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o linebot main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/linebot .
CMD ["./linebot"]
```

### 使用 Google Cloud Run
1. 建立 `cloudbuild.yaml`
2. 設定環境變數
3. 部署到 Cloud Run

## 開發指南

### 新增字詞資料
1. 在 Firestore 的 `characters` collection 中新增文件
2. 確保包含所有必要欄位
3. 更新相關課程的字詞列表

### 新增練習類型
1. 在 `models/practice.go` 中定義新的練習類型
2. 在 `services/practice.go` 中實作生成邏輯
3. 在 `handlers/postback.go` 中處理回答
4. 在 `utils/response.go` 中建立顯示格式

## 貢獻

歡迎提交 Issue 和 Pull Request 來改善這個專案。

## 授權

MIT License