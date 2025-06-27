package models

// PracticeQuestion 練習題目結構
type PracticeQuestion struct {
	ID            string   `json:"id"`            // 題目ID
	Type          string   `json:"type"`          // 題目類型 (phonetic, stroke, sentence)
	Character     string   `json:"character"`     // 相關字符
	Question      string   `json:"question"`      // 題目內容
	Options       []string `json:"options"`       // 選項（選擇題用）
	CorrectAnswer string   `json:"correctAnswer"` // 正確答案
	Explanation   string   `json:"explanation"`   // 解釋說明
	Difficulty    int      `json:"difficulty"`    // 難度等級
	CreatedAt     int64    `json:"createdAt"`     // 創建時間
}

// PracticeSession 練習會話結構
type PracticeSession struct {
	ID           string             `json:"id"`           // 會話ID
	UserID       string             `json:"userId"`       // 用戶ID
	Type         string             `json:"type"`         // 練習類型
	Questions    []PracticeQuestion `json:"questions"`    // 題目列表
	Answers      []PracticeAnswer   `json:"answers"`      // 答案列表
	Score        int                `json:"score"`        // 得分
	TotalScore   int                `json:"totalScore"`   // 總分
	StartTime    int64              `json:"startTime"`    // 開始時間
	EndTime      int64              `json:"endTime"`      // 結束時間
	Completed    bool               `json:"completed"`    // 是否完成
}

// PracticeAnswer 練習答案結構
type PracticeAnswer struct {
	QuestionID   string `json:"questionId"`   // 題目ID
	UserAnswer   string `json:"userAnswer"`   // 用戶答案
	CorrectAnswer string `json:"correctAnswer"` // 正確答案
	IsCorrect    bool   `json:"isCorrect"`    // 是否正確
	TimeSpent    int64  `json:"timeSpent"`    // 花費時間（毫秒）
	AnsweredAt   int64  `json:"answeredAt"`   // 答題時間
}

// PracticeStats 練習統計結構
type PracticeStats struct {
	UserID              string  `json:"userId"`              // 用戶ID
	TotalSessions       int     `json:"totalSessions"`       // 總練習次數
	TotalQuestions      int     `json:"totalQuestions"`      // 總題目數
	CorrectAnswers      int     `json:"correctAnswers"`      // 正確答案數
	AccuracyRate        float64 `json:"accuracyRate"`        // 正確率
	AverageScore        float64 `json:"averageScore"`        // 平均分數
	TotalTimeSpent      int64   `json:"totalTimeSpent"`      // 總花費時間
	AverageTimePerQuestion int64 `json:"averageTimePerQuestion"` // 平均每題時間
	LastPracticeTime    int64   `json:"lastPracticeTime"`    // 最後練習時間
	Streak              int     `json:"streak"`              // 連續練習天數
	BestStreak          int     `json:"bestStreak"`          // 最佳連續天數
}

// PracticeType 練習類型枚舉
type PracticeType string

const (
	PracticeTypePhonetic PracticeType = "phonetic" // 注音練習
	PracticeTypeStroke   PracticeType = "stroke"   // 筆畫練習
	PracticeTypeSentence PracticeType = "sentence" // 造句練習
	PracticeTypeMixed    PracticeType = "mixed"    // 混合練習
)

// QuestionDifficulty 題目難度枚舉
type QuestionDifficulty int

const (
	DifficultyEasy   QuestionDifficulty = 1 // 簡單
	DifficultyMedium QuestionDifficulty = 2 // 中等
	DifficultyHard   QuestionDifficulty = 3 // 困難
	DifficultyExpert QuestionDifficulty = 4 // 專家
)

// PracticeConfig 練習配置結構
type PracticeConfig struct {
	Type           PracticeType       `json:"type"`           // 練習類型
	Difficulty     QuestionDifficulty `json:"difficulty"`     // 難度等級
	QuestionCount  int                `json:"questionCount"`  // 題目數量
	TimeLimit      int64              `json:"timeLimit"`      // 時間限制（秒）
	RandomOrder    bool               `json:"randomOrder"`    // 是否隨機順序
	ShowExplanation bool              `json:"showExplanation"` // 是否顯示解釋
	Grade          *int               `json:"grade"`          // 指定年級（可選）
	Publisher      string             `json:"publisher"`      // 指定出版社（可選）
}