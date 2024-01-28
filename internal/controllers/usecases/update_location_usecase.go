package usecases

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UpdateLocationUsecase struct {
}

const (
	activatedKey = "activated"
	countryKey   = "country"
	cityKey      = "city"
)

func (u *UpdateLocationUsecase) Handle(message *tgbotapi.Message, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	invalidMsg := tgbotapi.NewMessage(message.Chat.ID, "Internal error")

	_, err := usecaseData.Get(activatedKey)
	if err != nil {
		err := usecaseData.Set(activatedKey, "activated")
		if err != nil {
			return &invalidMsg, Error
		}

		mes := tgbotapi.NewMessage(message.Chat.ID, "Enter desired country:\n")
		return &mes, Continue
	}

	country, err := usecaseData.Get(countryKey)

	if err != nil {
		err := usecaseData.Set(countryKey, message.Text)
		if err != nil {
			return &invalidMsg, Error
		}

		mes := tgbotapi.NewMessage(message.Chat.ID, "Enter desired city:\n")
		return &mes, Continue
	}

	err = usecaseData.Set(cityKey, message.Text)
	if err != nil {
		return &invalidMsg, Error
	}

	mes := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("You location has been updated: %v, %v", country, message.Text))

	return &mes, Finished
}
