package user

import (
	"fmt"
	"time"
	"zyosa/internal/domains/user/entity"
	"zyosa/internal/domains/user/repository"
	"zyosa/internal/helpers"
	"zyosa/internal/services"
	"zyosa/internal/types"
	"zyosa/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Handler struct {
	userRepo *repository.UserRepository
	viper    *viper.Viper
	JWTService *services.JWTService
}

func NewHandler(userRepo *repository.UserRepository, viper *viper.Viper, jwtServices *services.JWTService) *Handler {
	return &Handler{
		userRepo: userRepo,
		viper:    viper,
		JWTService: jwtServices,
	}
}

func (handler *Handler) Register(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*RegisterRequest)

	user := entity.User{Username: body.Username}
	email := entity.User{Email: body.Email}
	
	userErr, _ := handler.userRepo.FindByUsername(&user)
	emailErr, _ := handler.userRepo.FindByEmail(&email)
	if userErr != nil || emailErr != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusConflict, true, fmt.Errorf("username or email has been taken"))
	}

	hashedPassword, err := utils.HashPassword(body.Password, handler.viper.GetString("app.secret"))
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	newUser := entity.User{
		UUID:      types.UUID{ID: uuid},
		Username:  body.Username,
		Email:     body.Email,
		Password:  hashedPassword,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}

	if err := handler.userRepo.Create(&newUser); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusConflict, true, fmt.Errorf("register failed"))
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, handler.viper.GetBool("app.debug"), "Register success", map[string]interface{}{
		"id":         newUser.UUID.ID,
		"username":   newUser.Username,
		"email":      newUser.Email,
		"first_name": newUser.FirstName,
		"last_name":  newUser.LastName,
	})
}

func (handler *Handler) Login(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*LoginRequest)

	user := entity.User{Username: body.Username}
	exist, err := handler.userRepo.FindByUsername(&user) 
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("username or password is incorrect"))
	}

	if !utils.ComparePassword(exist.Password, body.Password, handler.viper.GetString("app.secret")) {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("username or password is incorrect"))
	}

	access_token, err := handler.JWTService.GenerateAccessToken(&exist.UUID, time.Hour)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	ctx.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: access_token,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return helpers.SuccessResponse[string](ctx, 200, handler.viper.GetBool("app.debug"), "Login success")
}

func (handler *Handler) Logout(ctx *fiber.Ctx) error {
	// Implementation for user logout
	return nil
}

func (handler *Handler) GetProfile(ctx *fiber.Ctx) error {
	// Implementation for getting user profile
	return nil
}

func (handler *Handler) UpdateProfile(ctx *fiber.Ctx) error {
	// Implementation for updating user profile
	return nil
}

func (handler *Handler) UpdateProfilePicture(ctx *fiber.Ctx) error {
	// Implementation for updating user profile picture
	return nil
}

func (handler *Handler) ChangePassword(ctx *fiber.Ctx) error {
	// Implementation for changing user password
	return nil
}
