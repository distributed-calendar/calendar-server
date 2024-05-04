package telegram

import (
	"context"
	"errors"
	"log/slog"

	"github.com/NicoNex/echotron/v3"
	"github.com/distributed-calendar/calendar-server/domain"
)

func (b *botAPI) handleStart(update *echotron.Update) stateFn {
	if update.Message.Text != commandStart {
		res, err := b.SendMessage("Для начала введите /start", b.chatID, nil)
		if err != nil {
			logSendEchotronError(res, err)
		}

		return b.handleStart
	}

	_, err := b.services.telegramService.CreateUserByTelegramID(
		context.Background(),
		b.chatID,
		update.Message.From.FirstName,
		update.Message.From.LastName,
	)
	if err != nil && !errors.Is(err, domain.ErrAlreadyExist) {
		slog.Error("error while creating user", err)

		res, e := b.SendMessage("Ошибка на сторона сервера, попробуйте позже.", b.chatID, nil)
		if e != nil {
			logSendEchotronError(res, err)
		}

		return b.handleStart
	}

	res, err := b.SendMessage("Привет!", b.chatID, nil)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleStart
	}

	_, err = b.SetMyCommands(
		nil,
		echotron.BotCommand{
			Command:     commandCreateEvent,
			Description: "Создать событие",
		},
		echotron.BotCommand{
			Command:     commandGetEvents,
			Description: "Получить все события",
		},
		echotron.BotCommand{
			Command:     commandAddTimepadEvent,
			Description: "Добавить событие из Timepad",
		},
	)
	if err != nil {
		slog.Error("cannot set bot commands", err)
	}

	return b.handleDefault
}
