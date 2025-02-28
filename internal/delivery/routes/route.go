package routes

import (
	"github.com/gofiber/fiber/v2"
)

type Route struct {
	App *fiber.App
}

func (r *Route) Init() {
	r.App.Get("/", RootHandler)
}

// @Summary Root Endpoint
// @Description Returns a simple hello message
// @Tags Root
// @Accept json
// @Produce json
// @Success 200 {string} string "Hello, World!"
// @Router / [get]
func RootHandler(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("Hello, World!")
}