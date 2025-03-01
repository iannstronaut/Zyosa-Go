package routes

import (
	"zyosa/internal/delivery/middleware"
	"zyosa/internal/domains/user"
	"zyosa/internal/helpers"
	"zyosa/internal/services"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	App            *fiber.App
	UserRoute      *user.Handler
	JWTService     *services.JWTService
	AuthMiddleware *middleware.AuthMiddleware
}

func (r *Route) Init() {
	r.App.Get("/", RootHandler)
	api := r.App.Group("/api")
	r.initializeUserRoutes(api.Group("/user"))
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

func (r *Route) initializeUserRoutes(router fiber.Router) {
	router.Post("/login",
		helpers.RateLimiterConfig(2, 5, "Too many requests, please try again later"),
		middleware.EnsureJsonValidRequest[user.LoginRequest],
		r.UserRoute.Login)

	router.Post("/register",
		helpers.RateLimiterConfig(5, 3, "Too many requests, please try again later"),
		middleware.EnsureJsonValidRequest[user.RegisterRequest],
		r.UserRoute.Register)

	router.Post("/logout",
		helpers.RateLimiterConfig(5, 5, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("user"),
		r.UserRoute.Logout)

	router.Get("/profile",
		helpers.RateLimiterConfig(1, 10, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("user"),
		r.UserRoute.GetProfile)

	router.Put("/profile/update",
		helpers.RateLimiterConfig(1, 10, "Too many requests, please try again later"),
		middleware.EnsureJsonValidRequest[user.UpdateProfileRequest],
		r.AuthMiddleware.EnsureAuthenticated("user"),
		r.UserRoute.UpdateProfile)

	router.Put("/profile/update/password",
		helpers.RateLimiterConfig(1, 10, "Too many requests, please try again later"),
		middleware.EnsureJsonValidRequest[user.ChangePasswordRequest],
		r.AuthMiddleware.EnsureAuthenticated("user"),
		r.UserRoute.ChangePassword)
}
