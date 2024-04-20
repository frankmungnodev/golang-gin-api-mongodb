package utils

import (
	"github.com/go-playground/validator/v10"
)

func GetValidationMessage(ve validator.ValidationErrors) string {
	for _, fe := range ve {
		return msgForTag(fe)
	}

	return ve.Error()
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.StructField() + " field is required"

	case "email":
		return "Invalid email"

	case "min":
		return fe.StructField() + " must be at least " + fe.Param() + " characters long"

	case "max":
		return fe.StructField() + " must be less than " + fe.Param() + " characters"
	}
	return fe.Error()
}
