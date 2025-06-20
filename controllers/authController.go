package controllers

import (
	"fmt"
	"go-api/config"
	"go-api/helpers"
	"go-api/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	data := new(models.User)
	if err := ctx.BodyParser(data); err != nil {
		return helpers.Error(
			ctx,
			fiber.StatusBadRequest,
			"Invalid request body",
			err,
		)
	}

	if errors := helpers.Validate(data); errors != nil {
		return helpers.Error(
			ctx,
			fiber.StatusBadRequest,
			"Validation failed",
			errors,
		)
	}

	file, _ := ctx.FormFile("avatar")
	fileExt := strings.ToLower(strings.Split(file.Filename, ".")[1])
	allowedExtensions := map[string]bool{"jpg": true, "jpeg": true, "png": true, "gif": true}
	if !allowedExtensions[fileExt] {
		return helpers.Error(
			ctx,
			fiber.StatusBadRequest,
			"Invalid file type. Only jpg, jpeg, png, and gif are allowed",
			nil,
		)
	}
	fileName := strings.Replace(uuid.New().String(), "-", "", -1)
	userAvatar := fmt.Sprintf("%s.%s", fileName, fileExt)

	if err := ctx.SaveFile(file, fmt.Sprintf("./public/avatars/%s", userAvatar)); err != nil {
		return helpers.Error(
			ctx,
			fiber.StatusBadRequest,
			"Failed to save avatar",
			nil,
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return helpers.Error(
			ctx,
			fiber.StatusBadRequest,
			"Could not hash password",
			nil,
		)
	}

	user := models.User{
		Avatar:   userAvatar,
		Name:     data.Name,
		Email:    data.Email,
		Dob:      data.Dob,
		Gender:   data.Gender,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user); err != nil {
		return helpers.Error(
			ctx,
			fiber.StatusInternalServerError,
			"An error occured",
			nil,
		)
	}

	return helpers.Success(
		ctx,
		fiber.StatusCreated,
		"User created successfully",
		"user",
		user,
	)
}

func Login(ctx *fiber.Ctx) error {
	return ctx.SendString("Login endpoint hit")
}
