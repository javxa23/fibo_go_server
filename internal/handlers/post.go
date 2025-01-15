package handlers

import (
	"database/sql"
	"fibo_go_server/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var post models.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec("INSERT INTO posts (title, content, author_id, rating) VALUES ($1, $2, $3, $4)",
			post.Title, post.Content, post.AuthorID, 0.0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
	}
}
