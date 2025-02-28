package helpers

import (
	"zyosa/internal/types"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse[T any](ctx *fiber.Ctx, code int, shouldNotify bool, msg string, content T) error {
	return ctx.Status(code).JSON(types.ApiResponse[T]{
		Success: 	true,
		Code: 		code,
		ShouldNotify: shouldNotify,
		Message: 	msg,
		Content: 	content,
	})
}

func ErrorResponse(ctx *fiber.Ctx, code int, shouldNotify bool, err error) error {
	return ctx.Status(code).JSON(types.ApiResponse[error]{
		Success: 	false,
		Code: 		code,
		ShouldNotify: shouldNotify,
		Message: 	err.Error(),
	})
}