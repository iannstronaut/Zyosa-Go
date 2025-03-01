package middleware

import (
	"fmt"
	"strings"
	"zyosa/internal/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

func EnsureJsonValidRequest[T any](ctx *fiber.Ctx) error {
	body := new(T)
	err := ctx.BodyParser(&body)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.ErrBadRequest.Code, false, err)
	}

	verr := Validator.Struct(body)
	if verr != nil {
		errors := make(map[string]string)

		if validationErrors, ok := verr.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				errors[fieldErr.Field()] = fieldErr.Tag()
			}
		}

		return helpers.ErrorResponse(ctx, fiber.ErrBadRequest.Code, true, ValidationError{Errors: errors}.Errors)
	}

	ctx.Locals("body", body)

	return ctx.Next()
}

type ValidationError struct {
	Errors map[string]string
}

func (v ValidationError) Error() string {
	var messages []string
	for field, err := range v.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, err))
	}
	return strings.Join(messages, ", ")
}
