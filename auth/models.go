package auth

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Db *sqlx.DB
}

// USERS TABLE
type Users struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

// REFRESH TOKENS TABLE
type RefreshTokens struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// PASSWORD RESETS TABLE
type PasswordResets struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// EMAIL VERIFICATIONS TABLE
type EmailVerifications struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}

// USER REGISTRATION
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LOGIN REQUEST
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// REFRESH TOKEN REQUEST
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
