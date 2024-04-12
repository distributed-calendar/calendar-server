package app

import "github.com/distributed-calendar/calendar-server/internal/port/telegram"

func (a *App) initPorts() {
	a.initTelegramBot()
}

func (a *App) initTelegramBot() {
	bot, err := telegram.NewBot(a.cfg.Telegram.BotToken, a.cfg.Telegram.WebhookURL, a.telegramService, a.eventService, a.httpServer)
	if err != nil {
		panic(err)
	}

	a.addOnRun(bot.StartBot)
}
