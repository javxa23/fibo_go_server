package models

type Salary struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Amount    int    `json:"amount"`
	Month     string `json:"month"`
	Year      string `json:"year"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
