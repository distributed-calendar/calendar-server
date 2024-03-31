package app

import (
	"github.com/distributed-calendar/calendar-server/internal/repo/event"
	"github.com/distributed-calendar/calendar-server/internal/repo/user"
)

func (a *App) initRepos() {
	a.initEventRepo()
	a.initUserRepo()
}

func (a *App) initEventRepo() {
	a.eventRepo = event.NewRepo(a.pgConnPool)
}

func (a *App) initUserRepo() {
	a.userRepo = user.NewRepo(a.pgConnPool)
}
