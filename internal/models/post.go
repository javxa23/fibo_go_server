package models

type Post struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	AuthorID int     `json:"author_id"`
	Rating   float64 `json:"rating"`
}
