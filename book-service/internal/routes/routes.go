package routes

import (
	"book-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRouter() (*gin.Engine, func()) {
	shutdown := handlers.InitTracer("book-service")

	r := gin.Default()

	r.Use(otelgin.Middleware("book-service"))
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

	r.POST("/books", handlers.AddBook)
	r.GET("/books", handlers.ListBooks)
	r.GET("/books/:id", handlers.GetBook)
	return r, shutdown
}
