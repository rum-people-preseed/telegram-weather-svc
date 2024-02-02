package controller

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/message_tools/message_reader"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
)

type MessageHandler struct {
	registeredFactories map[string]UsecaseFactory
	activeUsecases      map[int64]Usecase
	log                 models.Logger
	bot                 *models.Bot
}

func NewMessageHandler(bot *models.Bot, log models.Logger) *MessageHandler {
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
	message, chatID := update.Message, message_reader.GetChatId(update)
	command, err := message_reader.GetCommand(message)
	gotNewCommand := err == nil

	defer func() {
		err := h.EndCallback(update)
		if err != nil {
			h.log.Warnf("Error during closing callback for chat %v", chatID)
		}
	}()

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
	chatID := message_reader.GetChatId(update)
	activeUsecase, exists := h.activeUsecases[chatID]

	if !exists {
		h.log.Warnf("usecase does not exists %v", chatID)
		return errors.New("usecase does not exists")
	}

	msg, status := activeUsecase.Handle(update)

	if status == Finished || status == Error {
		delete(h.activeUsecases, chatID)
	}

	if msg != nil {
		return h.bot.SendMessage(msg)
	}

	return nil
}

func (h *MessageHandler) EndCallback(update *tgbotapi.Update) error {
	var err error
	if update.CallbackQuery != nil {
		_, err = h.bot.BotAPI.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
	}
	return err
}
