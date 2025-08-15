package routes

import (
	"payment-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/payments", handlers.ProcessPayment)
	return r
}
