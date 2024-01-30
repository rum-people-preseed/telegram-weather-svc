package sequences

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	c "github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
)

const (
	InitialState        string = "initial_state"
	GettingCountryState        = "getting_country"
	GettingCityState           = "getting_city"
)

type GetLocationSequence struct {
	state       string
	countryName string
	cityName    string
}

func CreateGetLocationSequence() GetLocationSequence {
	return GetLocationSequence{
		state: InitialState,
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
		return s.handleInitialState(extractChatId(update))
	case GettingCountryState:
		return s.handleGettingCountry(update.Message)
	case GettingCityState:
		return s.handleGettingCity(update.Message)
	default:
		return c.InvalidMessage(update.Message.Chat.ID), c.Error
	}
}

func (s *GetLocationSequence) handleInitialState(chatID int64) (*tgbotapi.MessageConfig, c.Status) {
	mes := tgbotapi.NewMessage(chatID, "Enter desired country:\n")
	s.state = GettingCountryState

	return &mes, c.Continue
}

func (s *GetLocationSequence) handleGettingCountry(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	// todo validate country
	s.countryName = message.Text
	mes := tgbotapi.NewMessage(message.Chat.ID, "Enter desired city:\n")
	s.state = GettingCityState

	return &mes, c.Continue
}

func (s *GetLocationSequence) handleGettingCity(message *tgbotapi.Message) (*tgbotapi.MessageConfig, c.Status) {
	// todo validate city
	s.cityName = message.Text

	return c.InvalidMessage(0), c.Finished
}

func extractChatId(update *tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}

	return 0
}
