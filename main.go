package main

import(
	"log"
	"go-api/config"
	"github.com/gofiber/fiber/v2"
)

func main()  {
	config.Connect()

	app := fiber.New()

    app.Get("/", func (ctx *fiber.Ctx) error {
        return ctx.SendString("Hello world!")
    })

    log.Fatal(app.Listen(":8000"))
}