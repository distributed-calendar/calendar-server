package provider

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
)

type CalendarProvider interface {
	Calendar(ctx context.Context, id string) (*domain.Calendar, error)
	Create(ctx context.Context, calendar *domain.Calendar) error
}
