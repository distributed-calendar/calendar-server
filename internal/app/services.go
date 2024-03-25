package app

import (
	telegrambot "github.com/distributed-calendar/calendar-server/internal/port/telegram"
	eventrepo "github.com/distributed-calendar/calendar-server/internal/repo/event"
	"github.com/distributed-calendar/calendar-server/internal/repo/user"
	eventservice "github.com/distributed-calendar/calendar-server/internal/service/event"
	"github.com/distributed-calendar/calendar-server/internal/service/telegram"
)

func (a *App) initEventService() {
	eventRepo := eventrepo.NewRepo(a.pgConnPool)
	a.eventService = eventservice.NewService(eventRepo)
}

func (a *App) initTelegramService() {
	userRepo := user.NewRepo(a.pgConnPool)
	a.telegramService = telegram.NewService(userRepo)
}

func (a *App) initTelegramBot() {
	bot, err := telegrambot.NewBot(a.cfg.Telegram.BotToken, a.telegramService, a.eventService)
	if err != nil {
		panic(err)
	}

	a.addOnRun(bot.StartBot)
}
