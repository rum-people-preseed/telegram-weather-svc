package sequences

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	c "github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_constructor"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models/location_chooser"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
)

const (
	InitialState        string = "initial_state"
	GettingCountryState string = "getting_country"
	GettingCityState    string = "getting_city"
)

type GetLocationSequence struct {
	geoService  services.GeoService
	chatID      int64
	state       string
	countryName string
	cityName    string
	coordinates models.Coordinates
}

func CreateGetLocationSequence(chatID int64, geoService services.GeoService) GetLocationSequence {
	return GetLocationSequence{
		state:      InitialState,
		geoService: geoService,
		chatID:     chatID,
	}
}

func (s *GetLocationSequence) GetCityName() string {
	return s.cityName
}

func (s *GetLocationSequence) GetCountryName() string {
	return s.countryName
}

func (s *GetLocationSequence) GetCoordinates() models.Coordinates {
	return s.coordinates
}

func (s *GetLocationSequence) Handle(update *tgbotapi.Update) (*tgbotapi.MessageConfig, c.Status) {
	switch s.state {
	case InitialState:
		return s.handleInitialState()
	case GettingCountryState:
		return s.handleGettingCountry(update.Message)
	case GettingCityState:
		return s.handleGettingCity(update.Message)
	default:
		return c.InvalidMessage(update.Message.Chat.ID), c.Error
	}
}

func (s *GetLocationSequence) handleInitialState() (*tgbotapi.MessageConfig, c.Status) {
	mes := tgbotapi.NewMessage(s.chatID, location_chooser.ResponseEnterCountry)
	s.state = GettingCountryState

	return &mes, c.Continue
}

func (s *GetLocationSequence) handleGettingCountry(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	s.countryName = message.Text

	countryName, err := s.geoService.ValidateCountry(s.countryName)
	if err != nil {
		errMsg := tgbotapi.NewMessage(s.chatID, location_chooser.CountryValidationError)
		return &errMsg, c.Continue
	}

	if countryName != s.countryName {
		mes := message_constructor.MakeTextMessage(s.chatID, fmt.Sprintf(location_chooser.DidYouMean, countryName))
		return &mes, c.Continue
	}

	mes := tgbotapi.NewMessage(s.chatID, location_chooser.ResponseEnterCity)
	s.state = GettingCityState

	return &mes, c.Continue
}

func (s *GetLocationSequence) handleGettingCity(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	s.cityName = message.Text

	cityName, coordinates, err := s.geoService.ValidateCity(s.cityName, s.countryName)
	if err != nil {
		errMsg := tgbotapi.NewMessage(s.chatID, location_chooser.CityValidationError)
		return &errMsg, c.Continue
	}

	if cityName != s.cityName {
		mes := message_constructor.MakeTextMessage(s.chatID, fmt.Sprintf(location_chooser.DidYouMean, cityName))
		return &mes, c.Continue
	}

	s.coordinates = coordinates
	return nil, c.Finished
}
