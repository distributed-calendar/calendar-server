package telegram

import (
	"context"
	"log/slog"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/distributed-calendar/calendar-server/domain"
)

type createEventState struct {
	name        string
	startTime   time.Time
	endTime     time.Time
	description string
	notifyTime  time.Time
}

func (b *botAPI) handleCreateEvent(_ *echotron.Update) stateFn {
	res, err := b.SendMessage("Введите название события", b.chatID, nil)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleDefault
	}

	return b.handleCreateEventName
}

func (b *botAPI) handleCreateEventName(update *echotron.Update) stateFn {
	name := update.Message.Text

	state := &createEventState{
		name: name,
	}

	res, err := b.SendMessage("Введите время начала события", b.chatID, nil)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleDefault
	}

	return wrapStateFn(b.handleCreateEventStartTime, state)
}

func (b *botAPI) handleCreateEventStartTime(update *echotron.Update, state *createEventState) stateFn {
	startTime, err := time.Parse(time.DateTime, update.Message.Text)
	if err != nil {
		if res, e := b.SendMessage("Некорректная дата. Укажите дату в формате '2006-01-02 15:04:05'", b.chatID, nil); e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return wrapStateFn(b.handleCreateEventStartTime, state)
	}

	state.startTime = startTime

	res, err := b.SendMessage("Введите время окончания события", b.chatID, nil)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleDefault
	}

	return wrapStateFn(b.handleCreateEventEndTime, state)
}

func (b *botAPI) handleCreateEventEndTime(update *echotron.Update, state *createEventState) stateFn {
	endTime, err := time.Parse(time.DateTime, update.Message.Text)
	if err != nil {
		if res, e := b.SendMessage("Некорректная дата. Укажите дату в формате '2006-01-02 15:04:05'", b.chatID, nil); e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return wrapStateFn(b.handleCreateEventEndTime, state)
	}

	if endTime.Before(state.startTime) {
		if res, e := b.SendMessage("Дата окончания раньше, чем дата начала. Укажите корректную дату", b.chatID, nil); e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return wrapStateFn(b.handleCreateEventEndTime, state)
	}

	state.endTime = endTime

	res, err := b.SendMessage("Введите описание события", b.chatID, nil)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleDefault
	}

	return wrapStateFn(b.handleCreateEventDescription, state)
}

func (b *botAPI) handleCreateEventDescription(update *echotron.Update, state *createEventState) stateFn {
	desc := update.Message.Text
	state.description = desc

	user, err := b.services.telegramService.GetUserByTelegramID(context.Background(), b.chatID)
	if err != nil {
		slog.Error("error getting user by telegram id", err)

		if res, e := b.SendMessage("Внутренняя ошибка. Попробуйте позже", b.chatID, nil); e != nil {
			logSendEchotronError(res, e)
		}

		return b.handleDefault
	}

	event := &domain.Event{
		Name:        state.name,
		StartTime:   state.startTime,
		EndTime:     state.endTime,
		Description: state.description,
		NotifyTime:  state.notifyTime,
		CreatorID:   user.ID,
	}
	err = b.services.eventService.AddEvent(context.Background(), event)
	if err != nil {
		slog.Error("error adding event", err)

		if res, e := b.SendMessage("Внутренняя ошибка. Попробуйте позже", b.chatID, nil); e != nil {
			logSendEchotronError(res, e)
		}

		return b.handleDefault
	}

	return b.handleDefault
}
