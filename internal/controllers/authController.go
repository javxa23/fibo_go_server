package controllers

import (
	"database/sql"
	"fibo_go_server/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM fiboblog.users WHERE username = $1", user.Username).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Username already taken"})
			return
		}

		err = db.QueryRow("SELECT COUNT(*) FROM fiboblog.users WHERE email = $1", user.Email).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email already taken"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		_, err = db.Exec("INSERT INTO fiboblog.users (username, email, password) VALUES ($1, $2, $3)",
			user.Username, user.Email, string(hashedPassword))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

func LoginUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var hashedPassword string
		var userID int
		var username string
		err := db.QueryRow("SELECT id, username, password FROM fiboblog.users WHERE email = $1", credentials.Email).Scan(&userID, &username, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Login successful",
			"user_id":  userID,
			"username": username,
		})
	}
}
