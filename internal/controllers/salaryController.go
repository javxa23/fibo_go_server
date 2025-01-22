package controllers

import (
	"database/sql"
	"net/http"

	"fibo_go_server/internal/models"

	"github.com/gin-gonic/gin"
)

// CalculateSalary handles the process of calculating and adding salary to a user
func CalculateSalary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var salary models.Salary
		if err := c.ShouldBindJSON(&salary); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Simple calculation for salary, this can be more complex based on the reputation points
		query := `INSERT INTO salaries (user_id, amount, month, year, created_at) 
		          VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
		err := db.QueryRow(query, salary.UserID, salary.Amount, salary.Month, salary.Year).Scan(&salary.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, salary)
	}
}
