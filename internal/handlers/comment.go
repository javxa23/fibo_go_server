package handlers

import (
	"database/sql"
	"fibo_go_server/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddCommentHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment models.Comment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Insert the comment into the database
		_, err := db.Exec(
			"INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)",
			comment.PostID, comment.UserID, comment.Content,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully"})
	}
}

func GetCommentsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("postID")

		rows, err := db.Query("SELECT id, post_id, user_id, content FROM comments WHERE post_id = $1", postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
			return
		}
		defer rows.Close()

		var comments []models.Comment
		for rows.Next() {
			var comment models.Comment
			if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse comment"})
				return
			}
			comments = append(comments, comment)
		}

		c.JSON(http.StatusOK, comments)
	}
}

func DeleteCommentHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("commentID")

		_, err := db.Exec("DELETE FROM comments WHERE id = $1", commentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
	}
}
