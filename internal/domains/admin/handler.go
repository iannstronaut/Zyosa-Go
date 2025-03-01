package admin

import (
	"fmt"
	"time"
	"zyosa/internal/domains/admin/entity"
	"zyosa/internal/domains/admin/repository"
	"zyosa/internal/helpers"
	"zyosa/internal/services"
	"zyosa/internal/types"
	"zyosa/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)


type Handler struct {
	adminRepo   *repository.AdminRepository
	viper      *viper.Viper
	JWTService *services.JWTService
}

func NewHandler(adminRepo *repository.AdminRepository, viper *viper.Viper, jwtServices *services.JWTService) *Handler {
	return &Handler{
		adminRepo:   adminRepo,
		viper:      viper,
		JWTService: jwtServices,
	}
}

func (handler *Handler) AddAdmin(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*AddAdminRequest)

	user := entity.Admin{Username: body.Username}
	email := entity.Admin{Email: body.Email}

	userErr, _ := handler.adminRepo.FindByUsername(&user)
	emailErr, _ := handler.adminRepo.FindByEmail(&email)
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

	newUser := entity.Admin{
		UUID:      types.UUID{ID: uuid},
		Username:  body.Username,
		Email:     body.Email,
		Password:  hashedPassword,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}

	if err := handler.adminRepo.Create(&newUser); err != nil {
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

	admin := entity.Admin{Username: body.Username}
	exist, err := handler.adminRepo.FindByUsername(&admin)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("username or password is incorrect"))
	}

	if !utils.ComparePassword(exist.Password, body.Password, handler.viper.GetString("app.secret")) {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, true, fmt.Errorf("username or password is incorrect"))
	}

	access_token, err := handler.JWTService.GenerateAccessToken(exist.UUID, "admin", time.Hour)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access_token,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	ctx.Locals("userId", exist.UUID)

	return ctx.Next()
}
func (handler *Handler) Logout(ctx *fiber.Ctx) error {
	ctx.ClearCookie("access_token")
	return helpers.SuccessResponse[string](ctx, fiber.StatusOK, handler.viper.GetBool("app.debug"), "Logout success")
}
func (handler *Handler) GetProfile(ctx *fiber.Ctx) error {
	req := ctx.Locals("user").(*services.UserClaims)

	userId := entity.Admin{UUID: req.UUID}

	exist, err := handler.adminRepo.FindByID(userId.ID.String(), &userId)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, true, fmt.Errorf("user not found"))
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, handler.viper.GetBool("app.debug"), "Get profile success", map[string]interface{}{
		"id":         exist.UUID.ID,
		"username":   exist.Username,
		"email":      exist.Email,
		"first_name": exist.FirstName,
		"last_name":  exist.LastName,
		"phone":      exist.Phone,
	})
}
func (handler *Handler) UpdateProfile(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*UpdateProfileRequest)
	req := ctx.Locals("user").(*services.UserClaims)

	admin := entity.Admin{}
	_, err := handler.adminRepo.FindByID(req.UUID.ID.String(), &admin)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, true, fmt.Errorf("user not found"))
	}

	if body.Username != nil && *body.Username != "" {
		existingUser := entity.Admin{Username: *body.Username}
		_, err := handler.adminRepo.FindByUsername(&existingUser)
		if err == nil && existingUser.ID != admin.ID {
			return helpers.ErrorResponse(ctx, fiber.StatusConflict, true, fmt.Errorf("username is already taken"))
		}
		admin.Username = *body.Username
	}

	if body.Email != nil && *body.Email != "" {
		existingUser := entity.Admin{Email: *body.Email}
		_, err := handler.adminRepo.FindByEmail(&existingUser)
		if err == nil && existingUser.ID != admin.ID {
			return helpers.ErrorResponse(ctx, fiber.StatusConflict, true, fmt.Errorf("email is already in use"))
		}
		admin.Email = *body.Email
	}

	// Update field lain jika ada
	if body.FirstName != nil && *body.FirstName != "" {
		admin.FirstName = *body.FirstName
	}
	if body.LastName != nil && *body.LastName != "" {
		admin.LastName = *body.LastName
	}
	if body.Phone != nil && *body.Phone != "" {
		admin.Phone = body.Phone
	}

	// Simpan perubahan ke database
	err = handler.adminRepo.Update(&admin)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("failed to update profile"))
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, handler.viper.GetBool("app.debug"), "Profile updated successfully", map[string]interface{}{
		"id":         admin.UUID.ID,
		"username":   admin.Username,
		"email":      admin.Email,
		"first_name": admin.FirstName,
		"last_name":  admin.LastName,
		"phone":      admin.Phone,
	})
}
func (handler *Handler) UpdateProfilePicture(ctx *fiber.Ctx) error {
	return nil
}
func (handler *Handler) ChangePassword(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*ChangePasswordRequest)
	req := ctx.Locals("user").(*services.UserClaims)

	if body.NewPassword != body.ConfirmPassword {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, true, fmt.Errorf("new password and confirm password must be the same"))
	}

	admin := entity.Admin{}
	_, err := handler.adminRepo.FindByID(req.UUID.ID.String(), &admin)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, true, fmt.Errorf("user not found"))
	}

	hashedPassword, err := utils.HashPassword(body.NewPassword, handler.viper.GetString("app.secret"))
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	admin.Password = hashedPassword

	err = handler.adminRepo.Update(&admin)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("failed to change password"))
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, handler.viper.GetBool("app.debug"), "Password changed successfully", map[string]interface{}{
		"id": admin.UUID.ID,
	})
}