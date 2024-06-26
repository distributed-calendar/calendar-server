package telegram

import (
	"context"
	"log/slog"
	"time"

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

	cacheCtx, cfn := context.WithTimeout(ctx, 2*time.Second)
	defer cfn()

	err = s.cacheProvider.SetUser(cacheCtx, telegramID, user)
	if err != nil {
		slog.Error("error setting user to cache", err)
	}

	return user, nil
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	cacheCtx, cfn := context.WithTimeout(ctx, 2*time.Second)
	defer cfn()

	user, err := s.cacheProvider.GetUser(cacheCtx, telegramID)
	slog.Info("got user")
	if err == nil {
		return user, nil
	}

	slog.Error("error getting user from cache", err)

	return s.userProvider.GetUserByTelegramID(ctx, telegramID)
}
