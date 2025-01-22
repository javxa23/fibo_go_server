package controllers

import (
	"database/sql"
	"net/http"

	"fibo_go_server/internal/models"

	"github.com/gin-gonic/gin"
)

// AddLike handles the process of adding a like to a post
func AddLike(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var like models.Like
		if err := c.ShouldBindJSON(&like); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO likes (post_id, user_id) 
		          VALUES ($1, $2) RETURNING id`
		err := db.QueryRow(query, like.PostID, like.UserID).Scan(&like.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, like)
	}
}
