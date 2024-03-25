package telegram

import "github.com/NicoNex/echotron/v3"

func (b *botAPI) handleDefault(update *echotron.Update) stateFn {
	switch update.Message.Text {
	case commandStart:
		return b.handleStart(update)
	case commandCreateEvent:
		return b.handleCreateEvent(update)
	case commandGetEvents:
		return b.handleGetEvents(update)
	default:
		res, err := b.SendMessage("Неизвестная команда", b.chatID, nil)
		if err != nil {
			logSendEchotronError(res, err)
		}
	}

	return b.handleDefault
}
