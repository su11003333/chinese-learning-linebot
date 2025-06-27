package utils

import (
	"regexp"
	"unicode"
)

// IsChineseCharacter 檢查是否為中文字符
func IsChineseCharacter(char rune) bool {
	return unicode.Is(unicode.Scripts["Han"], char)
}

// ExtractChineseCharacters 從文本中提取中文字符
func ExtractChineseCharacters(text string) []string {
	var characters []string
	for _, char := range text {
		if IsChineseCharacter(char) {
			characters = append(characters, string(char))
		}
	}
	return characters
}

// ContainsChineseCharacters 檢查文本是否包含中文字符
func ContainsChineseCharacters(text string) bool {
	for _, char := range text {
		if IsChineseCharacter(char) {
			return true
		}
	}
	return false
}

// GetFirstChineseCharacter 獲取文本中第一個中文字符
func GetFirstChineseCharacter(text string) string {
	for _, char := range text {
		if IsChineseCharacter(char) {
			return string(char)
		}
	}
	return ""
}

// ValidatePhoneticNotation 驗證注音符號格式
func ValidatePhoneticNotation(phonetic string) bool {
	// 基本的注音符號正則表達式
	phoneticPattern := `^[ㄅ-ㄩˊˇˋ˙]+$`
	matched, _ := regexp.MatchString(phoneticPattern, phonetic)
	return matched
}

// CleanText 清理文本，移除多餘的空白字符
func CleanText(text string) string {
	// 移除前後空白
	text = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(text, "")
	// 將多個連續空白替換為單個空格
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return text
}

// TruncateText 截斷文本到指定長度，並添加省略號
func TruncateText(text string, maxLength int) string {
	if len([]rune(text)) <= maxLength {
		return text
	}
	runes := []rune(text)
	if maxLength <= 3 {
		return string(runes[:maxLength])
	}
	return string(runes[:maxLength-3]) + "..."
}

// SplitByLength 按指定長度分割文本
func SplitByLength(text string, length int) []string {
	if length <= 0 {
		return []string{text}
	}

	runes := []rune(text)
	var result []string

	for i := 0; i < len(runes); i += length {
		end := i + length
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}

	return result
}