package main

import (
	"order-service/internal/db"
	"order-service/internal/models"
	"order-service/internal/routes"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&models.Order{})
	r := routes.SetupRouter()
	r.Run(":8003")
}
