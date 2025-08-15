package main

import (
	"notification-service/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := routes.SetupRouter()
	r.Run(":8000")
}
