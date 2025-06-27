package services

import (
	"fmt"
	"math/rand"
	"time"

	"chinese-learning-linebot/config"
	"chinese-learning-linebot/models"
)

type PracticeService struct {
	firebaseClient   *config.FirebaseClient
	characterService *CharacterService
	questionCache    map[string]*models.PracticeQuestion
}

func NewPracticeService(firebaseClient *config.FirebaseClient) *PracticeService {
	return &PracticeService{
		firebaseClient:   firebaseClient,
		characterService: NewCharacterService(firebaseClient),
		questionCache:    make(map[string]*models.PracticeQuestion),
	}
}

func (s *PracticeService) GeneratePhoneticQuestion() (*models.PracticeQuestion, error) {
	// 隨機選擇一個字符
	characters, err := s.characterService.GetRandomCharacters(1)
	if err != nil || len(characters) == 0 {
		return nil, fmt.Errorf("failed to get random character")
	}

	char := characters[0]
	questionID := fmt.Sprintf("phonetic_%d", time.Now().UnixNano())

	// 生成錯誤選項（簡化版本）
	wrongOptions := []string{"ㄅㄚ", "ㄆㄧ", "ㄇㄛ", "ㄈㄟ"}
	options := []string{char.Phonetic}

	// 添加錯誤選項，確保不重複
	for _, wrong := range wrongOptions {
		if wrong != char.Phonetic && len(options) < 4 {
			options = append(options, wrong)
		}
	}

	// 打亂選項順序
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	// 找到正確答案的位置
	correctIndex := 0
	for i, option := range options {
		if option == char.Phonetic {
			correctIndex = i
			break
		}
	}

	question := &models.PracticeQuestion{
		ID:            questionID,
		Type:          "phonetic",
		Character:     char.Character,
		Question:      fmt.Sprintf("請選擇「%s」的正確注音：", char.Character),
		Options:       options,
		CorrectAnswer: fmt.Sprintf("%d", correctIndex),
		Explanation:   fmt.Sprintf("「%s」的注音是「%s」", char.Character, char.Phonetic),
	}

	// 緩存問題
	s.questionCache[questionID] = question

	return question, nil
}

func (s *PracticeService) GenerateStrokeQuestion() (*models.PracticeQuestion, error) {
	// 隨機選擇一個字符
	characters, err := s.characterService.GetRandomCharacters(1)
	if err != nil || len(characters) == 0 {
		return nil, fmt.Errorf("failed to get random character")
	}

	char := characters[0]
	questionID := fmt.Sprintf("stroke_%d", time.Now().UnixNano())

	// 生成錯誤選項（正確答案±1-3）
	correctStrokes := char.StrokeCount
	options := []string{fmt.Sprintf("%d", correctStrokes)}

	// 添加錯誤選項
	for i := 1; i <= 3; i++ {
		if correctStrokes-i > 0 {
			options = append(options, fmt.Sprintf("%d", correctStrokes-i))
		}
		if len(options) < 4 {
			options = append(options, fmt.Sprintf("%d", correctStrokes+i))
		}
	}

	// 確保有4個選項
	for len(options) < 4 {
		options = append(options, fmt.Sprintf("%d", correctStrokes+len(options)))
	}

	// 打亂選項順序
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	// 找到正確答案的位置
	correctIndex := 0
	for i, option := range options {
		if option == fmt.Sprintf("%d", correctStrokes) {
			correctIndex = i
			break
		}
	}

	question := &models.PracticeQuestion{
		ID:            questionID,
		Type:          "stroke",
		Character:     char.Character,
		Question:      fmt.Sprintf("請選擇「%s」的筆畫數：", char.Character),
		Options:       options,
		CorrectAnswer: fmt.Sprintf("%d", correctIndex),
		Explanation:   fmt.Sprintf("「%s」的筆畫數是 %d 畫", char.Character, correctStrokes),
	}

	// 緩存問題
	s.questionCache[questionID] = question

	return question, nil
}

func (s *PracticeService) GenerateSentenceQuestion() (*models.PracticeQuestion, error) {
	// 隨機選擇一個字符
	characters, err := s.characterService.GetRandomCharacters(1)
	if err != nil || len(characters) == 0 {
		return nil, fmt.Errorf("failed to get random character")
	}

	char := characters[0]
	questionID := fmt.Sprintf("sentence_%d", time.Now().UnixNano())

	question := &models.PracticeQuestion{
		ID:          questionID,
		Type:        "sentence",
		Character:   char.Character,
		Question:    fmt.Sprintf("請用「%s」造句：", char.Character),
		Options:     []string{}, // 造句題沒有選項
		Explanation: fmt.Sprintf("很好！你用「%s」造了一個句子。", char.Character),
	}

	// 緩存問題
	s.questionCache[questionID] = question

	return question, nil
}

func (s *PracticeService) CheckAnswer(questionID, answer string) (bool, string, error) {
	question, exists := s.questionCache[questionID]
	if !exists {
		return false, "問題已過期，請重新開始練習", nil
	}

	switch question.Type {
	case "phonetic", "stroke":
		isCorrect := answer == question.CorrectAnswer
		return isCorrect, question.Explanation, nil

	case "sentence":
		// 造句題目前簡單判斷是否包含該字符
		containsChar := false
		for _, char := range answer {
			if string(char) == question.Character {
				containsChar = true
				break
			}
		}
		if containsChar {
			return true, fmt.Sprintf("很棒的句子！你成功使用了「%s」這個字。", question.Character), nil
		} else {
			return false, fmt.Sprintf("請確保句子中包含「%s」這個字。", question.Character), nil
		}

	default:
		return false, "未知的問題類型", nil
	}
}

func (s *PracticeService) CleanupExpiredQuestions() {
	// 清理過期的問題緩存（可以定期調用）
	currentTime := time.Now().UnixNano()
	for id := range s.questionCache {
		// 簡單的過期邏輯：ID中的時間戳超過1小時就刪除
		if currentTime-extractTimestampFromID(id) > int64(time.Hour) {
			delete(s.questionCache, id)
		}
	}
}

// 輔助函數
func extractTimestampFromID(id string) int64 {
	// 從ID中提取時間戳（簡化實現）
	// 實際實現可能需要更複雜的解析
	return time.Now().UnixNano() // 暫時返回當前時間
}