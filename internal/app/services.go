package app

import (
	eventrepo "github.com/distributed-calendar/calendar-server/internal/repo/event"
	"github.com/distributed-calendar/calendar-server/internal/repo/user"
	eventservice "github.com/distributed-calendar/calendar-server/internal/service/event"
	"github.com/distributed-calendar/calendar-server/internal/service/telegram"
	"github.com/distributed-calendar/calendar-server/internal/service/timepad"
)

func (a *App) initServices() {
	a.initEventService()
	a.initTelegramService()
	a.initTimepadService()
}

func (a *App) initEventService() {
	eventRepo := eventrepo.NewRepo(a.pgConnPool)
	a.eventService = eventservice.NewService(eventRepo)
}

func (a *App) initTelegramService() {
	userRepo := user.NewRepo(a.pgConnPool)
	a.telegramService = telegram.NewService(userRepo, a.cacheAdapter)
}

func (a *App) initTimepadService() {
	a.timepadService = timepad.NewService(
		a.timepadAdapter,
		a.telegramService,
		a.eventService,
	)
}
