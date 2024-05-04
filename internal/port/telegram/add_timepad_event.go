package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/NicoNex/echotron/v3"
)

func (b *botAPI) handleAddTimepadEvent(_ *echotron.Update) stateFn {
	res, err := b.SendMessage("Введите ID события из Timepad", b.chatID, nil)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleDefault
	}

	return b.handleAddTimepadEvent_GetID
}

func (b *botAPI) handleAddTimepadEvent_GetID(update *echotron.Update) stateFn {
	idString := update.Message.Text

	eventID, err := strconv.Atoi(idString)
	if err != nil {
		res, e := b.SendMessage("Вы ввели не число. Введите еще раз", b.chatID, nil)
		if e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return b.handleAddTimepadEvent_GetID
	}

	event, err := b.services.timepadService.GetEvent(context.Background(), b.chatID, eventID)
	if err != nil {
		slog.Error("cannot add timepad event from bot", err)

		res, e := b.SendMessage("Что-то пошло не так. Возможно, такого события не существует. Попробуйте еще раз позже", b.chatID, nil)
		if e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return b.handleDefault
	}

	eventMsg := fmt.Sprintf(
		"Событие успешно добавлено\nНазвание: %s\nДата начала: %s\nДата окончания: %s\nОписание: %s",
		event.Name,
		event.StartTime.String(),
		event.EndTime.String(),
		event.Description,
	)

	res, e := b.SendMessage(eventMsg, b.chatID, nil)
	if e != nil {
		logSendEchotronError(res, e)
	}

	return b.handleDefault
}
