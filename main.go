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

	//memoryStorage := temporal_storage.NewMemoryStorage()
	messagesController := controller.NewMessageHandler(bot, logger.Sugar())
	//weatherService := services.WeatherProvider{}
	helpFactory := usecases.HelpUsecaseFactory{}
	updateLocationFactory := usecases.UpdateLocationUsecaseFactory{}
	startFactory := usecases.StartUsecaseFactory{}
	predictFactory := usecases.PredictUsecaseFactory{}

	messagesController.RegisterUsecaseFactory(&helpFactory)
	messagesController.RegisterUsecaseFactory(&updateLocationFactory)
	messagesController.RegisterUsecaseFactory(&startFactory)
	messagesController.RegisterUsecaseFactory(&predictFactory)

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
