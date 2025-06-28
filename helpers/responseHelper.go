package helpers

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ApiResponse struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx *fiber.Ctx, statusCode int, status string, message string, dataKey string, data interface{}) error {
	response := fiber.Map{
		"status_code": statusCode,
		"status":      status,
		"message":     message,
	}

	if data != nil && dataKey != "" {
		response[dataKey] = data
	}

	return ctx.Status(statusCode).JSON(response)
}

func ErrorResponse(ctx *fiber.Ctx, statusCode int, status string, message string, error interface{}) error {
	if error != nil {
		log.Printf("Error: %v", error)
	}

	return ctx.Status(statusCode).JSON(fiber.Map{
		"status_code": statusCode,
		"status":      status,
		"message":     message,
	})
}

func ValidationErrorResponse[T any](ctx *fiber.Ctx, err error) error {
	errors := make(map[string]string)
	var requestType T // to extract JSON tags

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := getJSONFieldName(requestType, e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", capitalize(field))
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", capitalize(field), e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s must not exceed %s characters", capitalize(field), e.Param())
			case "oneof":
				errors[field] = fmt.Sprintf("%s must be one of: %s", capitalize(field), e.Param())
			case "datetime":
				errors[field] = fmt.Sprintf("%s must be in the format %s", capitalize(field), e.Param())
			default:
				errors[field] = fmt.Sprintf("%s is not valid", capitalize(field))
			}
		}
	}

	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"status_code": fiber.StatusUnprocessableEntity,
		"status":      "VALIDATION ERROR",
		"errors":      errors,
	})
}

func getJSONFieldName(structType interface{}, fieldName string) string {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(fieldName); ok {
		jsonTag := f.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			return splitJSONTag(jsonTag)
		}
	}
	return fieldName
}

func splitJSONTag(tag string) string {
	if commaIdx := strings.Index(tag, ","); commaIdx != -1 {
		return tag[:commaIdx]
	}
	return tag
}

func capitalize(field string) string {
	if len(field) == 0 {
		return field
	}
	return strings.ToUpper(field[:1]) + field[1:]
}
