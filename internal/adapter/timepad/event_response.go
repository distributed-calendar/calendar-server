package timepad

type EventResponse struct {
	ID               int     `json:"id"`
	StartsAt         string  `json:"starts_at"`
	EndsAt           *string `json:"ends_at,omitempty"`
	Name             string  `json:"name"`
	DescriptionShort *string `json:"description_short,omitempty"`
}
