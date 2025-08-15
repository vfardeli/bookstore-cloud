package routes

import (
	"bookstore-cloud/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/users/register", handlers.RegisterUser)
	r.POST("/users/login", handlers.LoginUser)

	r.GET("/books", handlers.ListBooks)
	r.GET("/books/:id", handlers.GetBook)

	r.POST("/orders", handlers.CreateOrder)

	return r
}
