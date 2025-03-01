package services

import (
	"fmt"
	"time"

	"zyosa/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	types.UUID
	Roles string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret}
}

// GenerateAccessToken membuat token JWT untuk user tertentu dengan waktu kedaluwarsa
func (j *JWTService) GenerateAccessToken(uuid types.UUID, role string, expTime time.Duration) (string, error) {
	claims := UserClaims{
		UUID:  uuid,
		Roles: role,
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
func (j *JWTService) ValidateAccessToken(tokenStr string) (*UserClaims, error) {
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

	return claims, nil
}

// Unserialize hanya melakukan parsing token JWT tanpa validasi
func (j *JWTService) Unserialize(tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}

	// Gunakan ParseWithClaims tetapi abaikan validasi expiry
	_, _, err := new(jwt.Parser).ParseUnverified(tokenStr, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	return claims, nil
}

