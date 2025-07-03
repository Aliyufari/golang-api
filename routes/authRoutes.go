package routes

import (
	"go-api/controllers"
	"go-api/dtos"
	"go-api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/api/register", middlewares.ValidateRequest[dtos.CreateUserRequest](), controllers.Register)
	app.Post("/api/login", middlewares.ValidateRequest[dtos.LoginUserRequest](), controllers.Login)
	app.Get("/api/me", middlewares.Authenticate, controllers.Me)
}
