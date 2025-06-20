package helpers

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(data interface{}) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field, _ := reflect.TypeOf(data).FieldByName(err.StructField())
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = strings.ToLower(err.StructField())
		}

		fieldName := strings.ToUpper(string(field.Name[0])) + strings.ToLower(field.Name[1:])

		switch err.Tag() {
		case "required":
			errors[jsonTag] = fieldName + " is required"
		case "email":
			errors[jsonTag] = "Invalid email format"
		case "unique":
			errors[jsonTag] = fieldName + " already exists"
		case "min":
			errors[jsonTag] = fieldName + " should be at least " + err.Param() + " characters"
		case "max":
			errors[jsonTag] = fieldName + " should be at most " + err.Param() + " characters"
		case "file_format":
			errors[jsonTag] = fieldName + " must be one of these formats: jpg, png, jpeg"
		default:
			errors[jsonTag] = "Invalid value"
		}
	}
	return errors
}
