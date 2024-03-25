package event

import (
	"context"
	"time"

	"github.com/distributed-calendar/calendar-server/domain"
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

func (r *Repo) GetEvent(ctx context.Context, eventID int) (*domain.Event, error) {
	// TODO
	return nil, nil
}

func (r *Repo) AddEvent(ctx context.Context, event *domain.Event) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO event(name, start_time, end_time, description, creator_id) VALUES ($1, $2, $3, $4, $5)",
		event.Name,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.CreatorID,
	)

	return err
}

func (r *Repo) GetUserEvents(ctx context.Context, userID int, fromDt, toDt time.Time) ([]*domain.Event, error) {
	var events []*domain.Event

	rows, err := r.db.Query(
		ctx,
		"SELECT id, name, start_time, end_time, description FROM event WHERE creator_id = $1 AND $2 <= start_time AND start_time <= $3",
		userID,
		fromDt,
		toDt,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		event := &domain.Event{}
		if err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *Repo) UpdateEvent(ctx context.Context, event *domain.Event) error {
	// TODO
	return nil
}
