package routes

import (
	"go-api/controllers"
	"go-api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RoleRoutes(app *fiber.App) {
	app.Get(
		"/api/roles",
		middlewares.Authenticate,
		middlewares.HasRole("admin"),
		controllers.GetRoles,
	)
}
