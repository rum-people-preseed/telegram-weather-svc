package sequences

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	c "github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_reader"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models/location_chooser"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
)

const (
	InitialState        string = "initial_state"
	GettingCountryState        = "getting_country"
	GettingCityState           = "getting_city"
)

type GetLocationSequence struct {
	geoService  services.GeoService
	state       string
	countryName string
	cityName    string
}

func CreateGetLocationSequence(geoService services.GeoService) GetLocationSequence {
	return GetLocationSequence{
		state:      InitialState,
		geoService: geoService,
	}
}

func (s *GetLocationSequence) GetCityName() string {
	return s.cityName
}

func (s *GetLocationSequence) GetCountryName() string {
	return s.countryName
}

func (s *GetLocationSequence) Handle(update *tgbotapi.Update) (*tgbotapi.MessageConfig, c.Status) {
	switch s.state {
	case InitialState:
		return s.handleInitialState(message_reader.GetChatId(update))
	case GettingCountryState:
		return s.handleGettingCountry(update.Message)
	case GettingCityState:
		return s.handleGettingCity(update.Message)
	default:
		return c.InvalidMessage(update.Message.Chat.ID), c.Error
	}
}

func (s *GetLocationSequence) handleInitialState(chatID int64) (*tgbotapi.MessageConfig, c.Status) {
	mes := tgbotapi.NewMessage(chatID, location_chooser.ResponseEnterCountry)
	s.state = GettingCountryState

	return &mes, c.Continue
}

func (s *GetLocationSequence) handleGettingCountry(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	s.countryName = message.Text

	err := s.geoService.ValidateCountry(s.countryName)
	if err != nil {
		errMsg := tgbotapi.NewMessage(message.Chat.ID, location_chooser.CountryValidationError)
		return &errMsg, c.Error
	}

	mes := tgbotapi.NewMessage(message.Chat.ID, location_chooser.ResponseEnterCity)
	s.state = GettingCityState

	return &mes, c.Continue
}

func (s *GetLocationSequence) handleGettingCity(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	s.cityName = message.Text

	err := s.geoService.ValidateCity(s.cityName, s.countryName)
	if err != nil {
		errMsg := tgbotapi.NewMessage(message.Chat.ID, location_chooser.CityValidationError)
		return &errMsg, c.Error
	}

	return nil, c.Finished
}
