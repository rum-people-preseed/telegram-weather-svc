package main

import (
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_reader"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	bot := models.NewBot(logger.Sugar())

	messagesController := controller.NewMessageHandler(bot, logger.Sugar())
	geoService := services.NewGeoNameService(logger.Sugar())
	helpFactory := usecases.HelpUsecaseFactory{}
	startFactory := usecases.StartUsecaseFactory{}
	weatherService := services.WeatherProvider{}
	predictFactory := usecases.PredictUsecaseFactory{WeatherService: &weatherService, GeoService: &geoService}

	messagesController.RegisterUsecaseFactory(&helpFactory)
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
			_ = bot.SendMessage(controller.InvalidMessage(message_reader.GetChatId(&update)))
		}
	}
}
