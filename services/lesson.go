package services

import (
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"

	"chinese-learning-linebot/config"
	"chinese-learning-linebot/models"
)

type LessonService struct {
	firebaseClient *config.FirebaseClient
}

func NewLessonService(firebaseClient *config.FirebaseClient) *LessonService {
	return &LessonService{
		firebaseClient: firebaseClient,
	}
}

func (s *LessonService) GetLearningProgress(publisher string, grade int, semester *int) (*models.LearningProgress, error) {
	// 構建查詢條件
	query := s.firebaseClient.Firestore.Collection("lessons").
		Where("publisher", "==", publisher).
		Where("grade", "==", grade)

	if semester != nil {
		query = query.Where("semester", "==", *semester)
	}

	// 執行查詢
	docs, err := query.Documents(s.firebaseClient.Ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query lessons: %v", err)
	}

	progress := &models.LearningProgress{
		Publisher: publisher,
		Grade:     grade,
		Semester:  semester,
		Lessons:   []models.LessonInfo{},
	}

	totalCharacters := make(map[string]bool)
	for _, doc := range docs {
		data := doc.Data()
		lesson := models.LessonInfo{
			ID:    doc.Ref.ID,
			Title: getStringFromData(data, "title"),
			Unit:  getStringFromData(data, "unit"),
		}

		// 獲取課程中的字符
		if chars, ok := data["characters"].([]interface{}); ok {
			for _, char := range chars {
				if charStr, ok := char.(string); ok {
					lesson.Characters = append(lesson.Characters, charStr)
					totalCharacters[charStr] = true
				}
			}
		}

		lesson.CharacterCount = len(lesson.Characters)
		progress.Lessons = append(progress.Lessons, lesson)
	}

	progress.TotalLessons = len(progress.Lessons)
	progress.TotalCharacters = len(totalCharacters)

	// 計算累積字數（可以從 cumulative_characters collection 獲取更精確的數據）
	cumulativeCount, err := s.getCumulativeCharacterCount(publisher, grade, semester)
	if err == nil {
		progress.CumulativeCharacters = cumulativeCount
	} else {
		progress.CumulativeCharacters = progress.TotalCharacters
	}

	return progress, nil
}

func (s *LessonService) getCumulativeCharacterCount(publisher string, grade int, semester *int) (int, error) {
	// 構建文檔ID
	docID := fmt.Sprintf("%s_%d", publisher, grade)
	if semester != nil {
		docID += "_" + strconv.Itoa(*semester)
	}

	doc, err := s.firebaseClient.Firestore.Collection("cumulative_characters").Doc(docID).Get(s.firebaseClient.Ctx)
	if err != nil {
		return 0, err
	}

	data := doc.Data()
	if count, ok := data["count"].(int64); ok {
		return int(count), nil
	}

	return 0, fmt.Errorf("count field not found or invalid type")
}

func (s *LessonService) GetLessonsByGrade(publisher string, grade int) ([]models.LessonInfo, error) {
	query := s.firebaseClient.Firestore.Collection("lessons").
		Where("publisher", "==", publisher).
		Where("grade", "==", grade).
		OrderBy("semester", firestore.Asc).
		OrderBy("unit", firestore.Asc)

	docs, err := query.Documents(s.firebaseClient.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var lessons []models.LessonInfo
	for _, doc := range docs {
		data := doc.Data()
		lesson := models.LessonInfo{
			ID:    doc.Ref.ID,
			Title: getStringFromData(data, "title"),
			Unit:  getStringFromData(data, "unit"),
		}

		if chars, ok := data["characters"].([]interface{}); ok {
			for _, char := range chars {
				if charStr, ok := char.(string); ok {
					lesson.Characters = append(lesson.Characters, charStr)
				}
			}
		}

		lesson.CharacterCount = len(lesson.Characters)
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (s *LessonService) GetCharactersFromLessons(publisher string, grade int, semester *int) ([]string, error) {
	progress, err := s.GetLearningProgress(publisher, grade, semester)
	if err != nil {
		return nil, err
	}

	characterSet := make(map[string]bool)
	for _, lesson := range progress.Lessons {
		for _, char := range lesson.Characters {
			characterSet[char] = true
		}
	}

	var characters []string
	for char := range characterSet {
		characters = append(characters, char)
	}

	return characters, nil
}

// 輔助函數
func getStringFromData(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}