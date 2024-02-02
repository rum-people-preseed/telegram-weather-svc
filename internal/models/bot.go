package models

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Bot struct {
	BotAPI *tgbotapi.BotAPI
	log    Logger
}

func NewBot(logger Logger) *Bot {
	telegramApiToken := os.Getenv("TELEGRAM_API_TOKEN")
	if telegramApiToken == "" {
		logger.Errorf("Failed to load env TELEGRAM_API_TOKEN")
		os.Exit(1)
	}

	apiBot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		logger.Errorf("Failed to bind to API bot with token")
		os.Exit(1)
	}

	apiBot.Debug = false
	logger.Infof("Authorized on account %s", apiBot.Self.UserName)
	return &Bot{BotAPI: apiBot, log: logger}
}

func (bot *Bot) SetUpUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.BotAPI.GetUpdatesChan(u)
	if err != nil {
		bot.log.Errorf("Failed to get updates chanel")
		os.Exit(1)
	}

	bot.log.Infof("Bot is ready to receive updates from channel")
	return updates
}

func (bot *Bot) SendMessage(msg *tgbotapi.MessageConfig) error {
	_, err := bot.BotAPI.Send(msg)
	if err != nil {
		bot.log.Errorf("Failed to send message to chat with id %s", msg.ChatID)
		return err
	}

	return nil
}
