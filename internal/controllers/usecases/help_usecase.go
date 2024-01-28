package usecases

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type HelpUsecase struct {
}

func (u *HelpUsecase) Handle(message *tgbotapi.Message, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	mes := tgbotapi.NewMessage(message.Chat.ID, "Our possible commands: ")
	return &mes, Finished
}
