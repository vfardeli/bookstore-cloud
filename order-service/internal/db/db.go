package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	DB_HOST_ENV_VAR     = "DB_HOST"
	DB_USER_ENV_VAR     = "DB_USER"
	DB_PASSWORD_ENV_VAR = "DB_PASSWORD"
	DB_NAME_ENV_VAR     = "DB_NAME"
	DB_PORT_ENV_VAR     = "DB_PORT"
)

// Connect to the database using the environment variables to grab all the sensitive database inputs.
func Connect() {
	// Load DB parameters from environment variables.
	godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv(DB_HOST_ENV_VAR),
		os.Getenv(DB_USER_ENV_VAR),
		os.Getenv(DB_PASSWORD_ENV_VAR),
		os.Getenv(DB_NAME_ENV_VAR),
		os.Getenv(DB_PORT_ENV_VAR),
	)

	var database *gorm.DB
	var err error

	for i := 0; i < 10; i++ { // retry up to 10 times
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to DB")
			break
		} else {
			log.Printf("⏳ DB not ready (%v). Retrying in 3s...\n", err)
			time.Sleep(3 * time.Second)
		}
	}

	DB = database
	fmt.Println("Order DB connected!")
}
