package telegram

import (
	"context"
	"log/slog"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/distributed-calendar/calendar-server/internal/service/telegram/provider"
)

type Service struct {
	userProvider  provider.UserProvider
	cacheProvider provider.CacheProvider
}

func NewService(
	userProvider provider.UserProvider,
	cacheProvider provider.CacheProvider,
) *Service {
	return &Service{
		userProvider:  userProvider,
		cacheProvider: cacheProvider,
	}
}

func (s *Service) CreateUserByTelegramID(ctx context.Context, telegramID int64, name, surname string) (*domain.User, error) {
	user, err := s.userProvider.CreateUser(ctx, telegramID, name, surname)
	if err != nil {
		return nil, err
	}

	err = s.cacheProvider.SetUser(ctx, telegramID, user)
	if err != nil {
		slog.Error("error setting user to cache", err)
	}

	return user, nil
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	user, err := s.cacheProvider.GetUser(ctx, telegramID)
	slog.Info("got user")
	if err == nil {
		return user, nil
	}

	slog.Error("error getting user from cache", err)

	return s.userProvider.GetUserByTelegramID(ctx, telegramID)
}
