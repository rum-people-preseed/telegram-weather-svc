package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	telegramApiToken := os.Getenv("TELEGRAM_API_TOKEN")
	if telegramApiToken == "" {
		sugar.Error("Failed to load env TELEGRAM_API_TOKEN")
		os.Exit(1)
	}

	apiBot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		sugar.Error("Failed to bind to API bot with token")
		os.Exit(1)
	}

	apiBot.Debug = true

	sugar.Info("Authorized on account %s", apiBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := apiBot.GetUpdatesChan(u)

	if err != nil {
		sugar.Error("Failed to get updates chanel")
		os.Exit(1)
	}

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		if _, err := apiBot.Send(msg); err != nil {
			sugar.Error("Failed to send message")
			os.Exit(1)
		}
	}
}
