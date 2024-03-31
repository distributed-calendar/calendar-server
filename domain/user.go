package domain

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	TelegramID int64  `json:"telegram_id"`
}
