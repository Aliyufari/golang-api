package main

import (
	"go-api/config"
	"go-api/helpers"
	"go-api/routes"
	"go-api/seeders"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func init() {
	helpers.LoadEnv()
	config.ConnectDB()

	//Seed Roles & Permissions
	if err := seeders.SeedRolesAndPermissions(config.DB); err != nil {
		panic(err)
	}
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
