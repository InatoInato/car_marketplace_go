package main

import (
	"github.com/InatoInato/car_marketplace_go.git/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	router.Setup(app)

	if err := app.Listen(":8080"); err != nil{
		panic(err)
	}
}