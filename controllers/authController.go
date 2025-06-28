package controllers

import (
	"go-api/config"
	"go-api/helpers"
	"go-api/models"
	"go-api/requests"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) (err error) { // <== named return value
	req := ctx.Locals("validated").(requests.CreateUserRequest)

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

	newUser := models.User{
		Avatar:   userAvatar,
		Name:     req.Name,
		Email:    req.Email,
		Dob:      req.Dob,
		Gender:   req.Gender,
		Password: string(hashedPassword),
	}

	if err = config.DB.Create(&newUser).Error; err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, "INTERNAL SERVER ERROR", "Could not create user", nil)
	}

	// prevent defer from deleting the file
	err = nil
	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "CREATED", "User created successfully", "user", newUser)
}

func Login(ctx *fiber.Ctx) error {
	req := ctx.Locals("validated").(requests.LoginUserRequest)

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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"status":      "SUCCESS",
		"message":     "Login successful",
		"user":        user,
		"token":       token,
	})
}
