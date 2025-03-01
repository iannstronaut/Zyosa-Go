package routes

import (
	"zyosa/internal/delivery/middleware"
	Admin "zyosa/internal/domains/admin"
	User "zyosa/internal/domains/user"
	"zyosa/internal/helpers"
	"zyosa/internal/services"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	App            *fiber.App
	UserRoute      *User.Handler
	AdminRoute     *Admin.Handler
	JWTService     *services.JWTService
	AuthMiddleware *middleware.AuthMiddleware
}

func (r *Route) Init() {
	r.App.Get("/", RootHandler)
	api := r.App.Group("/api")
	r.initializeUserRoutes(api.Group("/user"))
	r.initializeAdminRoutes(api.Group("/admin"))
}

func RootHandler(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("Hello, World!")
}

func (r *Route) initializeUserRoutes(router fiber.Router) {
	router.Post("/login",
		helpers.RateLimiterConfig(2, 5, "Too many requests, please try again later"),
		middleware.EnsureJsonValidRequest[User.LoginRequest],
		r.UserRoute.Login)

	router.Post("/register",
		helpers.RateLimiterConfig(5, 3, "Too many requests, please try again later"),
		middleware.EnsureJsonValidRequest[User.RegisterRequest],
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
		r.AuthMiddleware.EnsureAuthenticated("user"),
		middleware.EnsureJsonValidRequest[User.UpdateProfileRequest],
		r.UserRoute.UpdateProfile)

	router.Put("/profile/update/password",
		helpers.RateLimiterConfig(1, 10, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("user"),
		middleware.EnsureJsonValidRequest[User.ChangePasswordRequest],
		r.UserRoute.ChangePassword)
}

func (r *Route) initializeAdminRoutes(router fiber.Router){
	router.Post("/login",
		helpers.RateLimiterConfig(2, 5, "Too many requests, please try again later"),
		middleware.EnsureJsonValidRequest[Admin.LoginRequest],
		r.AdminRoute.Login)

	router.Post("/add",
		helpers.RateLimiterConfig(5, 3, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("admin"),
		middleware.EnsureJsonValidRequest[Admin.AddAdminRequest],
		r.AdminRoute.AddAdmin)
	
	router.Post("/logout",
		helpers.RateLimiterConfig(5, 5, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("admin"),
		r.AdminRoute.Logout)

	router.Get("/profile",
		helpers.RateLimiterConfig(1, 10, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("admin"),
		r.AdminRoute.GetProfile)

	router.Put("/profile/update",
		helpers.RateLimiterConfig(1, 10, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("admin"),
		middleware.EnsureJsonValidRequest[Admin.UpdateProfileRequest],
		r.AdminRoute.UpdateProfile)

	router.Put("/profile/update/password",
		helpers.RateLimiterConfig(1, 10, "Too many requests, please try again later"),
		r.AuthMiddleware.EnsureAuthenticated("admin"),
		middleware.EnsureJsonValidRequest[Admin.ChangePasswordRequest],
		r.AdminRoute.ChangePassword)
}
