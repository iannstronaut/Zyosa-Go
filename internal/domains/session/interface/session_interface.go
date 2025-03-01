package _interface

import "github.com/gofiber/fiber/v2"

type SessionInterface interface {
	RefreshToken(ctx *fiber.Ctx) error
	Callback(role string) fiber.Handler 
}