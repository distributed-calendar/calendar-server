package telegram

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/distributed-calendar/calendar-server/internal/service/telegram/provider"
)

type Service struct {
	userProvider provider.UserProvider
}

func NewService(
	userProvider provider.UserProvider,
) *Service {
	return &Service{
		userProvider: userProvider,
	}
}

func (s *Service) CreateUser(ctx context.Context, telegramID int64, name, surname string) (*domain.User, error) {
	return s.userProvider.CreateUser(ctx, telegramID, name, surname)
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	return s.userProvider.GetUserByTelegramID(ctx, telegramID)
}
