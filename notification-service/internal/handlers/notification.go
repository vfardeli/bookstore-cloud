package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"notification-service/internal/models"

	"github.com/gin-gonic/gin"
)

func SendNotification(c *gin.Context) {
	var notif models.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from := os.Getenv("EMAIL_FROM")
	if notif.Type == "" {
		notif.Type = "email"
	}

	// Simulate sending delay
	time.Sleep(1 * time.Second)

	// Log sending (mock sending)
	fmt.Printf("[Notification] Sent %s to %s from %s\nSubject: %s\nBody: %s\n\n",
		notif.Type, notif.To, from, notif.Subject, notif.Body)

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification sent successfully",
		"data":    notif,
	})
}
