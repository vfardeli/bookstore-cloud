package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	// userServiceURL := os.Getenv("USER_SERVICE_URL") + "/login"
	// userPayload, _ := json.Marshal(map[string]string{
	// 	"username": "alice", // Normally would be from auth token
	// 	"password": "pass123",
	// })
	// resp, err := http.Post(userServiceURL, "application/json", bytes.NewBuffer(userPayload))
	// if err != nil || resp.StatusCode != http.StatusOK {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
	// 	return
	// }

	// Validate Book
	bookServiceURL := os.Getenv("BOOK_SERVICE_URL") + "/books/" + strconv.Itoa(int(order.BookID))
	resp, err := http.Get(bookServiceURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book"})
		return
	}
	defer resp.Body.Close()

	var book models.Book
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order.Status = "PENDING"
	db.DB.Create(&order)
	log.Println("Order created: ", order)

	// Call Payment Service
	paymentPayload := map[string]interface{}{
		"order_id": order.ID,
		"method":   "CREDIT_CARD",
		"amount":   float64(order.Quantity) * book.Price,
	}
	payloadBytes, _ := json.Marshal(paymentPayload)

	paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")

	paymentResp, err := http.Post(fmt.Sprintf("%s/payments", paymentServiceURL), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil || paymentResp.StatusCode != http.StatusOK {
		log.Println("Error calling payment service:", err)
		order.Status = "PAYMENT_FAILED"
		c.JSON(http.StatusInternalServerError, order)
		return
	}
	defer resp.Body.Close()

	order.Status = "PAID"
	db.DB.Save(order)
	c.JSON(http.StatusOK, order)
}

func GetOrders(c *gin.Context) {
	var orders []models.Order
	db.DB.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
