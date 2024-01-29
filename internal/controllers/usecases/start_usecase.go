package usecases

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type StartUsecase struct {
}

func (u *StartUsecase) Handle(update *tgbotapi.Update, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	message := update.Message
	introText := "Hello, " + message.From.UserName + "!\n" +
		"This bot is made for help you do not get wet in the rain, do not die from the heat or be ready for an abnormal storm. \n\n" +
		"Please, follow /help to get information about all facilities we are provide."
	mes := tgbotapi.NewMessage(message.Chat.ID, introText)
	return &mes, Finished
}
