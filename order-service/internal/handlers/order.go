package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"order-service/internal/db"
	"order-service/internal/models"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate User
	userServiceURL := os.Getenv("USER_SERVICE_URL") + "/login"
	userPayload, _ := json.Marshal(map[string]string{
		"username": "alice", // Normally would be from auth token
		"password": "pass123",
	})
	resp, err := http.Post(userServiceURL, "application/json", bytes.NewBuffer(userPayload))
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	// Validate Book
	bookServiceURL := os.Getenv("BOOK_SERVICE_URL") + "/books/" + strconv.Itoa(int(order.BookID))
	resp, err = http.Get(bookServiceURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book"})
		return
	}

	order.Status = "PENDING"
	db.DB.Create(&order)
	c.JSON(http.StatusOK, order)
}

func GetOrders(c *gin.Context) {
	var orders []models.Order
	db.DB.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
