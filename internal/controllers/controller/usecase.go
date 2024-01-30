package controller

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type DataMap interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Del(key string) error
}

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
