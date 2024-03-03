package models

type UpdateCalendarRequest struct {
	CalendarID string `json:"calendarId"`
	NewOwnerID string `json:"newOwnerId"`
}
