package routes

import (
	"fibo_go_server/config"
	"fibo_go_server/internal/db"
	"fibo_go_server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// User routes
	router.POST("/users/register", handlers.RegisterUserHandler(database))
	router.POST("/users/login", handlers.LoginUserHandler(database))

	// Post routes
	router.POST("/posts", handlers.CreatePostHandler(database))

	// Comment routes
	router.POST("/comments", handlers.AddCommentHandler(database))
	router.GET("/posts/:postID/comments", handlers.GetCommentsHandler(database))
	router.DELETE("/comments/:commentID", handlers.DeleteCommentHandler(database))
}
