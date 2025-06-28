package middlewares

import (
	"go-api/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateRequest[T any]() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req T

		if err := ctx.BodyParser(&req); err != nil {
			return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "error", "Invalid JSON format", err)
		}

		if normalizer, ok := any(&req).(interface{ Normalize() }); ok {
			normalizer.Normalize()
		}

		if err := validate.Struct(&req); err != nil {
			return helpers.ValidationErrorResponse[T](ctx, err)
		}

		ctx.Locals("validated", req)
		return ctx.Next()
	}
}
