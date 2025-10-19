package main

import (
	"log"
	"os"

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

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Get server port from environment variables
	port := os.Getenv("SERVER_PORT")
	if port == ""{
		port = "8080"
	}

	app := fiber.New()
	router.Setup(app, userHandler)

	if err := app.Listen(":" + port); err != nil{
		panic(err)
	}
}