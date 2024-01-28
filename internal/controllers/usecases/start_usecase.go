package usecases

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type StartUsecase struct {
}

func (u *StartUsecase) Handle(message *tgbotapi.Message, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	mes := tgbotapi.NewMessage(message.Chat.ID, "Hello dear citizen!")
	return &mes, Finished
}
