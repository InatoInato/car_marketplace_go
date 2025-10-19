package handler

import (
	"github.com/InatoInato/car_marketplace_go.git/internal/model"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Car(c *fiber.Ctx) error{
		car := model.Car{
			Make: "Nissan",
			CarModel: "Primera P10",
			Year: 1995,
			Color: "Silver",
			EngineCapacity: 2.0,
			Transmission: "Manual",
			IsRunning: true,
			UserID: 1,
		}
	return c.JSON(car)
}