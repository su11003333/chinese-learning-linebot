package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"chinese-learning-linebot/config"
	"chinese-learning-linebot/handlers"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化 Firebase
	ctx := context.Background()
	firebaseClient, err := config.InitFirebase(ctx)
	if err != nil {
		log.Printf("Warning: Failed to initialize Firebase: %v", err)
		log.Println("Running without Firebase functionality")
		firebaseClient = nil
	} else {
		defer firebaseClient.Close()
	}

	// 初始化 LINE Bot
	bot, err := config.InitLineBot()
	if err != nil {
		log.Printf("Warning: Failed to initialize LINE Bot: %v", err)
		log.Println("Running without LINE Bot functionality")
		bot = nil
	}

	// 初始化 Gin 路由
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 健康檢查端點
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// LINE Bot Webhook 端點
	r.POST("/webhook", handlers.WebhookHandler(bot, firebaseClient))

	// 啟動服務器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
