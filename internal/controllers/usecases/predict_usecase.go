package usecases

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type PredictUsecase struct {
}

func (u *PredictUsecase) Handle(message *tgbotapi.Message, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	mes := tgbotapi.NewMessage(message.Chat.ID, "Will be done!!!")
	return &mes, Finished
}
