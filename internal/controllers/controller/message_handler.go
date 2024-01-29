package controller

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/repositories/temporal_storage"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/utils"
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
	bot                 *models.Bot
}

func NewMessageHandler(bot *models.Bot, log Logger, usecasesData temporal_storage.TemporalStorage) *MessageHandler {
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

func (h *MessageHandler) AcceptNewUpdate(update *tgbotapi.Update) error {
	message, chatID := update.Message, utils.GetChatId(update)

	activeCommand, errActive := h.usecasesData.Get(chatID, activeCommandKey)
	command, errCommand := utils.ExtractCommand(message)

	println("Active " + activeCommand)
	println("Command " + command)

	if errActive != nil {
		if errCommand == nil {
			//textError := "Unrecognized command! Please, see /help"
			//msgCfg := utils.GetSimpleMessage(chatID, textError)
			//_ = h.bot.SendMessage(&msgCfg)
			//return errors.New(textError)
			h.usecasesData.Set(chatID, activeCommandKey, command)
		}

		return h.ExecuteUsecase(update, command)
	}

	if errCommand == nil {
		// start new command
		err := h.usecasesData.Del(chatID)
		if err != nil {
			h.log.Warnf("failed to delete data with chatID %v", chatID)
		}

		h.usecasesData.Set(chatID, activeCommandKey, command)
		return h.ExecuteUsecase(update, command)
	}

	return h.ExecuteUsecase(update, activeCommand)
}

func (h *MessageHandler) ExecuteUsecase(update *tgbotapi.Update, command string) error {
	chatID := utils.GetChatId(update)

	usecase, exists := h.registeredCallbacks[command]
	if !exists {
		// todo: need to be simplified
		textError := "Unrecognized command! Please, see /help"
		msgCfg := utils.GetSimpleMessage(chatID, textError)
		_ = h.bot.SendMessage(&msgCfg)
		h.EndCallback(update)
		return errors.New(textError)
	}

	dataAccessor := temporal_storage.CreateTemporalStorageAccessor(h.usecasesData, chatID)
	msg, status := usecase.Handle(update, &dataAccessor)

	if status == usecases.Finished || status == usecases.Error {
		err := h.usecasesData.Del(chatID)
		if err != nil {
			h.log.Warnf("Failed to delete data with chatID %v", chatID)
		}
	}

	h.EndCallback(update)
	return h.bot.SendMessage(msg)
}

func (h *MessageHandler) EndCallback(update *tgbotapi.Update) {
	// todo: to think how to handle it in correct way
	if update.CallbackQuery != nil {
		_, _ = h.bot.BotAPI.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
	}
}
