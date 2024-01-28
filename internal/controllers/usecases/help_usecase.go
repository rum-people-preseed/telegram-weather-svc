package usecases

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type HelpUsecase struct {
}

func (u *HelpUsecase) Handle(message *tgbotapi.Message, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	mes := tgbotapi.NewMessage(message.Chat.ID,
		"This bot can help you to get information about weather conditions at any region on any date\n\n"+
			"List of available commands:\n"+
			"/help - show this message\n"+
			"/start - show introduction message\n"+
			"/update_location - update default location")
	return &mes, Finished
}
