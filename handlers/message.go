package handlers

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"chinese-learning-linebot/config"
)

// 用戶狀態管理
type UserState struct {
	Mode      string // "cumulative_query" 或 ""
	Publisher string
	Grade     int
	Semester  int
	Lesson    int
	Step      int // 0: 等待出版社, 1: 等待年級, 2: 等待學期, 3: 等待課次, 4: 等待查詢字詞
	// 用戶偏好設定（記憶半年）
	PreferredPublisher string
	PreferredGrade     int
	PreferredSemester  int
}

// 從 Firestore 獲取用戶狀態
func getUserState(firebaseClient *config.FirebaseClient, userID string) *UserState {
	doc, err := firebaseClient.Firestore.Collection("user_states").Doc(userID).Get(firebaseClient.Ctx)
	if err != nil {
		// 如果文檔不存在或發生錯誤，返回空狀態
		return &UserState{}
	}

	var state UserState
	if err := doc.DataTo(&state); err != nil {
		// 如果解析失敗，返回空狀態
		return &UserState{}
	}

	return &state
}

// 設置用戶狀態到 Firestore
func setUserState(firebaseClient *config.FirebaseClient, userID string, state *UserState) {
	_, err := firebaseClient.Firestore.Collection("user_states").Doc(userID).Set(firebaseClient.Ctx, state)
	if err != nil {
		log.Printf("Error setting user state: %v", err)
	}
}

// 清除用戶狀態從 Firestore
func clearUserState(firebaseClient *config.FirebaseClient, userID string) {
	_, err := firebaseClient.Firestore.Collection("user_states").Doc(userID).Delete(firebaseClient.Ctx)
	if err != nil {
		log.Printf("Error clearing user state: %v", err)
	}
}

func handleMessage(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		return handleTextMessage(event, message, bot, firebaseClient)
	default:
		return replyMessage(event, bot, "抱歉，我只能處理文字訊息。")
	}
}

func handleTextMessage(event *linebot.Event, message *linebot.TextMessage, bot *linebot.Client, firebaseClient *config.FirebaseClient) error {
	userText := strings.TrimSpace(message.Text)
	userID := event.Source.UserID
	state := getUserState(firebaseClient, userID)

	// 處理退出指令
	if userText == "退出" {
		// 只清除當前查詢狀態，保留用戶偏好設定
		state.Mode = ""
		state.Publisher = ""
		state.Grade = 0
		state.Semester = 0
		state.Lesson = 0
		state.Step = 0
		// 保留 PreferredPublisher, PreferredGrade, PreferredSemester
		setUserState(firebaseClient, userID, state)
		return replyMessage(event, bot, "已退出當前模式，請輸入新的指令。")
	}

	// 如果用戶在累積字詞查詢模式中
	if state.Mode == "cumulative_query" {
		return handleCumulativeQueryMode(event, userText, bot, firebaseClient, userID, state)
	}

	// 處理新指令
	switch userText {
	case "查詢累積字詞":
		return startCumulativeQuery(event, bot, firebaseClient, userID)
	case "重設偏好", "重設設定", "清除記憶":
		return resetUserPreferences(event, bot, firebaseClient, userID)
	case "使用者課程設定", "查看設定", "我的設定":
		return showUserSettings(event, bot, firebaseClient, userID)
	case "印字帖":
		return handlePrintWorksheet(event, bot, firebaseClient, userID)
	case "平板學寫字":
		return handleTabletPractice(event, bot)
	case "幫助", "help", "說明":
		return handleHelp(event, bot)
	default:
		return handleUnknownMessage(event, bot)
	}
}

// 開始累積字詞查詢模式
func startCumulativeQuery(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient, userID string) error {
	// 檢查是否有現有的用戶偏好設定
	existingState := getUserState(firebaseClient, userID)
	
	// 調試日誌
	log.Printf("User %s existing state: Publisher=%s, Grade=%d, Semester=%d", userID, existingState.PreferredPublisher, existingState.PreferredGrade, existingState.PreferredSemester)
	
	state := &UserState{
		Mode: "cumulative_query",
	}
	
	// 如果用戶已有偏好設定，提供三個選項
	if existingState.PreferredPublisher != "" && existingState.PreferredGrade > 0 && existingState.PreferredSemester > 0 {
		state.Publisher = existingState.PreferredPublisher
		state.Grade = existingState.PreferredGrade
		state.Semester = existingState.PreferredSemester
		state.PreferredPublisher = existingState.PreferredPublisher
		state.PreferredGrade = existingState.PreferredGrade
		state.PreferredSemester = existingState.PreferredSemester
		state.Step = -1 // 特殊步驟：等待用戶選擇操作
		setUserState(firebaseClient, userID, state)
		
		semesterText := "上學期"
		if state.Semester == 2 {
			semesterText = "下學期"
		}
		
		// 創建選項快速回覆
		quickReply := &linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{
					Action: &linebot.MessageAction{
						Label: "照用上次設定",
						Text:  "照用上次設定",
					},
				},
				{
					Action: &linebot.MessageAction{
						Label: "修改課程",
						Text:  "修改課程",
					},
				},
				{
					Action: &linebot.MessageAction{
						Label: "重新設定",
						Text:  "重新設定",
					},
				},
			},
		}
		
		return replyMessageWithQuickReply(event, bot, fmt.Sprintf("📚 累積字詞查詢\n\n已記憶的設定：%s %d年級%s\n\n請選擇操作：", state.Publisher, state.Grade, semesterText), quickReply)
	}
	
	// 沒有偏好設定，從頭開始
	// 保留現有的偏好設定，只重設查詢相關的欄位
	state.Publisher = ""
	state.Grade = 0
	state.Semester = 0
	state.Lesson = 0
	state.Step = 0
	// 保留現有的偏好設定
	state.PreferredPublisher = existingState.PreferredPublisher
	state.PreferredGrade = existingState.PreferredGrade
	state.PreferredSemester = existingState.PreferredSemester
	setUserState(firebaseClient, userID, state)

	// 創建出版社選擇快速回覆
	quickReply := &linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{
				Action: &linebot.MessageAction{
					Label: "康軒",
					Text:  "康軒",
				},
			},
			{
				Action: &linebot.MessageAction{
					Label: "南一",
					Text:  "南一",
				},
			},
			{
				Action: &linebot.MessageAction{
					Label: "翰林",
					Text:  "翰林",
				},
			},
		},
	}

	return replyMessageWithQuickReply(event, bot, "📚 累積字詞查詢\n\n請選擇出版社：", quickReply)
}

// 處理累積字詞查詢模式的狀態機
func handleCumulativeQueryMode(event *linebot.Event, userText string, bot *linebot.Client, firebaseClient *config.FirebaseClient, userID string, state *UserState) error {
	switch state.Step {
	case -1: // 等待用戶選擇操作（照用上次設定、修改課程、重新設定）
		switch userText {
		case "照用上次設定":
			// 直接跳到課次輸入
			state.Step = 3
			setUserState(firebaseClient, userID, state)
			semesterText := "上學期"
			if state.Semester == 2 {
				semesterText = "下學期"
			}
			return replyMessage(event, bot, fmt.Sprintf("✅ 使用已記憶的設定：%s %d年級%s\n\n請輸入課次（例如：5）：", state.Publisher, state.Grade, semesterText))
			
		case "修改課程":
			// 直接跳到課次輸入，但保持當前設定
			state.Step = 3
			setUserState(firebaseClient, userID, state)
			semesterText := "上學期"
			if state.Semester == 2 {
				semesterText = "下學期"
			}
			return replyMessage(event, bot, fmt.Sprintf("📝 修改課程模式\n\n當前設定：%s %d年級%s\n\n請輸入課次（例如：5）：", state.Publisher, state.Grade, semesterText))
			
		case "重新設定":
			// 清除當前設定，重新開始
			state.Publisher = ""
			state.Grade = 0
			state.Semester = 0
			state.Step = 0
			setUserState(firebaseClient, userID, state)
			
			// 創建出版社選擇快速回覆
			quickReply := &linebot.QuickReplyItems{
				Items: []*linebot.QuickReplyButton{
					{
						Action: &linebot.MessageAction{
							Label: "康軒",
							Text:  "康軒",
						},
					},
					{
						Action: &linebot.MessageAction{
							Label: "南一",
							Text:  "南一",
						},
					},
					{
						Action: &linebot.MessageAction{
							Label: "翰林",
							Text:  "翰林",
						},
					},
				},
			}
			return replyMessageWithQuickReply(event, bot, "🔄 重新設定\n\n📚 累積字詞查詢\n\n請選擇出版社：", quickReply)
			
		default:
			return replyMessage(event, bot, "請選擇：照用上次設定、修改課程、或重新設定")
		}
		
	case 0: // 等待出版社
		if userText == "康軒" || userText == "南一" || userText == "翰林" {
			state.Publisher = userText
			state.Step = 1
			setUserState(firebaseClient, userID, state)

			// 創建年級選擇快速回覆
			quickReply := &linebot.QuickReplyItems{
				Items: []*linebot.QuickReplyButton{
					{Action: &linebot.MessageAction{Label: "1年級", Text: "1"}},
					{Action: &linebot.MessageAction{Label: "2年級", Text: "2"}},
					{Action: &linebot.MessageAction{Label: "3年級", Text: "3"}},
					{Action: &linebot.MessageAction{Label: "4年級", Text: "4"}},
					{Action: &linebot.MessageAction{Label: "5年級", Text: "5"}},
					{Action: &linebot.MessageAction{Label: "6年級", Text: "6"}},
				},
			}
			return replyMessageWithQuickReply(event, bot, fmt.Sprintf("已選擇：%s\n\n請選擇年級：", userText), quickReply)
		} else {
			return replyMessage(event, bot, "請選擇正確的出版社：康軒、南一、翰林")
		}

	case 1: // 等待年級
		if grade := parseGrade(userText); grade > 0 && grade <= 6 {
			state.Grade = grade
			state.Step = 2
			setUserState(firebaseClient, userID, state)

			// 創建學期選擇快速回覆
			quickReply := &linebot.QuickReplyItems{
				Items: []*linebot.QuickReplyButton{
					{Action: &linebot.MessageAction{Label: "上學期", Text: "1"}},
					{Action: &linebot.MessageAction{Label: "下學期", Text: "2"}},
				},
			}
			return replyMessageWithQuickReply(event, bot, fmt.Sprintf("已選擇：%s %d年級\n\n請選擇學期：", state.Publisher, grade), quickReply)
		} else {
			return replyMessage(event, bot, "請輸入正確的年級數字（1-6）")
		}

	case 2: // 等待學期
		if semester := parseSemester(userText); semester == 1 || semester == 2 {
			state.Semester = semester
			// 保存用戶偏好設定
			state.PreferredPublisher = state.Publisher
			state.PreferredGrade = state.Grade
			state.PreferredSemester = state.Semester
			state.Step = 3
			setUserState(firebaseClient, userID, state)
			semesterText := "上學期"
			if semester == 2 {
				semesterText = "下學期"
			}
			return replyMessage(event, bot, fmt.Sprintf("已選擇：%s %d年級%s\n\n✅ 已記憶您的偏好設定，下次查詢將直接使用\n\n請輸入課次（例如：5）：", state.Publisher, state.Grade, semesterText))
		} else {
			return replyMessage(event, bot, "請選擇正確的學期：1（上學期）或 2（下學期）")
		}

	case 3: // 等待課次
		if lesson := parseLesson(userText); lesson > 0 {
			state.Lesson = lesson
			// 更新偏好設定（適用於修改課程模式）
			state.PreferredPublisher = state.Publisher
			state.PreferredGrade = state.Grade
			state.PreferredSemester = state.Semester
			state.Step = 4
			setUserState(firebaseClient, userID, state)
			semesterText := "上學期"
			if state.Semester == 2 {
				semesterText = "下學期"
			}
			return replyMessage(event, bot, fmt.Sprintf("已設定：%s %d年級%s第%d課\n\n✅ 已更新偏好設定\n\n請輸入要查詢的字詞（例如：我好喜歡吃飯配菜）：", state.Publisher, state.Grade, semesterText, lesson))
		} else {
			return replyMessage(event, bot, "請輸入正確的課次數字")
		}

	case 4: // 等待查詢字詞
		if isChineseCharacter(userText) {
			return performCumulativeQuery(event, userText, bot, firebaseClient, userID, state)
		} else {
			return replyMessage(event, bot, "請輸入中文字詞進行查詢")
		}

	default:
		clearUserState(firebaseClient, userID)
		return replyMessage(event, bot, "查詢過程出現錯誤，請重新開始")
	}
}

// 解析年級
func parseGrade(text string) int {
	switch text {
	case "1", "一":
		return 1
	case "2", "二":
		return 2
	case "3", "三":
		return 3
	case "4", "四":
		return 4
	case "5", "五":
		return 5
	case "6", "六":
		return 6
	default:
		return 0
	}
}

// 解析學期
func parseSemester(text string) int {
	switch text {
	case "1", "上", "上學期":
		return 1
	case "2", "下", "下學期":
		return 2
	default:
		return 0
	}
}

// 解析課次
func parseLesson(text string) int {
	// 簡單的數字解析
	var lesson int
	_, err := fmt.Sscanf(text, "%d", &lesson)
	if err != nil {
		return 0
	}
	return lesson
}

// 執行累積字詞查詢
func performCumulativeQuery(event *linebot.Event, queryText string, bot *linebot.Client, firebaseClient *config.FirebaseClient, userID string, state *UserState) error {
	// 分解查詢字詞為單個字符
	queryChars := []rune(queryText)
	learnedChars := []string{}
	notLearnedChars := []string{}

	// 獲取累積生字列表
	cumulativeChars, err := getCumulativeCharacters(firebaseClient, state.Publisher, state.Grade, state.Semester, state.Lesson)
	if err != nil {
		log.Printf("Error getting cumulative characters: %v", err)
		return replyMessage(event, bot, "查詢過程中發生錯誤，請稍後再試")
	}

	// 檢查每個字符是否已學過
	for _, char := range queryChars {
		charStr := string(char)
		if contains(cumulativeChars, charStr) {
			learnedChars = append(learnedChars, charStr)
		} else {
			notLearnedChars = append(notLearnedChars, charStr)
		}
	}

	// 構建回覆訊息（移除範圍和查詢字詞顯示）
	responseText := "📊 累積字詞查詢結果\n\n"

	if len(learnedChars) > 0 {
		responseText += fmt.Sprintf("✅ 已學過：%s\n", strings.Join(learnedChars, ""))
	}

	if len(notLearnedChars) > 0 {
		responseText += fmt.Sprintf("❌ 尚未學過：%s\n", strings.Join(notLearnedChars, ""))
	}

	responseText += fmt.Sprintf("\n📈 統計：已學 %d/%d 字", len(learnedChars), len(queryChars))
	responseText += "\n\n💡 輸入新的字詞繼續查詢，或輸入「退出」結束查詢"

	return replyMessage(event, bot, responseText)
}

// 重設用戶偏好設定
func resetUserPreferences(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient, userID string) error {
	state := getUserState(firebaseClient, userID)
	if state.PreferredPublisher != "" || state.PreferredGrade > 0 || state.PreferredSemester > 0 {
		// 清除偏好設定但保留其他狀態
		state.PreferredPublisher = ""
		state.PreferredGrade = 0
		state.PreferredSemester = 0
		setUserState(firebaseClient, userID, state)
		return replyMessage(event, bot, "✅ 已清除您的偏好設定記憶\n\n下次查詢時將重新選擇出版社、年級和學期")
	} else {
		return replyMessage(event, bot, "目前沒有已記憶的偏好設定")
	}
}

func handleHelp(event *linebot.Event, bot *linebot.Client) error {
	helpText := `🎓 中文學習小幫手使用說明

📝 功能介紹：
• 輸入「查詢累積字詞」開始累積字詞查詢
• 輸入「印字帖」前往印字帖網站下載練習字帖
• 輸入「平板學寫字」前往平板練字頁面
• 輸入「重設偏好」清除記憶的版本/年級/學期設定
• 輸入「退出」退出當前模式

💡 使用範例：
1. 輸入「查詢累積字詞」
2. 首次使用：選擇出版社（康軒/南一/翰林）、年級（1-6年級）、學期（上學期/下學期）
3. 再次使用：系統會記住您的設定，直接輸入課次
4. 輸入課次（例如：5）
5. 輸入要查詢的字詞（例如：我好喜歡吃飯配菜）

📝 印字帖功能：
• 輸入「印字帖」前往 hanziplay.com 下載練習字帖
• 如果已設定偏好版本，會自動帶入出版社/年級/學期參數

✍️ 平板學寫字功能：
• 輸入「平板學寫字」前往平板練字頁面
• 可在平板上直接練習寫字，提供即時筆劃指導

🔧 其他指令：
• 「重設偏好」- 清除記憶的出版社/年級/學期設定
• 「退出」- 退出當前查詢模式

❓ 需要協助請輸入「幫助」`

	return replyMessage(event, bot, helpText)
}

// 獲取累積生字列表（參考demo.js的邏輯）
func getCumulativeCharacters(firebaseClient *config.FirebaseClient, publisher string, grade int, semester int, lesson int) ([]string, error) {
	allCharacters := make(map[string]bool)

	// 查詢所有符合條件的課程
	lessonsRef := firebaseClient.Firestore.Collection("lessons")
	query := lessonsRef.Where("publisher", "==", publisher)

	docs, err := query.Documents(firebaseClient.Ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query lessons: %v", err)
	}

	// 遍歷所有課程，找出符合範圍的課程
	for _, doc := range docs {
		data := doc.Data()
		lessonGrade, ok := data["grade"].(int64)
		if !ok {
			continue
		}
		lessonSemester, ok := data["semester"].(int64)
		if !ok {
			continue
		}
		lessonNumber, ok := data["lesson"].(int64)
		if !ok {
			continue
		}

		// 判斷是否在累積範圍內
		isInRange := (int(lessonGrade) < grade) ||
			(int(lessonGrade) == grade && int(lessonSemester) < semester) ||
			(int(lessonGrade) == grade && int(lessonSemester) == semester && int(lessonNumber) <= lesson)

		if isInRange {
			// 提取課程中的字符
			if charactersData, exists := data["characters"]; exists {
				switch chars := charactersData.(type) {
				case []interface{}:
					for _, charInterface := range chars {
						if charMap, ok := charInterface.(map[string]interface{}); ok {
							if character, exists := charMap["character"]; exists {
								if charStr, ok := character.(string); ok {
									allCharacters[charStr] = true
								}
							}
						} else if charStr, ok := charInterface.(string); ok {
							// 如果直接是字符串
							allCharacters[charStr] = true
						}
					}
				}
			}
		}
	}

	// 轉換為字符串切片
	result := make([]string, 0, len(allCharacters))
	for char := range allCharacters {
		result = append(result, char)
	}

	return result, nil
}

// 檢查字符串切片是否包含指定字符
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func handleUnknownMessage(event *linebot.Event, bot *linebot.Client) error {
	return replyMessage(event, bot, "抱歉，我不太理解您的意思。請輸入「幫助」查看使用說明，或輸入「查詢累積字詞」開始查詢。")
}

func handleFollow(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient) error {
	welcomeText := `🎉 歡迎使用中文學習小幫手！

我可以幫助您：
📚 查詢累積字詞

請輸入「查詢累積字詞」開始使用，或輸入「幫助」查看詳細說明！`

	return replyMessage(event, bot, welcomeText)
}

func handleUnfollow(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient) error {
	// 清除用戶狀態
	userID := event.Source.UserID
	clearUserState(firebaseClient, userID)
	return nil
}

// 輔助函數
func isChineseCharacter(text string) bool {
	// 檢查是否包含中文字符
	matched, _ := regexp.MatchString(`[\p{Han}]+`, text)
	return matched && len([]rune(text)) <= 10 // 限制長度避免長句子被誤判
}

func replyMessage(event *linebot.Event, bot *linebot.Client, text string) error {
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do()
	return err
}

func replyMessageWithQuickReply(event *linebot.Event, bot *linebot.Client, text string, quickReply *linebot.QuickReplyItems) error {
	message := linebot.NewTextMessage(text).WithQuickReplies(quickReply)
	_, err := bot.ReplyMessage(event.ReplyToken, message).Do()
	return err
}

// 處理印字帖功能
func handlePrintWorksheet(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient, userID string) error {
	state := getUserState(firebaseClient, userID)
	
	// 建立基本 URL
	baseURL := "https://hanziplay.com/practice-sheet"
	
	// 如果有用戶偏好設定，加入查詢參數
	if state.PreferredPublisher != "" && state.PreferredGrade > 0 && state.PreferredSemester > 0 {
		// 轉換出版社名稱為英文參數
		var publisher string
		switch state.PreferredPublisher {
		case "康軒":
			publisher = "kang-hsuan"
		case "南一":
			publisher = "nan-i"
		case "翰林":
			publisher = "han-lin"
		default:
			publisher = ""
		}
		
		if publisher != "" {
			baseURL += fmt.Sprintf("?publisher=%s&grade=%d&semester=%d", 
				publisher, state.PreferredGrade, state.PreferredSemester)
		}
		
		semesterText := "上學期"
		if state.PreferredSemester == 2 {
			semesterText = "下學期"
		}
		
		responseText := fmt.Sprintf("📝 印字帖功能\n\n✅ 已使用您的偏好設定：\n📚 %s %d年級%s\n\n🔗 請點擊連結前往印字帖頁面：\n%s\n\n💡 您可以在網站上選擇要印製的字詞並下載字帖", 
			state.PreferredPublisher, state.PreferredGrade, semesterText, baseURL)
		
		return replyMessage(event, bot, responseText)
	} else {
		// 沒有偏好設定，直接提供基本連結
		responseText := fmt.Sprintf("📝 印字帖功能\n\n🔗 請點擊連結前往印字帖頁面：\n%s\n\n💡 建議您先使用「查詢累積字詞」功能設定版本年級學期，下次使用印字帖功能時會自動帶入您的設定", baseURL)
		
		return replyMessage(event, bot, responseText)
	}
}

// 處理平板學寫字功能
func handleTabletPractice(event *linebot.Event, bot *linebot.Client) error {
	url := "https://hanziplay.com/characters/practice"
	responseText := fmt.Sprintf("✍️ 平板學寫字\n\n🔗 請點擊連結前往平板練字頁面：\n%s\n\n💡 您可以在平板上直接練習寫字，提供即時筆劃指導", url)
	
	return replyMessage(event, bot, responseText)
}

// 顯示用戶設定
func showUserSettings(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient, userID string) error {
	state := getUserState(firebaseClient, userID)
	
	var response string
	
	if state.PreferredPublisher == "" && state.PreferredGrade == 0 && state.PreferredSemester == 0 {
		response = "📋 使用者課程設定\n\n❌ 尚未設定任何偏好\n\n請先使用『查詢累積字詞』功能來設定您的偏好設定。"
	} else {
		semesterText := "上學期"
		if state.PreferredSemester == 2 {
			semesterText = "下學期"
		}
		
		response = fmt.Sprintf("📋 使用者課程設定\n\n✅ 已記憶的偏好設定：\n📚 出版社：%s\n🎓 年級：%d年級\n📅 學期：%s\n\n", 
			state.PreferredPublisher, 
			state.PreferredGrade, 
			semesterText)
		
		// 如果用戶當前在累積查詢模式中，顯示當前狀態
		if state.Mode == "cumulative_query" {
			currentSemesterText := "上學期"
			if state.Semester == 2 {
				currentSemesterText = "下學期"
			}
			
			response += fmt.Sprintf("🔄 當前查詢狀態：\n📚 出版社：%s\n🎓 年級：%d年級\n📅 學期：%s\n", 
				state.Publisher, 
				state.Grade, 
				currentSemesterText)
			
			if state.Lesson > 0 {
				response += fmt.Sprintf("📖 課次：第%d課\n", state.Lesson)
			}
			
			switch state.Step {
			case -1:
				response += "⏳ 等待選擇操作"
			case 0:
				response += "⏳ 等待選擇出版社"
			case 1:
				response += "⏳ 等待選擇年級"
			case 2:
				response += "⏳ 等待選擇學期"
			case 3:
				response += "⏳ 等待輸入課次"
			case 4:
				response += "⏳ 等待輸入查詢字詞"
			}
		} else {
			response += "💡 提示：輸入『查詢累積字詞』開始查詢"
		}
	}
	
	return replyMessage(event, bot, response)
}