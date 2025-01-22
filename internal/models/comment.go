package models

type Comment struct {
	ID              int       `json:"id"`
	PostID          int       `json:"post_id"`
	UserID          *int      `json:"user_id"`
	ParentCommentID *int      `json:"parent_comment_id"`
	Content         string    `json:"content"`
	CreatedAt       string    `json:"created_at"`
	Children        []Comment `json:"children"`
}
