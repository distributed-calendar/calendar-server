package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateUser(ctx context.Context, telegramID int64, name, surname string) (*domain.User, error) {
	row := r.db.QueryRow(
		ctx,
		`INSERT INTO calendar_user (name, surname, telegram_id) VALUES ($1, $2, $3) RETURNING id`,
		name, surname, telegramID,
	)

	var id int
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAlreadyExist
		}

		return nil, fmt.Errorf("cannot create user: %w", err)
	}

	user := &domain.User{
		ID:         id,
		Name:       name,
		Surname:    surname,
		TelegramID: telegramID,
	}

	return user, nil
}

func (r *Repo) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	row := r.db.QueryRow(
		ctx,
		"SELECT id, name, surname FROM calendar_user WHERE telegram_id = $1",
		telegramID,
	)

	user := &domain.User{
		TelegramID: telegramID,
	}
	if err := row.Scan(&user.ID, &user.Name, &user.Surname); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}

		return nil, fmt.Errorf("cannot get user by telegram id: %w", err)
	}

	return user, nil
}
