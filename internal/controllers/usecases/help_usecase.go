package usecases

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
)

type HelpUsecaseFactory struct {
}

func (f *HelpUsecaseFactory) Create(chatID int64) controller.Usecase {
	return &HelpUsecase{
		chatID: chatID,
	}
}

func (f *HelpUsecaseFactory) Command() string {
	return "/help"
}

type HelpUsecase struct {
	chatID int64
}

func (u *HelpUsecase) Handle(_ *tgbotapi.Update) (tgbotapi.Chattable, controller.Status) {
	mes := tgbotapi.NewMessage(u.chatID,
		"This bot can help you to get information about weather conditions at any region on any date\n\n"+
			"List of available commands:\n"+
			"/help - show this message\n"+
			"/start - show introduction message\n"+
			"/predict - get weather forecast \n")
	return &mes, controller.Finished
}
