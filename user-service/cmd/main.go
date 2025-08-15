package main

import (
	"user-service/internal/db"
	"user-service/internal/models"
	"user-service/internal/routes"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&models.User{})
	r := routes.SetupRouter()
	r.Run(":8000")
}
