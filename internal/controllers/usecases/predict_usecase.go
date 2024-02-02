package usecases

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	c "github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/sequences"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_constructor"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_reader"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models/time_chooser"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
)

type PredictUsecaseFactory struct {
	WeatherService services.WeatherService
	GeoService     services.GeoService
}

func (f *PredictUsecaseFactory) Create() c.Usecase {
	statesWithCallbackData := make(map[string]bool)
	statesWithCallbackData[NextDayOrDateResponseState] = true

	return &PredictUsecase{
		locationSequence:       sequences.CreateGetLocationSequence(f.GeoService),
		state:                  LocationSequenceState,
		weatherService:         f.WeatherService,
		statesWithCallbackData: statesWithCallbackData,
	}

}

func (f *PredictUsecaseFactory) Command() string {
	return "/predict"
}

type PredictUsecase struct {
	weatherService         services.WeatherService
	state                  string
	locationSequence       sequences.GetLocationSequence
	statesWithCallbackData map[string]bool
	country                string
	city                   string
	date                   string
}

const (
	LocationSequenceState      = "location_sequence"
	NextDayOrDateResponseState = "next_day_or_date_response"
	EnterDateResponseState     = "enter_date_response"
)

func (u *PredictUsecase) Handle(update *tgbotapi.Update) (*tgbotapi.MessageConfig, c.Status) {

	err := u.CheckCorrectnessOfCallback(update)
	if err != nil {
		return c.InvalidCallbackData(message_reader.GetChatId(update)), c.Continue
	}

	switch u.state {
	case LocationSequenceState:
		return u.handleLocationSequenceState(update)
	case NextDayOrDateResponseState:
		return u.handleNextDayOrDateResponseState(update.CallbackQuery)
	case EnterDateResponseState:
		return u.handleEnterDateResponseState(update.Message)
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
		return mes, c.Continue
	}

	return u.handleNextDayOrDateState(update.Message.Chat.ID)
}
func (u *PredictUsecase) handleNextDayOrDateState(chatID int64) (*tgbotapi.MessageConfig, c.Status) {
	mes := message_constructor.MakeMessageWithButtons(chatID, time_chooser.MainMessage,
		message_constructor.MakeInlineButton(time_chooser.OptionNextDay, time_chooser.OptionNextDayCallbackData),
		message_constructor.MakeInlineButton(time_chooser.OptionEnterDate, time_chooser.OptionEnterDateCallbackData))
	u.state = NextDayOrDateResponseState
	return &mes, c.Continue
}

func (u *PredictUsecase) handleNextDayOrDateResponseState(callbackQuery *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, c.Status) {
	chatID, callbackData := callbackQuery.Message.Chat.ID, callbackQuery.Data

	switch callbackData {
	case time_chooser.OptionNextDayCallbackData:
		return u.RequestWeatherForecast(chatID)
	case time_chooser.OptionEnterDateCallbackData:
		return u.handleEnterDateState(chatID)
	default:
		return c.InvalidMessage(chatID), c.Error
	}
}

func (u *PredictUsecase) handleEnterDateState(chatID int64) (*tgbotapi.MessageConfig, c.Status) {
	mes := message_constructor.MakeTextMessage(chatID, time_chooser.ResponseEnterDate)
	u.state = EnterDateResponseState
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

func (u *PredictUsecase) CheckCorrectnessOfCallback(update *tgbotapi.Update) error {
	var err error
	if update.CallbackQuery == nil && u.statesWithCallbackData[u.state] {
		err = fmt.Errorf("callback was expected from chat %v", message_reader.GetChatId(update))
	}
	return err
}
