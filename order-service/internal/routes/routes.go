package routes

import (
	"order-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRouter() (*gin.Engine, func()) {
	shutdown := handlers.InitTracer("order-service")

	r := gin.Default()

	r.Use(otelgin.Middleware("order-service"))
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

	r.POST("/orders", handlers.CreateOrder)
	r.GET("/orders", handlers.GetOrders)
	return r, shutdown
}
