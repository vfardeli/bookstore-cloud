package main

import (
	"book-service/internal/db"
	"book-service/internal/models"
	"book-service/internal/routes"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&models.Book{})
	r := routes.SetupRouter()
	r.Run(":8002")
}
