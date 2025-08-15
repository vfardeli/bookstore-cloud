package routes

import (
	"notification-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/notifications", handlers.SendNotification)
	return r
}
