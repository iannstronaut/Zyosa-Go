package helpers

import (
	"zyosa/internal/types"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse[T any](ctx *fiber.Ctx, code int, shouldNotify bool, msg string, content ...T) error {
	var responseContent T

	if len(content) > 0 {
		responseContent = content[0]
	}

	return ctx.Status(code).JSON(types.ApiResponse[T]{
		Success:      true,
		Code:         code,
		ShouldNotify: shouldNotify,
		Message:      msg,
		Content:      responseContent,
	})
}

func ErrorResponse(ctx *fiber.Ctx, code int, shouldNotify bool, err any) error {
	var errMsg any

	switch e := err.(type) {
	case error:
		errMsg = e.Error()
	case string:
		errMsg = e
	case map[string]string:
		errMsg = e
	default:
		errMsg = "Unknown error"
	}

	return ctx.Status(code).JSON(types.ApiResponse[error]{
		Success:      false,
		Code:         code,
		ShouldNotify: shouldNotify,
		Message:      errMsg,
	})
}
