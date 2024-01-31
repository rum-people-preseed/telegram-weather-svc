package controller

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func InvalidMessage(chatID int64) *tgbotapi.MessageConfig {
	mes := tgbotapi.NewMessage(chatID, "Internal error")
	return &mes
}
