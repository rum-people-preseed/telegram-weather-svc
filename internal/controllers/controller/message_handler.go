package controller

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
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

type MessageHandler struct {
	registeredFactories map[string]UsecaseFactory
	activeUsecases      map[int64]Usecase
	log                 Logger
	bot                 *models.Bot
}

func NewMessageHandler(bot *models.Bot, log Logger) *MessageHandler {
	return &MessageHandler{
		bot:                 bot,
		log:                 log,
		activeUsecases:      make(map[int64]Usecase),
		registeredFactories: make(map[string]UsecaseFactory),
	}
}

func (h *MessageHandler) RegisterUsecaseFactory(usecaseFactory UsecaseFactory) error {
	command := usecaseFactory.Command()
	_, ok := h.registeredFactories[command]
	if ok {
		return errors.New(fmt.Sprintf("UsecaseFactory for command %v is already registered", command))
	}

	h.registeredFactories[command] = usecaseFactory
	return nil
}

func (h *MessageHandler) ActivateUsecase(chatID int64, command string) error {
	factory, exists := h.registeredFactories[command]
	if !exists {
		h.log.Warnf("factory for command %v does not exists", command)
		return errors.New("factory does not exists")
	}

	h.activeUsecases[chatID] = factory.Create()
	return nil
}

func (h *MessageHandler) AcceptNewUpdate(update *tgbotapi.Update) error {
	message, chatID := update.Message, utils.GetChatId(update)
	command, err := utils.ExtractCommand(message)
	gotNewCommand := err == nil

	if gotNewCommand {
		err = h.ActivateUsecase(chatID, command)
		if err != nil {
			h.log.Warnf("failed to activate usecase %v", err)
			return errors.New("failed to activate usecase")
		}
	}

	return h.ExecuteUsecase(update)
}

func (h *MessageHandler) ExecuteUsecase(update *tgbotapi.Update) error {
	chatID := utils.GetChatId(update)
	activeUsecase, exists := h.activeUsecases[chatID]

	if !exists {
		h.log.Warnf("usecase does not exists %v", chatID)
		return nil
	}

	msg, status := activeUsecase.Handle(update)

	if status == Finished || status == Error {
		delete(h.activeUsecases, chatID)
	}

	//what is that?
	//h.EndCallback(update)
	if msg != nil {
		return h.bot.SendMessage(msg)
	}

	return nil
}

func (h *MessageHandler) EndCallback(update *tgbotapi.Update) {
	// todo: to think how to handle it in correct way
	if update.CallbackQuery != nil {
		_, _ = h.bot.BotAPI.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
	}
}
