package telegram

import (
	"context"
	"log/slog"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/distributed-calendar/calendar-server/domain"
)

type createEventState struct {
	messageID   int
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

	state := &createEventState{
		messageID: res.Result.ID,
	}

	return withOpts(withState(b.handleCreateEventName, state), b.withCancel)
}

func (b *botAPI) handleCreateEventName(update *echotron.Update, state *createEventState) stateFn {
	name := update.Message.Text

	state.name = name

	// TODO
	_, err := b.DeleteMessage(b.chatID, update.Message.ID)
	if err != nil {
		slog.Error("cannot delete message", err)

		return b.handleDefault
	}

	editRes, err := b.EditMessageText(
		"Введите время начала события",
		echotron.NewMessageID(b.chatID, state.messageID),
		nil,
	)
	if err != nil {
		logSendEchotronError(editRes, err)

		return b.handleDefault
	}

	return withOpts(withState(b.handleCreateEventStartTime, state), b.withCancel)
}

func (b *botAPI) handleCreateEventStartTime(update *echotron.Update, state *createEventState) stateFn {
	// TODO
	_, err := b.DeleteMessage(b.chatID, update.Message.ID)
	if err != nil {
		slog.Error("cannot delete message", err)

		return b.handleDefault
	}

	startTime, err := time.Parse(time.DateTime, update.Message.Text)
	if err != nil {
		if res, e := b.EditMessageText(
			"Некорректная дата. Укажите дату в формате '2006-01-02 15:04:05'",
			echotron.NewMessageID(b.chatID, state.messageID),
			nil,
		); e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return withOpts(withState(b.handleCreateEventStartTime, state), b.withCancel)
	}

	state.startTime = startTime

	res, err := b.EditMessageText(
		"Введите время окончания события",
		echotron.NewMessageID(b.chatID, state.messageID),
		nil,
	)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleDefault
	}

	return withOpts(withState(b.handleCreateEventEndTime, state), b.withCancel)
}

func (b *botAPI) handleCreateEventEndTime(update *echotron.Update, state *createEventState) stateFn {
	// TODO
	_, err := b.DeleteMessage(b.chatID, update.Message.ID)
	if err != nil {
		slog.Error("cannot delete message", err)

		return b.handleDefault
	}

	endTime, err := time.Parse(time.DateTime, update.Message.Text)
	if err != nil {
		if res, e := b.EditMessageText(
			"Некорректная дата. Укажите дату в формате '2006-01-02 15:04:05'",
			echotron.NewMessageID(b.chatID, state.messageID),
			nil,
		); e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return withOpts(withState(b.handleCreateEventEndTime, state), b.withCancel)
	}

	if endTime.Before(state.startTime) {
		if res, e := b.EditMessageText(
			"Дата окончания раньше, чем дата начала. Укажите корректную дату",
			echotron.NewMessageID(b.chatID, state.messageID),
			nil,
		); e != nil {
			logSendEchotronError(res, e)

			return b.handleDefault
		}

		return withOpts(withState(b.handleCreateEventEndTime, state), b.withCancel)
	}

	state.endTime = endTime

	res, err := b.EditMessageText(
		"Введите описание события",
		echotron.NewMessageID(b.chatID, state.messageID),
		nil,
	)
	if err != nil {
		logSendEchotronError(res, err)

		return b.handleDefault
	}

	return withOpts(withState(b.handleCreateEventDescription, state), b.withCancel)
}

func (b *botAPI) handleCreateEventDescription(update *echotron.Update, state *createEventState) stateFn {
	// TODO
	_, err := b.DeleteMessage(b.chatID, update.Message.ID)
	if err != nil {
		slog.Error("cannot delete message", err)

		return b.handleDefault
	}

	desc := update.Message.Text
	state.description = desc

	user, err := b.services.telegramService.GetUserByTelegramID(context.Background(), b.chatID)
	if err != nil {
		slog.Error("error getting user by telegram id", err)

		if res, e := b.EditMessageText(
			"Внутренняя ошибка. Попробуйте позже",
			echotron.NewMessageID(b.chatID, state.messageID),
			nil,
		); e != nil {
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

		if res, e := b.EditMessageText(
			"Внутренняя ошибка. Попробуйте позже",
			echotron.NewMessageID(b.chatID, state.messageID),
			nil,
		); e != nil {
			logSendEchotronError(res, e)
		}

		return b.handleDefault
	}

	editRes, err := b.EditMessageText(
		"Событие успешно добавлено",
		echotron.NewMessageID(b.chatID, state.messageID),
		nil,
	)
	if err != nil {
		logSendEchotronError(editRes, err)

		return b.handleDefault
	}

	return b.handleDefault
}
