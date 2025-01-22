package models

type Post struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	CategoryID int    `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsApproved bool   `json:"is_approved"`
	ViewCount  int    `json:"view_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
