package helpers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Error(c *fiber.Ctx, status int, message string, errors interface{}) error {
	log.Printf("Error: %v", errors)

	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": message,
	})
}

func Success(c *fiber.Ctx, status int, message string, dataKey string, data interface{}) error {
	response := fiber.Map{
		"status":  status,
		"message": message,
	}

	if data != nil {
		response[dataKey] = data
	}

	return c.Status(status).JSON(response)
}
