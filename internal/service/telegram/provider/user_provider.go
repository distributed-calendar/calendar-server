package provider

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
)

type UserProvider interface {
	CreateUser(ctx context.Context, telegramID int64, name, surname string) (*domain.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
}
