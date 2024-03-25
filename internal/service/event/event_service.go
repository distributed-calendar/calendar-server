package event

import (
	"context"
	"time"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/distributed-calendar/calendar-server/internal/service/event/provider"
)

type Service struct {
	eventProvider provider.EventProvider
}

func NewService(
	eventProvider provider.EventProvider,
) *Service {
	return &Service{
		eventProvider: eventProvider,
	}
}

func (s *Service) GetEvent(ctx context.Context, eventID int) (*domain.Event, error) {
	return s.eventProvider.GetEvent(ctx, eventID)
}

func (s *Service) AddEvent(ctx context.Context, event *domain.Event) error {
	return s.eventProvider.AddEvent(ctx, event)
}

func (s *Service) GetUserEvents(ctx context.Context, userID int, fromDt, toDt time.Time) ([]*domain.Event, error) {
	return s.eventProvider.GetUserEvents(ctx, userID, fromDt, toDt)
}

func (s *Service) UpdateEvent(ctx context.Context, event *domain.Event) error {
	return s.UpdateEvent(ctx, event)
}
