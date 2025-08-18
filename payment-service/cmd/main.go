package main

import (
	"payment-service/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := routes.SetupRouter()
	r.Run(":8004")
}
