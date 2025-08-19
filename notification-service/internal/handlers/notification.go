package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"notification-service/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumePayments() {
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

	_ = ch.ExchangeDeclare(
		"payment.completed", // exchange name
		"fanout",            // type
		true,                // durable
		false,
		false,
		false,
		nil,
	)
	q, _ := ch.QueueDeclare("", false, true, true, false, nil)
	_ = ch.QueueBind(
		q.Name,
		"",
		"payment.completed",
		false,
		nil,
	)

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	log.Println("Notification Service waiting for payment.completed events...")

	for d := range msgs {
		var payment models.Payment
		json.Unmarshal(d.Body, &payment)

		log.Printf("Sending notification: Payment %s for Order %d\n", payment.Status, payment.OrderID)
	}
}

// func SendNotification(c *gin.Context) {
// 	var notif models.Notification
// 	if err := c.ShouldBindJSON(&notif); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	from := os.Getenv("EMAIL_FROM")

// 	// Simulate sending delay
// 	time.Sleep(1 * time.Second)

// 	// Log sending (mock sending)
// 	fmt.Printf("[Notification] Sent email to %s from %s\nSubject: %s\nBody: %s\n\n",
// 		notif.To, from, notif.Subject, notif.Body)

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Notification sent successfully",
// 		"data":    notif,
// 	})
// }
