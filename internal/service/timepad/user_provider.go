package timepad

import (
	"context"

	"github.com/distributed-calendar/calendar-server/domain"
)

type TelegramProvider interface {
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
}
