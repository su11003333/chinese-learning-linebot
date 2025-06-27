package config

import (
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func InitLineBot() (*linebot.Client, error) {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	if channelSecret == "" || channelToken == "" {
		return nil, fmt.Errorf("LINE_CHANNEL_SECRET and LINE_CHANNEL_ACCESS_TOKEN must be set")
	}

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}

	return bot, nil
}
