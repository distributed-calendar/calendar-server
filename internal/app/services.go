package app

import (
	calendarrepo "github.com/distributed-calendar/calendar-server/internal/repo/calendar"
	calendarservice "github.com/distributed-calendar/calendar-server/internal/service/calendar"
)

func (a *App) initCalendarService() {
	calendarRepo := calendarrepo.NewRepo(a.pgConnPool)
	a.calendarService = calendarservice.NewService(calendarRepo)
}
