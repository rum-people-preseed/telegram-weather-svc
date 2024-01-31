package controller

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func InvalidMessage(chatID int64) *tgbotapi.MessageConfig {
	mes := tgbotapi.NewMessage(chatID, "Invalid message/command. Please see /help")
	return &mes
}

func InvalidCallbackData(chatID int64) *tgbotapi.MessageConfig {
	mes := tgbotapi.NewMessage(chatID, "Please click button in order to make choice!")
	return &mes
}
