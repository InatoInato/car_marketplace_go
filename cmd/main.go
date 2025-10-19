package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/InatoInato/car_marketplace_go.git/internal/config"
	"github.com/InatoInato/car_marketplace_go.git/internal/handler"
	"github.com/InatoInato/car_marketplace_go.git/internal/repository"
	"github.com/InatoInato/car_marketplace_go.git/internal/router"
	"github.com/InatoInato/car_marketplace_go.git/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(); if err != nil{
		log.Println("Error loading .env file")
	}

	// Connect to the database
	config.ConnectDB()
	cache := config.ConnectRedis()



	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	carRepo := repository.NewCarRepository(config.DB)
	carService := service.NewCarService(carRepo, userRepo, cache)
	carHandler := handler.NewCarHandler(carService)

	// Get server port from environment variables
	port := os.Getenv("SERVER_PORT")
	if port == ""{
		port = "8080"
	}

	app := fiber.New()
	router.Setup(app, userHandler, carHandler)

	go func() {
		log.Printf("Server is running on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚠️  Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("⚠️  Server forced to shutdown: %v", err)
	}

	// Close Redis
	if err := cache.Close(); err != nil {
		log.Printf("⚠️  Error closing Redis: %v", err)
	} else {
		log.Println("✅ Redis connection closed")
	}

	// Close PostgreSQL
	db, _ := config.DB.DB()
	if err := db.Close(); err != nil {
		log.Printf("⚠️  Error closing PostgreSQL: %v", err)
	} else {
		log.Println("✅ PostgreSQL connection closed")
	}

	log.Println("👋 Graceful shutdown complete")
}