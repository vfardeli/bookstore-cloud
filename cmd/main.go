package main

import (
	"bookstore-cloud/internal/db"
	"bookstore-cloud/internal/models"
	"bookstore-cloud/internal/routes"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Order{})

	r := routes.SetupRouter()
	r.Run(":3000")
}
