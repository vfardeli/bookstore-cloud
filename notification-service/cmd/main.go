package main

import (
	"notification-service/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	handlers.ConsumePayments()
}
