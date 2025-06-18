package main

import(
	"log"
	"os" 
	"go-api/config"
	"go-api/helpers"
	"go-api/routes"
	"github.com/gofiber/fiber/v2"
)

func init(){
	helpers.LoadEnv()
	config.ConnectDB()
}

func main()  {
	app := fiber.New()

    routes.AuthRoutes(app)

	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}

    log.Fatal(app.Listen(":" + port))
}