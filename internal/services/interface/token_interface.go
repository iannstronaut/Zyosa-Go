package _interface

import (
	"time"
	"zyosa/internal/domains/user/entity"
	"zyosa/internal/services"

	"github.com/golang-jwt/jwt/v5"
)

type IJWTService interface {
	GenerateAccessToken(user entity.User, expTime time.Duration) (string, error)

	ValidateAccessToken(token string) (*services.UserClaims, error)

	Unserialize(token string) (*jwt.Token, error)
}