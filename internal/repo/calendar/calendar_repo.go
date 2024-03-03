package calendar

import (
	"context"
	"errors"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	createCalendarSQL = `INSERT INTO calendar(owner_id) VALUES ($1)`

	getCalendarSQL = `SELECT calendar_id FROM calendar WHERE calendar_id = $1`
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo *Repo) Calendar(ctx context.Context, id string) (*domain.Calendar, error) {
	calendar := &domain.Calendar{}

	row := repo.db.QueryRow(ctx, getCalendarSQL, id)

	err := row.Scan(
		&calendar.ID,
		&calendar.OwnerID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrDataNotFound
	}

	return calendar, err
}

func (repo *Repo) Create(ctx context.Context, calendar *domain.Calendar) error {
	_, err := repo.db.Exec(ctx, createCalendarSQL, calendar.OwnerID)

	return err
}
