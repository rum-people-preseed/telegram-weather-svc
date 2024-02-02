package usecases

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
)

type HelpUsecaseFactory struct {
}

func (f *HelpUsecaseFactory) Create() controller.Usecase {
	return &HelpUsecase{}
}

func (f *HelpUsecaseFactory) Command() string {
	return "/help"
}

type HelpUsecase struct {
}

func (u *HelpUsecase) Handle(update *tgbotapi.Update) (tgbotapi.Chattable, controller.Status) {
	message := update.Message
	mes := tgbotapi.NewMessage(message.Chat.ID,
		"This bot can help you to get information about weather conditions at any region on any date\n\n"+
			"List of available commands:\n"+
			"/help - show this message\n"+
			"/start - show introduction message\n"+
			"/predict - get weather forecast \n")
	return &mes, controller.Finished
}
