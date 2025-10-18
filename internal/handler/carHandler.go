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
			Model: "GTR R35",
			Year: 2020,
			Color: "Black",
			EngineCapacity: 3.8,
		}
	return c.JSON(car)
}