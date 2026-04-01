package auth

import (
	"errors"
	"time"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewAuthToken(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}

func (repository *Repository) SaveRefreshToken(userID int, token string, expiry time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := repository.Db.Exec(query, userID, token, expiry)
	return err
}

func (repository *Repository) GetRefreshToken(token string) (*RefreshTokens, error) {
	rt := RefreshTokens{}

	query := `
        SELECT * FROM refresh_tokens 
        WHERE token = ? AND expires_at > UTC_TIMESTAMP()
    `
	err := repository.Db.Get(&rt, query, token)
	return &rt, err
}

func (repository *Repository) DeleteRefreshToken(id int) error {
	_, err := repository.Db.Exec("DELETE FROM refresh_tokens WHERE id = ?", id)
	return err
}

func (repository *Repository) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	tokenHash := utils.HashToken(req.RefreshToken)

	access, refresh, err := repository.CheckRefreshToken(tokenHash)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (repository *Repository) CheckRefreshToken(oldToken string) (string, string, error) {

	// Check in DB
	rt, err := repository.GetRefreshToken(oldToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")

	// Check expiry
	if time.Now().In(loc).After(rt.ExpiresAt) {
		return "", "", errors.New("refresh token expired")
	}

	// Generate new tokens
	newAccess, _ := utils.GenerateAccessToken(rt.UserID, "user")
	newRefresh, _ := utils.GenerateRefreshToken(rt.UserID)

	// Rotate token (delete old)
	repository.DeleteRefreshToken(rt.Id)

	// Store new refresh token
	newHash := utils.HashToken(newRefresh)

	expiry := time.Now().In(loc).Add(7 * 24 * time.Hour)

	// Save the refresh token
	err = repository.SaveRefreshToken(rt.UserID, newHash, expiry)
	if err != nil {
		log.Error(err)
		return "", "", errors.New("failed to save refresh token")
	}

	return newAccess, newRefresh, nil
}
