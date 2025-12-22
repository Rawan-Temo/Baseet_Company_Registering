package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(req interface{} ) []string {
    if err := validate.Struct(req); err != nil {
        validationErrors := err.(validator.ValidationErrors)
        // Format validation errors
        var errors []string
        for _, err := range validationErrors {
            errors = append(errors, fmt.Sprintf("%s: %s", err.Field(), err.Tag()))
        }
        return errors
    }

	return nil
}