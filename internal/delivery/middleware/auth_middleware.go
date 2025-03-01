package middleware

import (
	"fmt"
	"zyosa/internal/helpers"
	"zyosa/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	JWTService *services.JWTService
}

func NewAuthMiddleware(jwtService *services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		JWTService: jwtService,
	}
}

func (m *AuthMiddleware) EnsureAuthenticated(requiredRole string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Cookies("access_token")
		if token == "" {
			return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("unauthorized"))
		}

		claims, err := m.JWTService.ValidateAccessToken(token)
		if err != nil {
			return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("unauthorized"))
		}

		userRole := claims.Roles
		if userRole != requiredRole {
			return helpers.ErrorResponse(ctx, fiber.StatusForbidden, true, fmt.Errorf("forbidden: insufficient role"))
		}

		ctx.Locals("user", claims)

		return ctx.Next()
	}
}
