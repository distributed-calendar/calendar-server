package timepad

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
)

type EventProvider interface {
	AddEvent(ctx context.Context, event *domain.Event) error
}
