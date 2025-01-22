package main

import (
	"fibo_go_server/config"
	"fibo_go_server/internal/routes"
)

func main() {
	config.LoadConfig()
	r := routes.SetupRouter()
	r.Run(":8080")
}
