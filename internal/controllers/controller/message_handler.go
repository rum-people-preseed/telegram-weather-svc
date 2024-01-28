package controller

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/repositories/temporal_storage"
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

const activeCommandKey = "active_command"

type MessageHandler struct {
	registeredCallbacks map[string]usecases.Usecase
	usecasesData        temporal_storage.TemporalStorage
	log                 Logger
	bot                 *tgbotapi.BotAPI
}

func NewMessageHandler(bot *tgbotapi.BotAPI, log Logger, usecasesData temporal_storage.TemporalStorage) *MessageHandler {
	return &MessageHandler{
		bot:                 bot,
		log:                 log,
		usecasesData:        usecasesData,
		registeredCallbacks: make(map[string]usecases.Usecase),
	}
}

func (h *MessageHandler) RegisterUsecase(usecase usecases.Usecase, command string) error {
	_, ok := h.registeredCallbacks[command]
	if ok {
		return errors.New(fmt.Sprintf("Callback for command %v is already registered", command))
	}

	h.registeredCallbacks[command] = usecase
	return nil
}

func extractCommand(text string) (string, error) {
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

func (h *MessageHandler) AcceptNewMessage(message *tgbotapi.Message) error {
	id := message.Chat.ID

	activeCommand, errActive := h.usecasesData.Get(id, activeCommandKey)
	command, errCommand := extractCommand(message.Text)

	if errActive != nil {
		// random message
		if errCommand != nil {
			return errors.New("need to activate separate usecase, ex. /help")
		}

		// new command
		h.usecasesData.Set(id, activeCommandKey, command)
		return h.ExecuteUsecase(message, command)
	}

	// new command during existing command
	if errCommand == nil {
		err := h.usecasesData.Del(id)

		if err != nil {
			h.log.Warnf("failed to delete data with id %v", id)
		}

		h.usecasesData.Set(id, activeCommandKey, command)
		return h.ExecuteUsecase(message, command)
	}

	return h.ExecuteUsecase(message, activeCommand)
}

func (h *MessageHandler) ExecuteUsecase(message *tgbotapi.Message, command string) error {
	id := message.Chat.ID
	usecase, exists := h.registeredCallbacks[command]
	if !exists {
		h.log.Debugf("Failed to find callback for command %v", command)
		return errors.New("unrecognised command, also create separate usecase")
	}

	dataAccessor := temporal_storage.CreateTemporalStorageAccessor(h.usecasesData, id)
	mes, status := usecase.Handle(message, &dataAccessor)

	if status == usecases.Finished || status == usecases.Error {
		err := h.usecasesData.Del(id)
		if err != nil {
			h.log.Warnf("failed to delete data with id %v", id)
		}
	}

	_, err := h.bot.Send(mes)
	if err != nil {
		h.log.Error("failed to send message, chat id %v", id)
		return errors.New("failed to send message")
	}

	return nil
}
