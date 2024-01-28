package main

import (
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"go.uber.org/zap"
)

func main() {

	bot := models.NewBot()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	messagesController := controller.NewMessageHandler(bot.BotAPI, logger.Sugar())
	startUsecase := usecases.StartUsecase{}

	messagesController.RegisterUsecase(&startUsecase, "/start")

	updates := bot.SetUpUpdates()
	for update := range updates {
		if update.Message == nil {
			continue
		}

		err := messagesController.AcceptNewMessage(update.Message)
		if err != nil {
			logger.Sugar().Warnf("error while handling message %v", err)
		}
	}
}
