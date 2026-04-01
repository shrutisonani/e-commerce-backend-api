package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var accessSecret = []byte("ACCESS_SECRET")
var refreshSecret = []byte("REFRESH_SECRET")

// GenerateAccessToken generates a JWT access token valid for 15 minutes.
func GenerateAccessToken(userID int, role string) (string, error) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().In(loc).Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// GenerateRefreshToken generates a JWT refresh token valid for 7 days.
func GenerateRefreshToken(userID int) (string, error) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().In(loc).Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}
