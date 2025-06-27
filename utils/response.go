
package utils

import (
	"fmt"
	"strings"
)

// CreateCumulativeQueryResultMessage 創建累積字詞查詢結果的訊息
func CreateCumulativeQueryResultMessage(publisher string, grade int, semester int, lesson int, learned []string, notLearned []string) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("📚 累積字詞查詢結果\n\n"))
	result.WriteString(fmt.Sprintf("📖 範圍：%s %d年級第%d學期第%d課\n\n", publisher, grade, semester, lesson))

	if len(learned) > 0 {
		result.WriteString("✅ 已學過的字詞：\n")
		result.WriteString(strings.Join(learned, "、"))
		result.WriteString("\n\n")
	}

	if len(notLearned) > 0 {
		result.WriteString("❌ 尚未學過的字詞：\n")
		result.WriteString(strings.Join(notLearned, "、"))
		result.WriteString("\n\n")
	}

	result.WriteString(fmt.Sprintf("📊 統計：已學 %d 字，未學 %d 字", len(learned), len(notLearned)))

	return result.String()
}

// CreatePublisherSelectionMessage 創建出版社選擇訊息
func CreatePublisherSelectionMessage() string {
	return "📚 請選擇出版社：\n\n" +
		"1️⃣ 康軒\n" +
		"2️⃣ 南一\n" +
		"3️⃣ 翰林\n\n" +
		"請輸入數字 1、2 或 3"
}

// CreateGradeSelectionMessage 創建年級選擇訊息
func CreateGradeSelectionMessage() string {
	return "📖 請選擇年級：\n\n" +
		"1️⃣ 一年級\n" +
		"2️⃣ 二年級\n" +
		"3️⃣ 三年級\n" +
		"4️⃣ 四年級\n" +
		"5️⃣ 五年級\n" +
		"6️⃣ 六年級\n\n" +
		"請輸入數字 1-6"
}

// CreateSemesterSelectionMessage 創建學期選擇訊息
func CreateSemesterSelectionMessage() string {
	return "📅 請選擇學期：\n\n" +
		"1️⃣ 上學期\n" +
		"2️⃣ 下學期\n\n" +
		"請輸入數字 1 或 2"
}

// CreateLessonSelectionMessage 創建課次選擇訊息
func CreateLessonSelectionMessage() string {
	return "📝 請輸入課次（例如：5）：\n\n" +
		"請輸入 1-20 之間的數字"
}

// CreateCharacterInputMessage 創建字詞輸入訊息
func CreateCharacterInputMessage(publisher string, grade int, semester int, lesson int) string {
	return fmt.Sprintf("✏️ 請輸入要查詢的字詞：\n\n"+
		"📖 查詢範圍：%s %d年級第%d學期第%d課\n\n"+
		"例如：我好喜歡吃飯配菜",
		publisher, grade, semester, lesson)
}