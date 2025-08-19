package main

import (
	"payment-service/internal/handlers"

	"github.com/joho/godotenv"
)

// func main() {
// 	godotenv.Load()
// 	r := routes.SetupRouter()
// 	r.Run(":8004")
// }

func main() {
	godotenv.Load()
	handlers.ConsumeOrders()
}
