package provider

import (
	"context"
	"time"

	"github.com/distributed-calendar/calendar-server/domain"
)

type EventProvider interface {
	GetEvent(ctx context.Context, eventID int) (*domain.Event, error)
	AddEvent(ctx context.Context, event *domain.Event) error
	GetUserEvents(ctx context.Context, userID int, fromDt, toDt time.Time) ([]*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) error
}
