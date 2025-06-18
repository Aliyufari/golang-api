package controllers

import "github.com/gofiber/fiber/v2"

func Register(ctx *fiber.Ctx) error {
	return ctx.SendString("Register endpoint hit")
}

func Login(ctx *fiber.Ctx) error {
	return ctx.SendString("Login endpoint hit")
}