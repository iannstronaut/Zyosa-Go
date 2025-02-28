package user

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=50"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Email     string  `json:"email" validate:"required,email"`
	Phone     *string `json:"phone" validate:"omitempty,e164"`
	Password  string  `json:"password" validate:"required,min=8"`
}

type UpdateProfileRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name" validate:"omitempty"`
	LastName  *string `json:"last_name" validate:"omitempty"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Phone     *string `json:"phone" validate:"omitempty,e164"`
}

type UpdateProfilePictureRequest struct {
	Image string `json:"image" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}