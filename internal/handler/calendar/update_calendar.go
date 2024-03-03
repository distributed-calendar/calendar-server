package calendar

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/distributed-calendar/calendar-server/internal/handler/calendar/models"
)

func (h *Handler) updateCalendar(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateCalendarRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("cannot decode update calendar request", err)

		return
	}

	ctx := r.Context()

	calendar := &domain.Calendar{
		ID:      req.CalendarID,
		OwnerID: req.NewOwnerID,
	}

	err = h.calendarService.Update(ctx, calendar)
	if err != nil {
		slog.Error("cannot update calendar", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
