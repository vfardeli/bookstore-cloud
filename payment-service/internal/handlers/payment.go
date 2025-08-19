package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"payment-service/internal/models"
	"payment-service/internal/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func publishPaymentCompleted(payment models.Payment, requestID string) {
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
		ContentType:   "application/json",
		Body:          body,
		CorrelationId: requestID,
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
		reqID := d.CorrelationId
		var order models.Order
		json.Unmarshal(d.Body, &order)

		utils.SendLog("payment-service", reqID, "info", fmt.Sprintf("Processing payment for order: %d", order.ID), nil)
		time.Sleep(2 * time.Second) // simulate payment delay

		payment := models.Payment{OrderID: order.ID, Status: "COMPLETED"}
		publishPaymentCompleted(payment, reqID)
	}
}
