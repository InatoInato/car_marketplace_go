package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	env := os.Getenv("APP_ENV") // dev / prod

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	// ---------- Logger GORM ----------
	var logLevel logger.LogLevel
	var ignoreNotFound bool

	if env == "prod" {
		logLevel = logger.Warn
		ignoreNotFound = true
	} else {
		logLevel = logger.Info
		ignoreNotFound = false
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: ignoreNotFound,
			Colorful:                  true,
		},
	)

	// ---------- Connecting to db ----------
	var db *gorm.DB
	var err error
	maxAttempts := 10

	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err == nil {
			break
		}
		log.Printf("❌ Failed to connect to DB, attempt %d/%d. Retrying in 2s...", i, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("💀 Failed to connect to database after multiple attempts:", err)
	}

	log.Println("✅ Database connected successfully!")

	// ---------- Auto migrations ----------
	err = db.AutoMigrate(
		&model.User{},
		&model.Car{},
		&model.CarPhoto{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	DB = db
}
