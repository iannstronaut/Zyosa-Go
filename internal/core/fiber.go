package core

import (
	"zyosa/internal/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

func NewFiber(viper *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName: viper.GetString("app.name"),
		ErrorHandler: NewErrorHandler(viper.GetBool("app.debug")),
	})

	log := logger.New(logger.Config{
		Format: "${time} [${method}] ${path} \t-\t ${status} ${latency}\n",
		TimeFormat: "[01-Jan-2006] [15:04:05]",
		TimeZone: "Asia/Jakarta",
	})

	app.Use(log)

	return app
}

func NewErrorHandler(shouldNotify bool) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if err, ok := err.(*fiber.Error); ok {
			code = err.Code
		}

		return helpers.ErrorResponse(ctx, code, shouldNotify, err)
	}
}