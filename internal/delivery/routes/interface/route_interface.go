package _interface

import "github.com/gofiber/fiber/v2"

type IRoute interface {
	Init()
	initializeUserRoutes(router fiber.Router)
	initializeAdminRoutes(router fiber.Router)
}
