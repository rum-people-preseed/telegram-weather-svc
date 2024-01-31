package controller

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Status int64

const (
	Continue Status = 0
	Finished Status = 1
	Error    Status = 2
)

type Usecase interface {
	Handle(update *tgbotapi.Update) (*tgbotapi.MessageConfig, Status)
}

type UsecaseFactory interface {
	Create() Usecase
	Command() string
}
