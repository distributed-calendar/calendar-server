package telegram

import (
	"strings"

	"github.com/NicoNex/echotron/v3"
)

func (b *botAPI) handleSetupYandexCalendar(_ *echotron.Update) stateFn {
	_, _ = b.SendMessage("Перейдите по ссылке https://oauth.yandex.ru/authorize?response_type=token&client_id=5d42422560224e1ea7adedaaa5052aa1", b.chatID, nil)

	return b.handeSetupYandexCalendarToken
}

func (b *botAPI) handeSetupYandexCalendarToken(update *echotron.Update) stateFn {
	token := strings.Split(update.Message.Text, " ")[1]

	_, _ = b.SendMessage(token, b.chatID, nil)

	return b.handleDefault
}
