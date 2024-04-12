package cache

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/redis/go-redis/v9"
)

type Adapter struct {
	client redis.UniversalClient
}

func NewAdapter(
	addrs string,
	password string,
	certPath string,
) (*Adapter, error) {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, fmt.Errorf("failed to append PEM")
	}

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    strings.Split(addrs, ","),
		Password: password,
		TLSConfig: &tls.Config{
			RootCAs:            rootCertPool,
		},
	})

	return &Adapter{
		client: client,
	}, nil
}

func (a *Adapter) Set(ctx context.Context, key string, value any) error {
	return a.client.Set(ctx, key, value, 0).Err()
}

func (a *Adapter) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key).Result()
}

func (a *Adapter) SetUser(ctx context.Context, telegramID int64, user *domain.User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return a.Set(ctx, userCacheKeyByTelegramID(telegramID), b)
}

func (a *Adapter) GetUser(ctx context.Context, telegramID int64) (*domain.User, error) {
	data, err := a.Get(ctx, userCacheKeyByTelegramID(telegramID))
	if err != nil {
		return nil, err
	}

	var user domain.User

	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func userCacheKeyByTelegramID(telegramID int64) string {
	return fmt.Sprintf("user.telegram.%d", telegramID)
}
