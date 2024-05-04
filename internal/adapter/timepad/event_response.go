package timepad

import "time"

type EventResponse struct {
	ID               int        `json:"id"`
	StartsAt         time.Time  `json:"starts_at"`
	EndsAt           *time.Time `json:"ends_at,omitempty"`
	Name             string     `json:"name"`
	DescriptionShort *string    `json:"description_short,omitempty"`
}
