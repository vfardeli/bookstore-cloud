package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"notification-service/internal/models"
	"notification-service/internal/utils"

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
		reqID := d.CorrelationId
		var payment models.Payment
		json.Unmarshal(d.Body, &payment)

		utils.SendLog("notification-service", reqID, "info", fmt.Sprintf("Sending notification: Payment %s for Order %d\n", payment.Status, payment.OrderID), nil)
	}
}
