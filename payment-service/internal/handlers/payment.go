package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"payment-service/internal/models"

	"github.com/gin-gonic/gin"
)

func ProcessPayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// orderServiceURL := os.Getenv("ORDER_SERVICE_URL")

	// 1. Verify order exists and is pending
	// orderResp, err := http.Get(fmt.Sprintf("%s/orders", orderServiceURL))
	// if err != nil || orderResp.StatusCode != http.StatusOK {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot verify order"})
	// 	return
	// }

	// (In a real app, you'd GET /orders/:id and check status == "PENDING")

	// 2. Fake Payment Processing
	time.Sleep(2 * time.Second) // pretend to talk to Stripe/PayPal
	payment.Status = "SUCCESS"
	log.Printf("Payment processed for Order %d, Amount %.2f", payment.OrderID, payment.Amount)

	// 3. Update order status to PAID
	// updatePayload, _ := json.Marshal(map[string]string{"status": "PAID"})
	// req, _ := http.NewRequest(http.MethodPut,
	// 	fmt.Sprintf("%s/orders/%d", orderServiceURL, payment.OrderID),
	// 	bytes.NewBuffer(updatePayload),
	// )
	// req.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil || resp.StatusCode != http.StatusOK {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment processed but failed to update order"})
	// 	return
	// }

	// Call Notification Service
	notificationPayload := map[string]string{
		"to":      "recipient@example.com",
		"subject": "Payment Successful",
		"body":    fmt.Sprintf("Your payment for order %d was successful!", payment.OrderID),
	}
	payloadBytes, _ := json.Marshal(notificationPayload)

	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")

	_, err := http.Post(fmt.Sprintf("%s/notifications", notificationServiceURL), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Error calling notification service:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment processed successfully",
		"payment": payment,
	})
}
