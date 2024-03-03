package calendar

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) getCalendar(w http.ResponseWriter, r *http.Request) {
	calendarID := chi.URLParam(r, "calendarId")
	if calendarID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	calendar, err := h.calendarService.Calendar(ctx, calendarID)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		slog.Error("cannot get calendar", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	resp := &CalendarResponse{
		CalendarID:      calendarID,
		CalendarOwnerID: calendar.OwnerID,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Error("cannot encode calendar response", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
