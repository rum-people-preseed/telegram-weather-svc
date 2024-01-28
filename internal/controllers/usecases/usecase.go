package usecases

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Usecase interface {
	Handle(message *tgbotapi.Message)
}
