package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpirationTime = 24 * time.Hour
)

type AuthClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func SignToken(userID uint, secret string) (string, error) {
	claims := AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "task-master-go",
			Subject:   "user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(secret)

	return ss, err
}
