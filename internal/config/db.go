package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	var db *gorm.DB
	var err error


	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database, attempt %d/%d. Retrying in 2s...", i, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after multiple attempts:", err)
	}

	log.Println("DB Connected")

	err = db.AutoMigrate(
		&model.User{},
		&model.Car{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
}
