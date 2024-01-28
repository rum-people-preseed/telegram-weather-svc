package main

import (
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"go.uber.org/zap"
	"os"
)

func main() {

	bot := models.NewBot()
	updates := bot.SetUpUpdates()
	messagesController := controller.NewMessageHandler(bot.BotAPI)
	startUsecase := usecases.StartUsecase{}
	messagesController.RegisterUsecase(&startUsecase, "/start")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		err := messagesController.AcceptNewMessage(update.Message)
		if err != nil {
			zap.String("error", "Failed to handle message")
			os.Exit(1)
		}
	}
}
