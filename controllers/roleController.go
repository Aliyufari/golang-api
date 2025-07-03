package controllers

import (
	"go-api/config"
	"go-api/helpers"
	"go-api/models"

	"github.com/gofiber/fiber/v2"
)

func GetRoles(ctx *fiber.Ctx) error {
	var roles []models.Role

	if err := config.DB.Preload("Users", "Permissions").Find(&roles).Error; err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "ERROR", "Failed to fetch roles", err)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "OK", "Roles", "roles", roles)
}
