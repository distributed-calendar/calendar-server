package calendar

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/distributed-calendar/calendar-server/internal/handler/calendar/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) createCalendar(w http.ResponseWriter, r *http.Request) {
	ownerID := chi.URLParam(r, "ownerId")
	if ownerID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	calendar := &domain.Calendar{
		OwnerID: ownerID,
	}

	calendarID, err := h.calendarService.Create(ctx, calendar)
	if err != nil {
		slog.Error("cannot create calendar", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	resp := &models.CalendarResponse{
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
