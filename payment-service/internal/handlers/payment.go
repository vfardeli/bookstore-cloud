package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"payment-service/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func publishPaymentCompleted(payment models.Payment) {
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")

	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ { // retry up to 10 times
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@rabbitmq:5672/", rabbitmqUser, rabbitmqPassword))
		if err == nil {
			log.Println("Connected to RabbitMQ")
			break
		} else {
			log.Printf("⏳ RabbitMQ not ready (%v). Retrying in 3s...\n", err)
			time.Sleep(3 * time.Second)
		}
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	// q, _ := ch.QueueDeclare("payment.completed", false, false, false, false, nil)
	_ = ch.ExchangeDeclare(
		"payment.completed", // exchange name
		"fanout",            // type
		true,                // durable
		false,
		false,
		false,
		nil,
	)

	body, _ := json.Marshal(payment)
	ch.Publish("payment.completed", "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	log.Println("Published payment.completed event:", string(body))
}

func ConsumeOrders() {
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")

	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ { // retry up to 10 times
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@rabbitmq:5672/", rabbitmqUser, rabbitmqPassword))
		if err == nil {
			log.Println("Connected to RabbitMQ")
			break
		} else {
			log.Printf("⏳ RabbitMQ not ready (%v). Retrying in 3s...\n", err)
			time.Sleep(3 * time.Second)
		}
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	q, _ := ch.QueueDeclare("order.created", false, false, false, false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	log.Println("Payment Service waiting for order.created events...")

	for d := range msgs {
		var order models.Order
		json.Unmarshal(d.Body, &order)

		log.Println("Processing payment for order:", order.ID)
		time.Sleep(2 * time.Second) // simulate payment delay

		payment := models.Payment{OrderID: order.ID, Status: "COMPLETED"}
		publishPaymentCompleted(payment)
	}
}

// func ProcessPayment(c *gin.Context) {
// 	var payment models.Payment
// 	if err := c.ShouldBindJSON(&payment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// orderServiceURL := os.Getenv("ORDER_SERVICE_URL")

// 	// 1. Verify order exists and is pending
// 	// orderResp, err := http.Get(fmt.Sprintf("%s/orders", orderServiceURL))
// 	// if err != nil || orderResp.StatusCode != http.StatusOK {
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot verify order"})
// 	// 	return
// 	// }

// 	// (In a real app, you'd GET /orders/:id and check status == "PENDING")

// 	// 2. Fake Payment Processing
// 	time.Sleep(2 * time.Second) // pretend to talk to Stripe/PayPal
// 	payment.Status = "SUCCESS"
// 	log.Printf("Payment processed for Order %d, Amount %.2f", payment.OrderID, payment.Amount)

// 	// 3. Update order status to PAID
// 	// updatePayload, _ := json.Marshal(map[string]string{"status": "PAID"})
// 	// req, _ := http.NewRequest(http.MethodPut,
// 	// 	fmt.Sprintf("%s/orders/%d", orderServiceURL, payment.OrderID),
// 	// 	bytes.NewBuffer(updatePayload),
// 	// )
// 	// req.Header.Set("Content-Type", "application/json")
// 	// client := &http.Client{}
// 	// resp, err := client.Do(req)
// 	// if err != nil || resp.StatusCode != http.StatusOK {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment processed but failed to update order"})
// 	// 	return
// 	// }

// 	// Call Notification Service
// 	notificationPayload := map[string]string{
// 		"to":      "recipient@example.com",
// 		"subject": "Payment Successful",
// 		"body":    fmt.Sprintf("Your payment for order %d was successful!", payment.OrderID),
// 	}
// 	payloadBytes, _ := json.Marshal(notificationPayload)

// 	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")

// 	_, err := http.Post(fmt.Sprintf("%s/notifications", notificationServiceURL), "application/json", bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		log.Println("Error calling notification service:", err)
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Payment processed successfully",
// 		"payment": payment,
// 	})
// }
