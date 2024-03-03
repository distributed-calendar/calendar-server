package calendar

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/distributed-calendar/calendar-server/internal/service/calendar/provider"
)

type Service struct {
	calendarProvider provider.CalendarProvider
}

func NewService(
	calendarProvider provider.CalendarProvider,
) *Service {
	return &Service{
		calendarProvider: calendarProvider,
	}
}

func (service *Service) Calendar(ctx context.Context, id string) (*domain.Calendar, error) {
	return service.calendarProvider.Calendar(ctx, id)
}

func (service *Service) Create(ctx context.Context, calendar *domain.Calendar) error {
	return service.calendarProvider.Create(ctx, calendar)
}
