package usecases

import (
	"errors"
	"time"

	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/date_parser"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	c "github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/sequences"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_constructor"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models/time_chooser"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
)

type PredictUsecaseFactory struct {
	WeatherService services.WeatherService
	GeoService     services.GeoService
	DateParser     date_parser.DateParser
}

func (f *PredictUsecaseFactory) Create(chatID int64) c.Usecase {
	statesWithCallbackData := make(map[string]string)
	statesWithCallbackData[NextDayOrDateResponseState] = CallbackQuery
	statesWithCallbackData[EnterDateResponseState] = Message
	statesWithCallbackData[LocationSequenceState] = Message

	return &PredictUsecase{
		chatID:                 chatID,
		locationSequence:       sequences.CreateGetLocationSequence(chatID, f.GeoService),
		state:                  LocationSequenceState,
		weatherService:         f.WeatherService,
		statesWithCallbackData: statesWithCallbackData,
		dateParser:             f.DateParser,
	}

}

func (f *PredictUsecaseFactory) Command() string {
	return "/predict"
}

type PredictUsecase struct {
	chatID                 int64
	weatherService         services.WeatherService
	state                  string
	locationSequence       sequences.GetLocationSequence
	dateParser             date_parser.DateParser
	statesWithCallbackData map[string]string
	date                   time.Time
}

const (
	LocationSequenceState      = "location_sequence"
	NextDayOrDateResponseState = "next_day_or_date_response"
	EnterDateResponseState     = "enter_date_response"
	CallbackQuery              = "CallbackQuery"
	Message                    = "Message"

	TypeMessage = "Please type required message"
	PressButton = "Please press required button"
)

func (u *PredictUsecase) Handle(update *tgbotapi.Update) (tgbotapi.Chattable, c.Status) {

	mes, err := u.CheckCorrectnessOfCallback(update)
	if err != nil {
		return mes, c.Continue
	}

	switch u.state {
	case LocationSequenceState:
		return u.handleLocationSequenceState(update)
	case NextDayOrDateResponseState:
		return u.handleNextDayOrDateResponseState(update.CallbackQuery)
	case EnterDateResponseState:
		return u.handleEnterDateResponseState(update.Message)
	default:
		return c.InvalidMessage(u.chatID), c.Error
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

	return u.handleNextDayOrDateState()
}
func (u *PredictUsecase) handleNextDayOrDateState() (*tgbotapi.MessageConfig, c.Status) {
	mes := message_constructor.MakeMessageWithButtons(u.chatID, time_chooser.MainMessage,
		message_constructor.MakeInlineButton(time_chooser.OptionNextDay, time_chooser.OptionNextDayCallbackData),
		message_constructor.MakeInlineButton(time_chooser.OptionEnterDate, time_chooser.OptionEnterDateCallbackData))
	u.state = NextDayOrDateResponseState
	return &mes, c.Continue
}

func (u *PredictUsecase) handleNextDayOrDateResponseState(callbackQuery *tgbotapi.CallbackQuery) (tgbotapi.Chattable, c.Status) {
	switch callbackQuery.Data {
	case time_chooser.OptionNextDayCallbackData:
		u.date = time.Now().AddDate(0, 0, 1)
		return u.RequestWeatherForecast()
	case time_chooser.OptionEnterDateCallbackData:
		return u.handleEnterDateState()
	default:
		return message_constructor.MakeTextMessage(u.chatID, PressButton), c.Continue
	}
}

func (u *PredictUsecase) handleEnterDateState() (tgbotapi.Chattable, c.Status) {
	mes := message_constructor.MakeTextMessage(u.chatID, time_chooser.ResponseEnterDate)
	u.state = EnterDateResponseState
	return &mes, c.Continue
}

func (u *PredictUsecase) handleEnterDateResponseState(message *tgbotapi.Message) (tgbotapi.Chattable, c.Status) {
	date, err := u.dateParser.ParseDateString(message.Text)
	if err != nil {
		mes := message_constructor.MakeTextMessage(u.chatID, time_chooser.DateValidationError)
		return &mes, c.Continue
	}
	u.date = date
	return u.RequestWeatherForecast()
}

func (u *PredictUsecase) RequestWeatherForecast() (tgbotapi.Chattable, c.Status) {
	city, country := u.locationSequence.GetCityName(), u.locationSequence.GetCountryName()
	weatherData := models.NewWeatherData(city, country, u.locationSequence.GetCoordinates(), u.date)
	responseForecast, err := u.weatherService.GetWeatherPredictionMessage(weatherData, u.chatID)
	if err != nil {
		return c.InvalidMessage(u.chatID), c.Error
	}

	return responseForecast, c.Finished
}

func (u *PredictUsecase) CheckCorrectnessOfCallback(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var err error = nil
	var mes tgbotapi.MessageConfig
	if u.statesWithCallbackData[u.state] == Message && update.Message == nil {
		err = errors.New("expected Message type")
		mes = message_constructor.MakeTextMessage(u.chatID, TypeMessage)
	} else if u.statesWithCallbackData[u.state] == CallbackQuery && update.CallbackQuery == nil {
		err = errors.New("expected CallbackQuery type")
		mes = message_constructor.MakeTextMessage(u.chatID, PressButton)
	}

	return mes, err
}
