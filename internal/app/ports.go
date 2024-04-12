package app

import "github.com/distributed-calendar/calendar-server/internal/port/telegram"

func (a *App) initPorts() {
	a.initTelegramBot()
}

func (a *App) initTelegramBot() {
	_, fn, err := telegram.NewBot(a.cfg.Telegram.BotToken, a.telegramService, a.eventService)
	if err != nil {
		panic(err)
	}

	a.mux.Post("/", fn)
}
