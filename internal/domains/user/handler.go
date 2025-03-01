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
	userRepo   *repository.UserRepository
	viper      *viper.Viper
	JWTService *services.JWTService
}

func NewHandler(userRepo *repository.UserRepository, viper *viper.Viper, jwtServices *services.JWTService) *Handler {
	return &Handler{
		userRepo:   userRepo,
		viper:      viper,
		JWTService: jwtServices,
	}
}

// @Summary Register a new user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{} "Username or email has been taken"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/register [post]
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

// @Summary Login user
// @Description Authenticates user and returns an access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "User login data"
// @Success 200 {string} string "Login success"
// @Failure 401 {object} map[string]interface{} "Username or password is incorrect"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/login [post]
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

	access_token, err := handler.JWTService.GenerateAccessToken(exist.UUID, "user", time.Hour)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access_token,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return helpers.SuccessResponse(ctx, 200, handler.viper.GetBool("app.debug"), "Login success", map[string]interface{}{
		"username": exist.Username,
		"email":    exist.Email,
		"first_name": exist.FirstName,
		"last_name":  exist.LastName,
	})
}

// @Summary Logout user
// @Description Logs out the user by clearing the access token
// @Tags Auth
// @Success 200 {string} string "Logout success"
// @Router /user/logout [post]
func (handler *Handler) Logout(ctx *fiber.Ctx) error {
	ctx.ClearCookie("access_token")
	return helpers.SuccessResponse[string](ctx, fiber.StatusOK, handler.viper.GetBool("app.debug"), "Logout success")
}

// @Summary Get user profile
// @Description Retrieves the profile of the authenticated user
// @Tags User
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /user/profile [get]
func (handler *Handler) GetProfile(ctx *fiber.Ctx) error {
	req := ctx.Locals("user").(*services.UserClaims)

	userId := entity.User{UUID: req.UUID}

	exist, err := handler.userRepo.FindByID(userId.ID.String(), &userId)
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

// @Summary Update user profile
// @Description Updates the authenticated user's profile information
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body UpdateProfileRequest true "User profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 409 {object} map[string]interface{} "Username or email is already in use"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile [put]
func (handler *Handler) UpdateProfile(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*UpdateProfileRequest)
	req := ctx.Locals("user").(*services.UserClaims)

	user := entity.User{}
	_, err := handler.userRepo.FindByID(req.UUID.ID.String(), &user)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, true, fmt.Errorf("user not found"))
	}

	if body.Username != nil && *body.Username != "" {
		existingUser := entity.User{Username: *body.Username}
		_, err := handler.userRepo.FindByUsername(&existingUser)
		if err == nil && existingUser.ID != user.ID {
			return helpers.ErrorResponse(ctx, fiber.StatusConflict, true, fmt.Errorf("username is already taken"))
		}
		user.Username = *body.Username
	}

	if body.Email != nil && *body.Email != "" {
		existingUser := entity.User{Email: *body.Email}
		_, err := handler.userRepo.FindByEmail(&existingUser)
		if err == nil && existingUser.ID != user.ID {
			return helpers.ErrorResponse(ctx, fiber.StatusConflict, true, fmt.Errorf("email is already in use"))
		}
		user.Email = *body.Email
	}

	// Update field lain jika ada
	if body.FirstName != nil && *body.FirstName != "" {
		user.FirstName = *body.FirstName
	}
	if body.LastName != nil && *body.LastName != "" {
		user.LastName = *body.LastName
	}
	if body.Phone != nil && *body.Phone != "" {
		user.Phone = body.Phone
	}

	// Simpan perubahan ke database
	err = handler.userRepo.Update(&user)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("failed to update profile"))
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, handler.viper.GetBool("app.debug"), "Profile updated successfully", map[string]interface{}{
		"id":         user.UUID.ID,
		"username":   user.Username,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"phone":      user.Phone,
	})
}

func (handler *Handler) UpdateProfilePicture(ctx *fiber.Ctx) error {
	// Implementation for updating user profile picture
	return nil
}

// @Summary Change user password
// @Description Changes the authenticated user's password
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "New password and confirm password must be the same"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile/password [put]
func (handler *Handler) ChangePassword(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*ChangePasswordRequest)
	req := ctx.Locals("user").(*services.UserClaims)

	if body.NewPassword != body.ConfirmPassword {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, true, fmt.Errorf("new password and confirm password must be the same"))
	}

	user := entity.User{}
	_, err := handler.userRepo.FindByID(req.UUID.ID.String(), &user)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, true, fmt.Errorf("user not found"))
	}

	hashedPassword, err := utils.HashPassword(body.NewPassword, handler.viper.GetString("app.secret"))
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("internal server error"))
	}

	user.Password = hashedPassword

	err = handler.userRepo.Update(&user)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, true, fmt.Errorf("failed to change password"))
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, handler.viper.GetBool("app.debug"), "Password changed successfully", map[string]interface{}{
		"id": user.UUID.ID,
	})
}
