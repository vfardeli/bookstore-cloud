package routes

import (
	"user-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)
	return r
}
