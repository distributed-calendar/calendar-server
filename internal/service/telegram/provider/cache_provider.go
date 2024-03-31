package provider

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
)

type CacheProvider interface {
	SetUser(ctx context.Context, telegramID int64, user *domain.User) error
	GetUser(ctx context.Context, telegramID int64) (*domain.User, error)
}
