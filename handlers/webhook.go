package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"chinese-learning-linebot/config"
)

func WebhookHandler(bot *linebot.Client, firebaseClient *config.FirebaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		for _, event := range events {
			if err := handleEvent(event, bot, firebaseClient); err != nil {
				log.Printf("Error handling event: %v", err)
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func handleEvent(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient) error {
	switch event.Type {
	case linebot.EventTypeMessage:
		return handleMessage(event, bot, firebaseClient)
	case linebot.EventTypePostback:
		return handlePostback(event, bot, firebaseClient)
	case linebot.EventTypeFollow:
		return handleFollow(event, bot, firebaseClient)
	case linebot.EventTypeUnfollow:
		return handleUnfollow(event, bot, firebaseClient)
	default:
		log.Printf("Unknown event type: %s", event.Type)
	}
	return nil
}

func handlePostback(event *linebot.Event, bot *linebot.Client, firebaseClient *config.FirebaseClient) error {
	// 目前不處理postback事件
	return nil
}