package router

import (
	"github.com/InatoInato/car_marketplace_go.git/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, userHandler *handler.UserHandler){
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
}