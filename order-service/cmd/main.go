package main

import (
	"order-service/internal/db"
	"order-service/internal/handlers"
	"order-service/internal/models"
	"order-service/internal/routes"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&models.Order{})

	// Run RabbitMQ consumer in a separate goroutine
	go handlers.ConsumePayments()

	r := routes.SetupRouter()
	r.Run(":8003")
}
