package main

import (
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/date_parser"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_reader"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
)

func main() {
	logger := models.GetNewLogger()
	bot := models.NewBot(logger)

	messagesController := controller.NewMessageHandler(bot, logger)
	geoService := services.NewGeoNameService(logger)
	helpFactory := usecases.HelpUsecaseFactory{}
	startFactory := usecases.StartUsecaseFactory{}
	weatherService := services.WeatherProvider{}
	dateParser := date_parser.MultiformatDateParser{}
	predictFactory := usecases.PredictUsecaseFactory{WeatherService: &weatherService, GeoService: &geoService, DateParser: &dateParser}

	_ = messagesController.RegisterUsecaseFactory(&helpFactory)
	_ = messagesController.RegisterUsecaseFactory(&startFactory)
	_ = messagesController.RegisterUsecaseFactory(&predictFactory)

	updates := bot.SetUpUpdates()
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		err := messagesController.AcceptNewUpdate(&update)
		if err != nil {
			logger.Warnf("Error while handling message %v", err)
			_ = bot.SendMessage(controller.InvalidMessage(message_reader.GetChatId(&update)))
		}
	}
}
