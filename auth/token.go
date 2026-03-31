package auth

import (
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

	// c.ShouldBindJSON(&req)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	tokenHash := utils.HashToken(req.RefreshToken)

	refresh, err := repository.GetRefreshToken(tokenHash)
	if err != nil {
		log.Error(err)
		c.JSON(401, gin.H{"error": "invalid refresh token"})
		return
	}

	// Rotate token (delete old)
	repository.DeleteRefreshToken(refresh.Id)

	// Generate new tokens
	newAccess, _ := utils.GenerateAccessToken(refresh.UserID)
	newRefresh, _ := utils.GenerateRefreshToken(refresh.UserID)

	// Store new refresh token
	newHash := utils.HashToken(newRefresh)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	expiry := time.Now().In(loc).Add(7 * 24 * time.Hour)

	// repository.SaveRefreshToken(refresh.UserID, newHash, expiry)

	// Save the refresh token
	err = repository.SaveRefreshToken(refresh.UserID, newHash, expiry)
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "failed to save refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}
