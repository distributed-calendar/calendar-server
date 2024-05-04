package timepad

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
)

type Service struct {
	timepadProvider  TimepadProvider
	telegramProvider TelegramProvider
	eventProvider    EventProvider
}

func NewService(
	timepadProvider TimepadProvider,
	telegramProvider TelegramProvider,
	eventProvider EventProvider,
) *Service {
	return &Service{
		timepadProvider:  timepadProvider,
		telegramProvider: telegramProvider,
		eventProvider:    eventProvider,
	}
}

func (s *Service) GetEvent(ctx context.Context, telegramUserID int64, eventID int) (*domain.Event, error) {
	event, err := s.timepadProvider.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	user, err := s.telegramProvider.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, err
	}

	event.CreatorID = user.ID

	return event, s.eventProvider.AddEvent(ctx, event)
}
