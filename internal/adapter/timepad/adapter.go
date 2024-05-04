package timepad

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/go-resty/resty/v2"
)

type Adapter struct {
	client *resty.Client
}

func NewAdapter(url, token string) *Adapter {
	client := resty.New().SetBaseURL(url).SetAuthToken(token)

	return &Adapter{
		client: client,
	}
}

func (a *Adapter) GetEvent(ctx context.Context, eventID int) (*domain.Event, error) {
	var resp EventResponse

	_, err := a.client.R().
		SetContext(ctx).
		SetResult(&resp).
		Get(fmt.Sprintf("/events/%d", eventID))
	if err != nil {
		slog.Error("cannot get timepad event", err)

		return nil, domain.ErrExternalUnavailable
	}

	return mapTimepadEventToDomain(&resp), nil
}

func mapTimepadEventToDomain(e *EventResponse) *domain.Event {
	domainEvent := &domain.Event{
		Name:      e.Name,
		StartTime: e.StartsAt,
	}

	if e.EndsAt != nil {
		domainEvent.EndTime = *e.EndsAt
	}

	if e.DescriptionShort != nil {
		domainEvent.Description = *e.DescriptionShort
	}

	return domainEvent
}
