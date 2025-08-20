package routes

import (
	"user-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRouter() (*gin.Engine, func()) {
	shutdown := handlers.InitTracer("user-service")

	r := gin.Default()

	r.Use(otelgin.Middleware("user-service"))
	// Middleware: ensure every request has a Request ID
	r.Use(func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("RequestID", reqID)
		c.Writer.Header().Set("X-Request-ID", reqID)
		c.Next()
	})

	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)

	return r, shutdown
}
