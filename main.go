package main

import (
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/repositories/temporal_storage"
	"go.uber.org/zap"
)

func main() {

	bot := models.NewBot()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	memoryStorage := temporal_storage.NewMemoryStorage()
	messagesController := controller.NewMessageHandler(bot.BotAPI, logger.Sugar(), memoryStorage)
	startUsecase := usecases.StartUsecase{}
	helpUsecase := usecases.HelpUsecase{}
	predictUsecase := usecases.PredictUsecase{}
	updateLocationUsecase := usecases.UpdateLocationUsecase{}

	messagesController.RegisterUsecase(&startUsecase, "/start")
	messagesController.RegisterUsecase(&helpUsecase, "/help")
	messagesController.RegisterUsecase(&predictUsecase, "/predict")
	messagesController.RegisterUsecase(&updateLocationUsecase, "/update_location")

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
