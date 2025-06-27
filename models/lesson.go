package models

// LessonInfo 課程資訊結構
type LessonInfo struct {
	ID             string   `json:"id" firestore:"id"`                         // 課程ID
	Title          string   `json:"title" firestore:"title"`                   // 課程標題
	Unit           string   `json:"unit" firestore:"unit"`                     // 單元
	Publisher      string   `json:"publisher" firestore:"publisher"`           // 出版社
	Grade          int      `json:"grade" firestore:"grade"`                   // 年級
	Semester       int      `json:"semester" firestore:"semester"`             // 學期
	Characters     []string `json:"characters" firestore:"characters"`         // 課程中的字符
	CharacterCount int      `json:"characterCount"`                            // 字符數量（計算得出）
	Description    string   `json:"description" firestore:"description"`       // 課程描述
	Objectives     []string `json:"objectives" firestore:"objectives"`         // 學習目標
	Difficulty     int      `json:"difficulty" firestore:"difficulty"`         // 難度等級
	Order          int      `json:"order" firestore:"order"`                   // 課程順序
	CreatedAt      int64    `json:"createdAt" firestore:"createdAt"`           // 創建時間
	UpdatedAt      int64    `json:"updatedAt" firestore:"updatedAt"`           // 更新時間
}

// LearningProgress 學習進度結構
type LearningProgress struct {
	Publisher            string       `json:"publisher"`            // 出版社
	Grade                int          `json:"grade"`                // 年級
	Semester             *int         `json:"semester"`             // 學期（可為空表示全年）
	Lessons              []LessonInfo `json:"lessons"`              // 課程列表
	TotalLessons         int          `json:"totalLessons"`         // 總課程數
	TotalCharacters      int          `json:"totalCharacters"`      // 總字符數
	CumulativeCharacters int          `json:"cumulativeCharacters"` // 累積字符數
	CompletedLessons     int          `json:"completedLessons"`     // 已完成課程數
	ProgressPercentage   float64      `json:"progressPercentage"`   // 進度百分比
	LastUpdated          int64        `json:"lastUpdated"`          // 最後更新時間
}

// CumulativeCharacters 累積字數結構
type CumulativeCharacters struct {
	ID         string `json:"id" firestore:"id"`                 // 文檔ID (publisher_grade_semester)
	Publisher  string `json:"publisher" firestore:"publisher"`   // 出版社
	Grade      int    `json:"grade" firestore:"grade"`           // 年級
	Semester   int    `json:"semester" firestore:"semester"`     // 學期
	Count      int    `json:"count" firestore:"count"`           // 累積字數
	Characters []string `json:"characters" firestore:"characters"` // 字符列表
	CreatedAt  int64  `json:"createdAt" firestore:"createdAt"`   // 創建時間
	UpdatedAt  int64  `json:"updatedAt" firestore:"updatedAt"`   // 更新時間
}

// LessonSearchCriteria 課程搜索條件
type LessonSearchCriteria struct {
	Publisher string `json:"publisher"`
	Grade     *int   `json:"grade"`
	Semester  *int   `json:"semester"`
	Keyword   string `json:"keyword"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// PublisherInfo 出版社資訊
type PublisherInfo struct {
	Name        string `json:"name"`        // 出版社名稱
	DisplayName string `json:"displayName"` // 顯示名稱
	Description string `json:"description"` // 描述
	Grades      []int  `json:"grades"`      // 支援的年級
	Active      bool   `json:"active"`      // 是否啟用
}