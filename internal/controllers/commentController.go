package controllers

import (
	"database/sql"
	"net/http"

	"fibo_go_server/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment models.Comment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO comments (post_id, user_id, parent_comment_id, content, created_at) 
		          VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
		err := db.QueryRow(query, comment.PostID, comment.UserID, comment.ParentCommentID, comment.Content).Scan(&comment.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, comment)
	}
}
