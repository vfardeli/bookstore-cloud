package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")

	// 1. Verify order exists and is pending
	orderResp, err := http.Get(fmt.Sprintf("%s/orders", orderServiceURL))
	if err != nil || orderResp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot verify order"})
		return
	}

	// (In a real app, you'd GET /orders/:id and check status == "PENDING")

	// 2. Simulate payment gateway processing
	time.Sleep(2 * time.Second) // pretend to talk to Stripe/PayPal
	payment.Status = "SUCCESS"

	// 3. Update order status to PAID
	updatePayload, _ := json.Marshal(map[string]string{"status": "PAID"})
	req, _ := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/orders/%d", orderServiceURL, payment.OrderID),
		bytes.NewBuffer(updatePayload),
	)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment processed but failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment processed successfully",
		"payment": payment,
	})
}
