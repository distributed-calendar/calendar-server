package calendar

import (
	"net/http"

	"github.com/distributed-calendar/calendar-server/internal/service/calendar"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	calendarService *calendar.Service
}

func NewHandler(
	calendarService *calendar.Service,
) http.Handler {
	mux := chi.NewMux()

	handler := &Handler{
		calendarService: calendarService,
	}

	mux.Get("/calendar/{calendarId}", handler.getCalendar)
	mux.Put("/calendar/{ownerId}", handler.createCalendar)
	mux.Post("/calendar", handler.updateCalendar)

	return mux
}
