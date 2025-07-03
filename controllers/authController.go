package controllers

import (
	"go-api/config"
	"go-api/dtos"
	"go-api/helpers"
	"go-api/models"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) (err error) { // <== named return value
	req := ctx.Locals("validated").(dtos.CreateUserRequest)

	var user models.User
	if err = config.DB.Where("email = ?", req.Email).First(&user).Error; err == nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", "Email already exists", nil)
	}

	userAvatar, fullPath, err := helpers.SaveAvatar(ctx, "avatar", 2*1024*1024, "./public/avatars")
	if err != nil {
		log.Println("Avatar upload error:", err)
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", "Could not save avatar", nil)
	}

	defer func() {
		if err != nil {
			_ = os.Remove(fullPath)
		}
	}()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing failed:", err)
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "ERROR", "Could not hash password!", nil)
	}

	var role models.Role
	if err = config.DB.Where("name = ?", "user").First(&role).Error; err == nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, "NOT FOUND", "Role Not Found", nil)
	}

	newUser := models.User{
		Avatar:   userAvatar,
		Name:     req.Name,
		Email:    req.Email,
		Dob:      req.Dob,
		Gender:   req.Gender,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	if err = config.DB.Create(&newUser).Error; err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "INTERNAL SERVER ERROR", "Could not create user", nil)
	}

	response := dtos.UserResponse{
		ID:     user.ID,
		Avatar: userAvatar,
		Name:   req.Name,
		Email:  req.Email,
		Dob:    req.Dob,
		Gender: req.Gender,
		Role:   role.Name,
	}

	// prevent defer from deleting the file
	err = nil
	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "CREATED", "User created successfully", "user", response)
}

func Login(ctx *fiber.Ctx) error {
	req := ctx.Locals("validated").(dtos.LoginUserRequest)

	var user models.User
	result := config.DB.First(&user, "email = ?", req.Email)
	if result.Error != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", "Invalid Credentials", nil)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", "Invalid Credentials", nil)
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 5).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.SignedString([]byte(secret))
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "BAD REQUEST", "A erro occured", nil)
	}

	var role models.Role
	if err = config.DB.Where("id = ?", user.RoleID).First(&role).Error; err == nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, "NOT FOUND", "Role Not Found", nil)
	}

	response := dtos.UserResponse{
		ID:     user.ID,
		Avatar: user.Avatar,
		Name:   user.Name,
		Email:  user.Email,
		Dob:    user.Dob,
		Gender: user.Gender,
		Role:   role.Name,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"status":      "OK",
		"message":     "Login successful",
		"user":        response,
		"token":       token,
	})
}

func Me(ctx *fiber.Ctx) error {
	authUser, ok := ctx.Locals("user").(*models.User)
	if !ok || authUser.ID == uuid.Nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, "NOT FOUND", "User Not Found", nil)
	}

	var user models.User
	if err := config.DB.Preload("Role").First(&user, "id = ?", authUser.ID).Error; err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, "NOT FOUND", "User Not Found", err)
	}

	response := dtos.UserResponse{
		ID:     user.ID,
		Avatar: user.Avatar,
		Name:   user.Name,
		Email:  user.Email,
		Dob:    user.Dob,
		Gender: user.Gender,
		Role:   user.Role.Name,
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "OK", "User retreived", "user", response)
}

func UpdateMe(ctx *fiber.Ctx) error {
	req := ctx.Locals("validated").(dtos.UpdateProfileRequest)

	return req

}

func UpdateMyPassword(ctx *fiber.Ctx) error {
	return nil
}
