package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`    // Exclude from JSON responses
	Role     string `json:"role"` // Admin || User
}
