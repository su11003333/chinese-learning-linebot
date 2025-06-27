package services

import (
	"fmt"

	"chinese-learning-linebot/config"
	"chinese-learning-linebot/models"
)

type CharacterService struct {
	firebaseClient *config.FirebaseClient
}

func NewCharacterService(firebaseClient *config.FirebaseClient) *CharacterService {
	return &CharacterService{
		firebaseClient: firebaseClient,
	}
}

func (s *CharacterService) LookupCharacter(char string) (*models.CharacterInfo, error) {
	// 查詢單個字符
	doc, err := s.firebaseClient.Firestore.Collection("characters").Doc(char).Get(s.firebaseClient.Ctx)
	if err != nil {
		return nil, fmt.Errorf("character not found: %v", err)
	}

	var character models.CharacterInfo
	if err := doc.DataTo(&character); err != nil {
		return nil, fmt.Errorf("failed to parse character data: %v", err)
	}

	// 設置字符本身
	character.Character = char

	// 查詢該字符出現的課程
	lessons, err := s.getLessonsForCharacter(char)
	if err != nil {
		// 即使查詢課程失敗，仍然返回字符基本信息
		character.Lessons = []string{}
	} else {
		character.Lessons = lessons
	}

	return &character, nil
}

func (s *CharacterService) getLessonsForCharacter(char string) ([]string, error) {
	// 查詢包含該字符的課程
	query := s.firebaseClient.Firestore.Collection("lessons").Where("characters", "array-contains", char).Limit(10)
	docs, err := query.Documents(s.firebaseClient.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var lessons []string
	for _, doc := range docs {
		data := doc.Data()
		if title, ok := data["title"].(string); ok {
			lessons = append(lessons, title)
		}
	}

	return lessons, nil
}

func (s *CharacterService) SearchCharacters(keyword string, limit int) ([]*models.CharacterInfo, error) {
	// 模糊搜索字符（可以根據注音、部首等搜索）
	if limit <= 0 {
		limit = 10
	}

	// 這裡可以實現更複雜的搜索邏輯
	// 目前先實現簡單的精確匹配
	character, err := s.LookupCharacter(keyword)
	if err != nil {
		return []*models.CharacterInfo{}, nil
	}

	return []*models.CharacterInfo{character}, nil
}

func (s *CharacterService) GetRandomCharacters(count int) ([]*models.CharacterInfo, error) {
	// 隨機獲取字符（用於練習）
	if count <= 0 {
		count = 5
	}

	// 這裡可以實現隨機查詢邏輯
	// 由於 Firestore 的限制，這裡使用簡化的實現
	commonChars := []string{"學", "習", "中", "文", "字", "詞", "語", "言", "書", "本"}

	var characters []*models.CharacterInfo
	for i := 0; i < count && i < len(commonChars); i++ {
		char, err := s.LookupCharacter(commonChars[i])
		if err == nil {
			characters = append(characters, char)
		}
	}

	return characters, nil
}