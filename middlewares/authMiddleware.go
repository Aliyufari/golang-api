package middlewares

import (
	"errors"
	"fmt"
	"go-api/config"
	"go-api/helpers"
	"go-api/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Authenticate(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid token", nil)
	}
	userToken := authHeader[7:]

	token, err := jwt.Parse(userToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", err)
	}

	user, err := getUserFromToken(token)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not found", err)
	}

	ctx.Locals("user", user)
	return ctx.Next()
}

func getUserFromToken(token *jwt.Token) (*models.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, errors.New("invalid UUID format in token")
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func HasRole(requiredRole string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*models.User)
		if !ok || user.ID == uuid.Nil {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Access denied", nil)
		}

		var role models.Role
		if err := config.DB.First(&role, "id = ?", user.ID).Error; err != nil {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Role not found", err)
		}

		if role.Name != requiredRole {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Insufficient role", nil)
		}

		return ctx.Next()
	}
}

func HasAnyRole(allowedRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*models.User)
		if !ok || user.ID == uuid.Nil {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Access denied", nil)
		}

		var role models.Role
		if err := config.DB.First(&role, "id = ?", user.ID).Error; err != nil {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Role not found", err)
		}

		for _, r := range allowedRoles {
			if role.Name == r {
				return ctx.Next()
			}
		}

		return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "No matching role found", nil)
	}
}

func HasPermission(permissionName string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*models.User)
		if !ok || user.ID == uuid.Nil {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Access denied", nil)
		}

		var role models.Role
		if err := config.DB.Preload("Permissions").First(&role, "id = ?", user.ID).Error; err != nil {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Role not found", err)
		}

		for _, perm := range role.Permissions {
			if perm.Name == permissionName {
				return ctx.Next()
			}
		}

		return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Permission denied", nil)
	}
}

// func HasAnyPermission(allowedPermissions ...string) fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		user, ok := ctx.Locals("user").(*models.User)
// 		if !ok || user.RoleID == uuid.Nil {
// 			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Access denied", nil)
// 		}

// 		var role models.Role
// 		if err := config.DB.Preload("Permissions").First(&role, "id = ?", user.RoleID).Error; err != nil {
// 			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "Role not found", err)
// 		}

// 		for _, perm := range role.Permissions {
// 			for _, allowed := range allowedPermissions {
// 				if perm.Name == allowed {
// 					return ctx.Next()
// 				}
// 			}
// 		}

// 		return helpers.ErrorResponse(ctx, fiber.StatusForbidden, "FORBIDDEN", "No required permission found", nil)
// 	}
// }
