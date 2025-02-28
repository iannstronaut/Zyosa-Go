package _interface

import "github.com/gofiber/fiber/v2"

type UserInterface interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	GetProfile(ctx *fiber.Ctx) error
	UpdateProfile(ctx *fiber.Ctx) error
	UpdateProfilePicture(ctx *fiber.Ctx) error
	ChangePassword(ctx *fiber.Ctx) error
}