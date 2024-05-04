package timepad

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
)

type TimepadProvider interface {
	GetEvent(ctx context.Context, eventID int) (*domain.Event, error)
}
