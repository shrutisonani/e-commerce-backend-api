package auth

import (
	"time"
	"errors"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewEmailVerification (db *sqlx.DB) *Repository {
	return &Repository{Db:db}
}


// Save the email token when user create
func (repository *Repository) SaveEmailToken(userID int, token string, expiry time.Time) error {
	query := `INSERT INTO email_verifications (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := repository.Db.Exec(query, userID, token, expiry)
	return err
}

// Get token with validation
func (repository *Repository) GetValidEmailToken(token string) (*EmailVerifications, error) {
	emailVerified := EmailVerifications{}
	// err := repository.Db.Get(&emailVerified, "SELECT * FROM email_verifications WHERE token = ?", token)
	err := repository.Db.Get(&emailVerified, "SELECT * FROM email_verifications WHERE token = ? AND is_used = false AND expires_at > UTC_TIMESTAMP()", token)
	return &emailVerified, err
}

// Mark valid emali token used
func (repository *Repository) MarkEmailTokenUsed(id int) error {
	_, err := repository.Db.Exec("UPDATE email_verifications SET is_used = true WHERE user_id = ?", id)
	return err
}

// Verify the email token
func (repository *Repository) VerifyEmailToken(token string) error {

	emailToken, err := repository.GetValidEmailToken(token)
	if err != nil {
		log.Error(err)
		return errors.New("invalid or expired token")
	}

	// Mark user verified
	err = repository.VerifyUser(emailToken.UserID)
	if err != nil {
		log.Error(err)
		return err
	}

	// Mark token used
	return repository.MarkEmailTokenUsed(emailToken.UserID)
}

// Verify the email token
func (repository *Repository) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(400, gin.H{"error": "token required"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := repository.VerifyEmailToken(token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(200, gin.H{"message": "email verified successfully"})
}