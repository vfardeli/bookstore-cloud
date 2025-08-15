package routes

import (
	"book-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/books", handlers.AddBook)
	r.GET("/books", handlers.ListBooks)
	r.GET("/books/:id", handlers.GetBook)
	return r
}
