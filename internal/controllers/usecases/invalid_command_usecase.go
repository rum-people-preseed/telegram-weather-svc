package usecases

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_constructor"
)

type InvalidCommandUsecaseFactory struct {
}

func (f *InvalidCommandUsecaseFactory) Create(chatID int64) controller.Usecase {
	return &InvalidCommandUsecase{
		chatID: chatID,
	}
}

func (f *InvalidCommandUsecaseFactory) Command() string {
	return "/invalid_command"
}

type InvalidCommandUsecase struct {
	chatID int64
}

func (u *InvalidCommandUsecase) Handle(_ *tgbotapi.Update) (tgbotapi.Chattable, controller.Status) {
	mes := message_constructor.MakeTextMessage(u.chatID,
		"Please refer to /help to use correct command")
	return &mes, controller.Finished
}
