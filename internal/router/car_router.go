package router

import (
	"github.com/InatoInato/car_marketplace_go.git/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){
	app.Get("/", handler.Home)
	app.Get("/car", handler.Car)
}