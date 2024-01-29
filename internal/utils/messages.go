package utils

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetMessageWithButtons(chatID int64, msgText string, buttons ...tgbotapi.InlineKeyboardButton) tgbotapi.MessageConfig {
	msgCfg := tgbotapi.NewMessage(chatID, msgText)
	var dateChooseButtons = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			buttons...,
		),
	)
	msgCfg.ReplyMarkup = dateChooseButtons
	return msgCfg
}

func GetInlineButton(text string, callbackData string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, callbackData)
}

func GetSimpleMessage(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, text)
}

func GetChatId(update *tgbotapi.Update) int64 {
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	return chatID
}

func ExtractCommand(msg *tgbotapi.Message) (string, error) {
	if msg == nil {
		return "", errors.New("Empty message")
	}
	text := msg.Text
	err := errors.New("command not found")
	if text[0] != '/' {
		return "", err
	}
	var command string
	for _, alpha := range text {
		if alpha == ' ' {
			break
		}
		command += string(alpha)
	}
	return command, nil
}
