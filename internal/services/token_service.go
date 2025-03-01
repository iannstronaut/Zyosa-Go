package services

import (
	"fmt"
	"time"
	
	"zyosa/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UUID     string`json:"uuid"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret}
}

// GenerateAccessToken membuat token JWT untuk user tertentu dengan waktu kedaluwarsa
func (j *JWTService) GenerateAccessToken(uuid *types.UUID, expTime time.Duration) (string, error) {
	claims := UserClaims{
		UUID:     uuid.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateAccessToken memvalidasi token JWT dan mengembalikan UserClaims jika valid
func (j *JWTService) ValidateAccessToken(tokenStr string) (*string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return &claims.UUID, nil
}

// Unserialize hanya melakukan parsing token JWT tanpa validasi
func (j *JWTService) Unserialize(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
}
