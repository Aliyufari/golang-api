package routes

import (
	"go-api/controllers"
	"go-api/middlewares"
	"go-api/requests"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {

	app.Post("/api/register", middlewares.ValidateRequest[requests.CreateUserRequest](), controllers.Register)
	app.Post("/api/login", middlewares.ValidateRequest[requests.LoginUserRequest](), controllers.Login)

}
