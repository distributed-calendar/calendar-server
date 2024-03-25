package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/NicoNex/echotron/v3"
)

const (
	yearDuration = 365 * 24 * time.Hour
)

func (b *botAPI) handleGetEvents(_ *echotron.Update) stateFn {
	user, err := b.services.telegramService.GetUserByTelegramID(context.Background(), b.chatID)
	if err != nil {
		slog.Error("error getting user by telegram id", err)

		if res, e := b.SendMessage("Внутренняя ошибка. Попробуйте позже", b.chatID, nil); e != nil {
			logSendEchotronError(res, e)
		}

		return b.handleDefault
	}

	slog.Info("user", slog.Int("id", user.ID))

	timeNow := time.Now()
	fromDt := timeNow.Add(-yearDuration)
	toDt := timeNow.Add(yearDuration)

	events, err := b.services.eventService.GetUserEvents(context.Background(), user.ID, fromDt, toDt)
	if err != nil {
		slog.Error("error getting user events", err)

		if res, e := b.SendMessage("Внутренняя ошибка. Попробуйте позже", b.chatID, nil); e != nil {
			logSendEchotronError(res, e)
		}

		return b.handleDefault
	}

	if len(events) == 0 {
		res, err := b.SendMessage("Нету событий", b.chatID, nil)
		if err != nil {
			logSendEchotronError(res, err)
		}
	}

	for _, event := range events {
		eventMsg := fmt.Sprintf(
			"Название: %s\nДата начала: %s\nДата окончания: %s\nОписание: %s\nСоздатель: %s %s",
			event.Name,
			event.StartTime.String(),
			event.EndTime.String(),
			event.Description,
			user.Name, user.Surname,
		)
		res, err := b.SendMessage(eventMsg, b.chatID, nil)
		if err != nil {
			logSendEchotronError(res, err)
		}
	}

	return b.handleDefault
}
