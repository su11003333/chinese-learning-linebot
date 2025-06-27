# 中文學習 LINE Bot 專案概述

## 專案簡介

這是一個基於 Go 語言開發的中文學習 LINE Bot，旨在幫助用戶學習中文字詞、練習發音和筆劃，並追蹤學習進度。專案整合了 LINE Bot SDK、Firebase Firestore 資料庫，提供互動式的中文學習體驗。

## 技術架構

### 核心技術棧
- **後端語言**: Go (Golang)
- **Web 框架**: Gin
- **資料庫**: Firebase Firestore
- **訊息平台**: LINE Bot SDK v7
- **部署**: Docker + Google Cloud Run
- **環境管理**: godotenv

### 專案結構
```
chinese-learning-linebot/
├── main.go                 # 應用程式入口點
├── config/                 # 配置檔案
│   ├── firebase.go         # Firebase 初始化設定
│   └── line.go            # LINE Bot 初始化設定
├── handlers/              # HTTP 請求處理器
│   ├── webhook.go         # Webhook 路由處理
│   ├── message.go         # 訊息事件處理
│   └── postback.go        # 回調事件處理
├── models/                # 資料模型定義
│   ├── character.go       # 字詞資料結構
│   ├── lesson.go          # 課程資料結構
│   └── practice.go        # 練習資料結構
├── services/              # 業務邏輯服務
│   ├── character.go       # 字詞查詢服務
│   ├── lesson.go          # 課程管理服務
│   └── practice.go        # 練習生成服務
├── utils/                 # 工具函數
│   ├── response.go        # LINE 訊息回應工具
│   └── string.go          # 字串處理工具
├── .env.example           # 環境變數範例
├── Dockerfile             # Docker 容器配置
├── go.mod                 # Go 模組依賴
└── README.md              # 專案說明文件
```

## 核心功能

### 1. 字詞查詢功能
- **功能描述**: 用戶可以輸入中文字詞進行查詢
- **回應內容**: 提供注音、拼音、筆劃數、部首、解釋說明
- **實作位置**: `services/character.go`, `handlers/message.go`

### 2. 學習進度追蹤
- **功能描述**: 記錄用戶學習歷史和進度
- **資料追蹤**: 學習時間、完成課程、練習成績
- **實作位置**: `models/practice.go`, `services/practice.go`

### 3. 互動式練習
- **注音練習**: 測試用戶對字詞注音的掌握
- **筆劃練習**: 練習字詞的筆劃順序和數量
- **選擇題模式**: 提供多選項答案供用戶選擇
- **實作位置**: `handlers/postback.go`, `utils/response.go`

### 4. 課程管理
- **分級學習**: 按年級和學期組織學習內容
- **單元劃分**: 將學習內容分為不同單元
- **進度管理**: 追蹤各課程的學習進度
- **實作位置**: `models/lesson.go`, `services/lesson.go`

## 資料模型

### Character (字詞資料)
```go
type Character struct {
    ID           string   `json:"id"`           // 字詞ID
    Character    string   `json:"character"`    // 中文字詞
    Phonetic     string   `json:"phonetic"`     // 注音符號
    Pinyin       string   `json:"pinyin"`       // 拼音
    StrokeCount  int      `json:"strokeCount"`  // 筆劃數
    Radical      string   `json:"radical"`      // 部首
    Explanation  string   `json:"explanation"`  // 解釋說明
    Difficulty   int      `json:"difficulty"`   // 難度等級
    CreatedAt    int64    `json:"createdAt"`    // 創建時間
}
```

### Lesson (課程資料)
```go
type Lesson struct {
    ID          string   `json:"id"`          // 課程ID
    Title       string   `json:"title"`       // 課程標題
    Grade       int      `json:"grade"`       // 年級
    Semester    int      `json:"semester"`    // 學期
    Unit        int      `json:"unit"`        // 單元
    Characters  []string `json:"characters"`  // 包含的字詞ID列表
    CreatedAt   int64    `json:"createdAt"`   // 創建時間
}
```

### PracticeSession (練習會話)
```go
type PracticeSession struct {
    ID           string             `json:"id"`           // 會話ID
    UserID       string             `json:"userId"`       // 用戶ID
    Type         string             `json:"type"`         // 練習類型
    Questions    []PracticeQuestion `json:"questions"`    // 練習題目
    CurrentIndex int                `json:"currentIndex"` // 當前題目索引
    Score        int                `json:"score"`        // 得分
    CreatedAt    int64              `json:"createdAt"`    // 創建時間
}
```

## LINE Bot 互動流程

### 1. 訊息處理流程
1. 用戶發送訊息到 LINE
2. LINE 平台轉發到 Webhook (`/webhook`)
3. `handlers/webhook.go` 接收並分發事件
4. `handlers/message.go` 處理文字訊息
5. 根據訊息內容調用相應服務
6. 回傳格式化的回應訊息

### 2. 回調處理流程
1. 用戶點擊快速回覆或按鈕
2. 觸發 Postback 事件
3. `handlers/postback.go` 解析回調資料
4. 執行對應的業務邏輯
5. 更新練習狀態或進度
6. 回傳下一步操作或結果

## 環境配置

### 必要環境變數
```env
# LINE Bot 設定
CHANNEL_SECRET=your_channel_secret
CHANNEL_ACCESS_TOKEN=your_channel_access_token

# Firebase 設定
GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account.json
FIREBASE_PROJECT_ID=your_project_id

# 應用程式設定
PORT=8080
GIN_MODE=release
```

## 部署說明

### Docker 部署
1. 建立 Docker 映像檔
2. 設定環境變數
3. 部署到 Google Cloud Run
4. 設定 Webhook URL

### 本地開發
1. 安裝 Go 1.19+
2. 複製 `.env.example` 為 `.env`
3. 設定環境變數
4. 執行 `go run main.go`

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

### 擴展功能建議
- 語音辨識練習
- 字詞聯想遊戲
- 學習統計圖表
- 多人競賽模式
- 自訂學習計畫

## 專案特色

1. **模組化架構**: 清晰的分層設計，易於維護和擴展
2. **互動式學習**: 豐富的 LINE Bot 互動功能
3. **資料驅動**: 基於 Firestore 的靈活資料管理
4. **雲端部署**: 支援 Docker 和 Cloud Run 部署
5. **中文特化**: 專門針對中文學習需求設計

這個專案為中文學習提供了一個完整的數位化解決方案，結合了現代化的技術棧和教育需求，是一個具有實用價值的學習工具。