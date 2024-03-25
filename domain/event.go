package domain

import "time"

type Event struct {
	ID          int
	Name        string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	NotifyTime  time.Time
	CreatorID   int
}
