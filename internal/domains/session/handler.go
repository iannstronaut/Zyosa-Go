package session

import (
	"fmt"
	"time"
	"zyosa/internal/domains/session/entity"
	"zyosa/internal/domains/session/repository"
	"zyosa/internal/helpers"
	"zyosa/internal/services"
	"zyosa/internal/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Handler struct {
	sessionRepo *repository.SessionRepository
	viper	   	*viper.Viper
	JWTService  *services.JWTService
}

func NewHandler(sessionRepo *repository.SessionRepository, viper *viper.Viper, jwtService *services.JWTService) *Handler {
	return &Handler{
		sessionRepo: sessionRepo,
		viper:	   viper,
		JWTService: jwtService,
	}
}

func (handler *Handler) RefreshToken(ctx *fiber.Ctx) error {
	access_token := ctx.Locals("access_token").(string)
	refresh_token := ctx.Cookies("refresh_token")

	claims, err := handler.JWTService.Unserialize(access_token)
	if err != nil {
		fmt.Printf("Error Validate: %v", err)
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	exist, err := handler.sessionRepo.FindByToken(refresh_token)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("unauthorized"))
	}

	if exist.UserId != claims.UUID.ID.String() {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("unauthorized"))
	}

	expiredAt, err := time.Parse(time.RFC3339, exist.ExpiredAt)
	if err != nil {
		fmt.Println("Parsing error:", err)
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("invalid expiration date: %v", err))
	}	

	today := time.Now().UTC()

	if expiredAt.Before(today) {
    	return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("unauthorized"))
	}

	access_token, err = handler.JWTService.GenerateAccessToken(claims.UUID, claims.Roles, time.Hour * 24)
	if err != nil {
		fmt.Printf("Error Create: %v", err)
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	exist.UsedAt = time.Now().Format("2006-01-02 15:04:05")

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access_token,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return helpers.SuccessResponse(ctx, fiber.StatusOK, true, "refresh token success", map[string]interface{}{
		"id": claims.UUID.ID.String(),
	})
}

func (handler *Handler) GenerateRefreshToken(ctx *fiber.Ctx) error{
	userId := ctx.Locals("userId").(types.UUID)
	userAgent := ctx.Get("User-Agent")
	token, err := uuid.NewV7()
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}
	expiredAt := time.Now().AddDate(0, 1, 0).Format("2006-01-02 15:04:05")

	newSession := entity.Session{
		UserId: userId.ID.String(),
		Token: token.String(),
		UserAgent: userAgent,
		ExpiredAt: expiredAt,
		UsedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := handler.sessionRepo.Create(&newSession); err != nil {
		fmt.Printf("Error: %v", err)
		return helpers.ErrorResponse(ctx, fiber.StatusConflict, true, fmt.Errorf("create session failed"))
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token.String(),
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return helpers.SuccessResponse(ctx, fiber.StatusOK, true, "Login Success", map[string]interface{}{
		"id": userId.ID.String(),
	})
} 

func (handler *Handler) Callback(ctx *fiber.Ctx) error {
	return helpers.SuccessResponse[string](ctx, fiber.StatusOK, true, "success")
}