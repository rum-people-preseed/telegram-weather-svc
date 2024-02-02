package message_constructor

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func MakeMessageWithButtons(chatID int64, msgText string, buttons ...tgbotapi.InlineKeyboardButton) tgbotapi.MessageConfig {
	msgCfg := tgbotapi.NewMessage(chatID, msgText)
	var dateChooseButtons = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			buttons...,
		),
	)
	msgCfg.ReplyMarkup = dateChooseButtons
	return msgCfg
}

func MakeInlineButton(text string, callbackData string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, callbackData)
}

func MakeTextMessage(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, text)
}
