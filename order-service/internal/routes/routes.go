package routes

import (
	"order-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/orders", handlers.CreateOrder)
	r.GET("/orders", handlers.GetOrders)
	return r
}
