package main

import (
	"fibo_go_server/config"
	"fibo_go_server/internal/routes"
	"log"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	router := gin.Default()
	routes.SetupRoutes(router, cfg)

	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
