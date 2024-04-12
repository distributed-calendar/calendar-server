package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NicoNex/echotron/v3"
	"github.com/distributed-calendar/calendar-server/internal/service/event"
	"github.com/distributed-calendar/calendar-server/internal/service/telegram"
)

type stateFn func(*echotron.Update) stateFn

const (
	commandStart       = "/start"
	commandCreateEvent = "/create_event"
	commandGetEvents   = "/get_events"
	commandCancel      = "/cancel"
)

type Bot struct {
	dispatcher *echotron.Dispatcher

	token      string
	webhookURL string
}

func (b *Bot) StartBot() error {
	err := b.dispatcher.ListenWebhook(fmt.Sprintf("%s/%s", b.webhookURL, b.token))
	if err != nil {
		slog.Error("error starting bot", err)
	}

	return err
}

type botAPI struct {
	echotron.API

	chatID int64
	state  stateFn

	services *services
}

type services struct {
	telegramService *telegram.Service
	eventService    *event.Service
}

func NewBot(
	token string,
	webhookURL string,
	telegramService *telegram.Service,
	eventService *event.Service,
) (*Bot, error) {
	api, err := newAPI(token)
	if err != nil {
		return nil, fmt.Errorf("cannot create echotron API: %w", err)
	}

	services := &services{
		telegramService: telegramService,
		eventService:    eventService,
	}

	dispatcher := echotron.NewDispatcher(
		token,
		newBotAPICreator(
			api,
			services,
		),
	)

	bot := &Bot{
		dispatcher: dispatcher,
		token:      token,
		webhookURL: webhookURL,
	}

	return bot, nil
}

func newAPI(token string) (echotron.API, error) {
	api := echotron.NewAPI(token)

	return api, nil
}

func newBotAPICreator(api echotron.API, services *services) echotron.NewBotFn {
	return func(chatID int64) echotron.Bot {
		return &botAPI{
			api,
			chatID,
			nil,
			services,
		}
	}
}

func (b *botAPI) Update(update *echotron.Update) {
	slog.Info("got new message")

	if b.state == nil {
		if _, err := b.services.telegramService.GetUserByTelegramID(context.Background(), b.chatID); err != nil {
			b.state = b.handleStart
		} else {
			b.state = b.handleDefault
		}
	}

	b.state = b.state(update)
}

func logSendEchotronError(res echotron.APIResponseMessage, err error) {
	slog.Error("cannot send echotron message", "error", err, "error_code", res.ErrorCode, "description", res.Description)
}
