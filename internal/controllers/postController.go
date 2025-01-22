package controllers

import (
	"database/sql"
	"net/http"

	"fibo_go_server/internal/models"

	"github.com/gin-gonic/gin"
)

func CreatePost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var post models.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO fiboblog.posts (user_id, category_id, title, content) 
		          VALUES ($1, $2, $3, $4) RETURNING id`
		err := db.QueryRow(query, post.UserID, post.CategoryID, post.Title, post.Content).Scan(&post.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}
func GetPostsList(db *sql.DB) gin.HandlerFunc {
	type PostListType struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Username     string `json:"username"`
		CategoryName string `json:"category_name"`
		IsApproved   bool   `json:"is_approved"`
		ViewCount    int    `json:"view_count"`
		CreatedAt    string `json:"created_at"`
	}
	return func(c *gin.Context) {
		var posts []PostListType

		query := `
			SELECT 
				posts.id, 
				posts.title, 
				users.username, 
				categories.name AS category_name, 
				posts.is_approved, 
				posts.view_count, 
				posts.created_at 
			FROM 
				fiboblog.posts
			LEFT JOIN 
				fiboblog.users ON posts.user_id = users.id
			LEFT JOIN 
				fiboblog.categories ON posts.category_id = categories.id
		`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var post PostListType
			if err := rows.Scan(
				&post.ID,
				&post.Title,
				&post.Username,
				&post.CategoryName,
				&post.IsApproved,
				&post.ViewCount,
				&post.CreatedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			posts = append(posts, post)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, posts)
	}
}

type PostType struct {
	ID           int           `json:"id"`
	UserID       int           `json:"user_id"`
	UserName     string        `json:"user_name"` // User's username
	CategoryID   *int          `json:"category_id"`
	CategoryName string        `json:"category_name"` // Category name
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	IsApproved   bool          `json:"is_approved"`
	ViewCount    int           `json:"view_count"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	LikeCount    int           `json:"like_count"`
	Comments     []CommentType `json:"comments"`
}

type CommentType struct {
	ID              int           `json:"id"`
	PostID          int           `json:"post_id"`
	UserID          *int          `json:"user_id"`
	UserName        string        `json:"user_name"` // User's username
	ParentCommentID *int          `json:"parent_comment_id"`
	Content         string        `json:"content"`
	CreatedAt       string        `json:"created_at"`
	Children        []CommentType `json:"children"`
}

func GetPostDetails(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse post_id from query parameters
		postID := c.DefaultQuery("post_id", "")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
			return
		}

		// Fetch post details including the user username and category name
		var post PostType
		err := db.QueryRow(`
			SELECT p.id, p.user_id, u.username, p.category_id, c.name, p.title, p.content, p.is_approved, p.view_count, p.created_at, p.updated_at
			FROM fiboblog.posts p
			JOIN fiboblog.users u ON p.user_id = u.id
			JOIN fiboblog.categories c ON p.category_id = c.id
			WHERE p.id = $1`, postID).Scan(
			&post.ID, &post.UserID, &post.UserName, &post.CategoryID, &post.CategoryName, &post.Title, &post.Content,
			&post.IsApproved, &post.ViewCount, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		// Fetch like count
		err = db.QueryRow(`SELECT COUNT(*) FROM fiboblog.likes WHERE post_id = $1`, postID).Scan(&post.LikeCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch likes"})
			return
		}

		// Fetch comments including the user username
		comments := []CommentType{}
		rows, err := db.Query(`
			SELECT c.id, c.post_id, c.user_id, u.username, c.parent_comment_id, c.content, c.created_at
			FROM fiboblog.comments c
			JOIN fiboblog.users u ON c.user_id = u.id
			WHERE c.post_id = $1`, postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var comment CommentType
			if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.UserName, &comment.ParentCommentID, &comment.Content, &comment.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading comments"})
				return
			}
			comments = append(comments, comment)
		}
		post.Comments = comments

		// Return the result as JSON
		c.JSON(http.StatusOK, post)
	}
}
func GetCategories(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category

		rows, err := db.Query("SELECT id, name, description FROM fiboblog.categories")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var category models.Category
			if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			categories = append(categories, category)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}
