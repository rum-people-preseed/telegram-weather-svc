package models

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"os"
)

type Bot struct {
	BotAPI *tgbotapi.BotAPI
}

func NewBot() Bot {
	telegramApiToken := os.Getenv("TELEGRAM_API_TOKEN")
	if telegramApiToken == "" {
		zap.String("error", "Failed to load env TELEGRAM_API_TOKEN")
		os.Exit(1)
	}

	apiBot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		zap.String("error", "Failed to bind to API bot with token")
		os.Exit(1)
	}

	apiBot.Debug = true

	zap.String("info", fmt.Sprintf("Authorized on account %s", apiBot.Self.UserName))
	return Bot{BotAPI: apiBot}
}

func (bot *Bot) SetUpUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.BotAPI.GetUpdatesChan(u)
	if err != nil {
		zap.String("error", "Failed to get updates chanel")
		os.Exit(1)
	}

	zap.String("info", "Bot is ready to receive updates from channel")
	return updates
}

func (bot *Bot) SendMessage(msg tgbotapi.MessageConfig) {
	if _, err := bot.BotAPI.Send(msg); err != nil {
		zap.String("error", fmt.Sprintf("Failed to send message %s", msg.Text))
		os.Exit(1)
	}
}
