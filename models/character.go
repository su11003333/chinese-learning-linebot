package models

// CharacterInfo 字詞資訊結構
type CharacterInfo struct {
	Character   string   `json:"character" firestore:"character"`     // 字符本身
	Phonetic    string   `json:"phonetic" firestore:"phonetic"`       // 注音
	StrokeCount int      `json:"strokeCount" firestore:"strokeCount"` // 筆畫數
	Radical     string   `json:"radical" firestore:"radical"`         // 部首
	Meaning     string   `json:"meaning" firestore:"meaning"`         // 字義
	Examples    []string `json:"examples" firestore:"examples"`       // 例句
	Lessons     []string `json:"lessons"`                             // 出現的課程（查詢時填入）
	Frequency   int      `json:"frequency" firestore:"frequency"`     // 使用頻率
	Difficulty  int      `json:"difficulty" firestore:"difficulty"`   // 難度等級 (1-5)
	CreatedAt   int64    `json:"createdAt" firestore:"createdAt"`     // 創建時間
	UpdatedAt   int64    `json:"updatedAt" firestore:"updatedAt"`     // 更新時間
}

// CharacterSearchResult 字詞搜索結果
type CharacterSearchResult struct {
	Characters []*CharacterInfo `json:"characters"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
}

// CharacterStats 字詞統計信息
type CharacterStats struct {
	TotalCharacters     int `json:"totalCharacters"`
	LearnedCharacters   int `json:"learnedCharacters"`
	RemainingCharacters int `json:"remainingCharacters"`
	ProgressPercentage  int `json:"progressPercentage"`
}