package main

import (
	"api-gateway/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Routes → services
	// Users
	r.POST("/register", handlers.ProxyToRegisterUserService)
	r.POST("/login", handlers.ProxyToLoginService)

	// Books
	r.Any("/books", handlers.ProxyToBookService)
	r.Any("/books/*path", handlers.ProxyToBookService)

	// Orders
	r.Any("/orders", handlers.ProxyToOrderService)
	r.Any("/orders/*path", handlers.ProxyToOrderService)

	log.Println("API Gateway running on :8080")
	r.Run(":8080")
}
