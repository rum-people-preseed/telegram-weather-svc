package main

import (
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/repositories/temporal_storage"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
	"go.uber.org/zap"
)

func main() {

	bot := models.NewBot()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	memoryStorage := temporal_storage.NewMemoryStorage()
	messagesController := controller.NewMessageHandler(bot, logger.Sugar(), memoryStorage)
	weatherService := services.WeatherProvider{}

	startUsecase := usecases.StartUsecase{}
	helpUsecase := usecases.HelpUsecase{}
	predictUsecase := usecases.PredictUsecase{&weatherService}
	updateLocationUsecase := usecases.UpdateLocationUsecase{}

	messagesController.RegisterUsecase(&startUsecase, "/start")
	messagesController.RegisterUsecase(&helpUsecase, "/help")
	messagesController.RegisterUsecase(&predictUsecase, "/predict")
	messagesController.RegisterUsecase(&updateLocationUsecase, "/update_location")

	updates := bot.SetUpUpdates()
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		err := messagesController.AcceptNewUpdate(&update)
		if err != nil {
			logger.Sugar().Warnf("Error while handling message %v", err)
		}
	}
}
