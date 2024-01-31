package usecases

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	c "github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/sequences"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_constructor"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
)

type PredictUsecaseFactory struct {
	WeatherService services.WeatherService
}

func (f *PredictUsecaseFactory) Create() c.Usecase {
	return &PredictUsecase{
		locationSequence: sequences.CreateGetLocationSequence(),
		state:            InitialState,
		weatherService:   f.WeatherService,
	}
}

func (f *PredictUsecaseFactory) Command() string {
	return "/predict"
}

type PredictUsecase struct {
	weatherService   services.WeatherService
	state            string
	locationSequence sequences.GetLocationSequence
	country          string
	city             string
	date             string
}

const (
	InitialState               string = "initial_state"
	LocationResponse                  = "location_response"
	LocationSequenceState             = "location_sequence"
	NextDayOrDateResponseState        = "next_day_or_date_response"
	EnterDateResponse                 = "enter_date_response"
)

func (u *PredictUsecase) Handle(update *tgbotapi.Update) (*tgbotapi.MessageConfig, c.Status) {
	switch u.state {
	case InitialState:
		return u.handleInitialState(update.Message)
	case LocationResponse:
		return u.handleLocationResponseState(update)
	case LocationSequenceState:
		return u.handleLocationSequenceState(update)
	case NextDayOrDateResponseState:
		return u.handleNextDayOrDateResponseState(update.CallbackQuery)
	case EnterDateResponse:
		return u.handleEnterDateResponseState(update.Message)
	default:
		return c.InvalidMessage(update.Message.Chat.ID), c.Error
	}
}

func (u *PredictUsecase) handleInitialState(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	mes := message_constructor.MakeMessageWithButtons(message.Chat.ID,
		"Do you wanna receive prediction based on your current location?",
		message_constructor.MakeInlineButton("Yes", "Please use current location"),
		message_constructor.MakeInlineButton("No", "I wanna choose myself"))

	u.state = LocationResponse
	return &mes, c.Continue
}

func (u *PredictUsecase) handleLocationResponseState(update *tgbotapi.Update) (*tgbotapi.MessageConfig, c.Status) {
	chatID, callbackData := update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data

	switch callbackData {
	case "Please use current location":
		return u.handleNextDayOrDateState(chatID)
	case "I wanna choose myself":
		u.state = LocationSequenceState
		return u.locationSequence.Handle(update)
	default:
		return c.InvalidMessage(update.Message.Chat.ID), c.Error
	}
}

func (u *PredictUsecase) handleLocationSequenceState(update *tgbotapi.Update) (*tgbotapi.MessageConfig, c.Status) {
	mes, err := u.locationSequence.Handle(update)

	if err == c.Continue {
		return mes, err
	}
	if err == c.Error {
		// todo handle error state
	}

	return u.handleNextDayOrDateState(update.Message.Chat.ID)
}
func (u *PredictUsecase) handleNextDayOrDateState(chatID int64) (*tgbotapi.MessageConfig, c.Status) {
	mes := message_constructor.MakeMessageWithButtons(chatID, "Please choose the day you want to receive the weather prediction:",
		message_constructor.MakeInlineButton("Next day", "Next day"),
		message_constructor.MakeInlineButton("Enter a date", "Enter a date"))
	u.state = NextDayOrDateResponseState
	return &mes, c.Continue
}

func (u *PredictUsecase) handleNextDayOrDateResponseState(callbackQuery *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, c.Status) {
	chatID, callbackData := callbackQuery.Message.Chat.ID, callbackQuery.Data

	switch callbackData {
	case "Next day":
		return u.RequestWeatherForecast(chatID)
	case "Enter a date":
		return u.handleEnterDateState(chatID)
	default:
		return c.InvalidMessage(chatID), c.Error
	}
}

func (u *PredictUsecase) handleEnterDateState(chatID int64) (*tgbotapi.MessageConfig, c.Status) {
	mes := message_constructor.MakeTextMessage(chatID, "Enter desired day:\n")
	u.state = EnterDateResponse
	return &mes, c.Continue
}

func (u *PredictUsecase) handleEnterDateResponseState(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	// todo validate and translate message
	u.date = message.Text
	return u.RequestWeatherForecast(message.Chat.ID)
}

func (u *PredictUsecase) RequestWeatherForecast(chatID int64) (*tgbotapi.MessageConfig, c.Status) {
	weatherData := models.WeatherData{Country: u.country, City: u.city}
	responseForecast, err := u.weatherService.GetWeather(weatherData)
	if err != nil {
		return c.InvalidMessage(chatID), c.Error
	}
	mes := message_constructor.MakeTextMessage(chatID, responseForecast)
	return &mes, c.Finished
}
