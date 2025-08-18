package main

import (
	"io"
	"log"
	"net/http"

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
	r.Any("/login", func(c *gin.Context) {
		forwardRequest(c, "http://user-service:8001/login")
	})
	r.Any("/register", func(c *gin.Context) {
		forwardRequest(c, "http://user-service:8001/register")
	})
	r.Any("/books/*path", func(c *gin.Context) {
		forwardRequest(c, "http://book-service:8002/books"+c.Param("path"))
	})
	r.Any("/orders/*path", func(c *gin.Context) {
		forwardRequest(c, "http://order-service:8003/orders"+c.Param("path"))
	})
	r.Any("/payments/*path", func(c *gin.Context) {
		forwardRequest(c, "http://payment-service:8004/payments"+c.Param("path"))
	})
	r.Any("/notifications/*path", func(c *gin.Context) {
		forwardRequest(c, "http://notification-service:8005/notifications"+c.Param("path"))
	})

	log.Println("API Gateway running on :8080")
	r.Run(":8080")
}
