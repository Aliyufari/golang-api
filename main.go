package main

import (
	"go-api/config"
	"go-api/helpers"
	"go-api/routes"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func init() {
	helpers.LoadEnv()
	config.ConnectDB()
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	routes.AuthRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen(":" + port))
}
