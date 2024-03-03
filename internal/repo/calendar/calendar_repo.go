package calendar

import (
	"context"
	"errors"
	"log/slog"

	"github.com/distributed-calendar/calendar-server/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	createCalendarSQL = `INSERT INTO calendar(owner_id) VALUES ($1) RETURNING calendar_id`

	getCalendarSQL = `SELECT owner_id FROM calendar WHERE calendar_id = $1`

	updateCalendarSQL = `UPDATE calendar SET owner_id = $2 WHERE calendar_id = $1`
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
	calendar := &domain.Calendar{
		ID: id,
	}

	row := repo.db.QueryRow(ctx, getCalendarSQL, id)

	err := row.Scan(
		&calendar.OwnerID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrDataNotFound
	}

	return calendar, err
}

func (repo *Repo) Create(ctx context.Context, calendar *domain.Calendar) (string, error) {
	var calendarID string

	row := repo.db.QueryRow(ctx, createCalendarSQL, calendar.OwnerID)
	err := row.Scan(&calendarID)

	return calendarID, err
}

func (repo *Repo) Update(ctx context.Context, calendar *domain.Calendar) error {
	slog.Debug("", "id", calendar.ID, "ownerid", calendar.OwnerID)
	_, err := repo.db.Exec(ctx, updateCalendarSQL, calendar.ID, calendar.OwnerID)

	return err
}
