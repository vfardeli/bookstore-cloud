package main

import (
	"io"
	"log"
	"net/http"

	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

func forwardRequest(c *gin.Context, target string) {
	resp, err := http.NewRequest(c.Request.Method, target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Copy headers
	for k, v := range c.Request.Header {
		for _, val := range v {
			resp.Header.Add(k, val)
		}
	}

	client := &http.Client{}
	response, err := client.Do(resp)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	c.Data(response.StatusCode, response.Header.Get("Content-Type"), body)
}

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
