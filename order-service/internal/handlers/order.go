package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-service/internal/db"
	"order-service/internal/models"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishOrderCreated(order models.Order) {
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@rabbitmq:5672/", rabbitmqUser, rabbitmqPassword))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	q, _ := ch.QueueDeclare(
		"order.created", // queue name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	body, _ := json.Marshal(order)
	ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	log.Println("Published order.created event:", string(body))
}

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

	msgs, _ := ch.Consume(q.Name, "", true, true, false, false, nil)

	log.Println("Order Service waiting for payment.completed events...")

	for d := range msgs {
		var payment models.Payment
		json.Unmarshal(d.Body, &payment)

		var order models.Order
		db.DB.First(&order, payment.OrderID)

		order.Status = "PAID"
		db.DB.Save(order)

		log.Printf("Payment paid for Order %d\n", order.ID)
	}
}

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

	order.Amount = float64(order.Quantity) * book.Price
	order.Status = "PENDING"
	db.DB.Create(&order)
	log.Println("Order created: ", order)

	// Call Payment Service
	// paymentPayload := map[string]interface{}{
	// 	"order_id": order.ID,
	// 	"method":   "CREDIT_CARD",
	// 	"amount":   float64(order.Quantity) * book.Price,
	// }
	// payloadBytes, _ := json.Marshal(paymentPayload)

	// paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")

	// paymentResp, err := http.Post(fmt.Sprintf("%s/payments", paymentServiceURL), "application/json", bytes.NewBuffer(payloadBytes))
	// if err != nil || paymentResp.StatusCode != http.StatusOK {
	// 	log.Println("Error calling payment service:", err)
	// 	order.Status = "PAYMENT_FAILED"
	// 	c.JSON(http.StatusInternalServerError, order)
	// 	return
	// }
	// defer resp.Body.Close()

	// order.Status = "PAID"
	// db.DB.Save(order)
	// c.JSON(http.StatusOK, order)

	publishOrderCreated(order)
	c.JSON(http.StatusCreated, order)
}

func GetOrders(c *gin.Context) {
	var orders []models.Order
	db.DB.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
