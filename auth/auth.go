package auth

import (
	"net/http"
	"time"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewAuth(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}

// User registration
func (repository *Repository) Register(c *gin.Context) {
	user := Users{}

	// check auth register request
	if err := c.BindJSON(&user); err != nil {
		log.Error(err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// password
	hashed, _ := utils.HashPassword(user.Password, 10)
	user.Password = hashed

	// create user
	_, err := repository.CreateUser(&user)
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}	

	// generate email token
	token, _ := utils.GenerateToken()
	expiry := time.Now().Add(15 * time.Minute)

	// save token in email_verifications
	emailToken := repository.SaveEmailToken(user.Id, token, expiry)
	
	if emailToken != nil {
		log.Error(emailToken)
		c.JSON(500, gin.H{"error": emailToken.Error()})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// update the user verify
	isUserVerify := repository.VerifyUser(user.Id)

	if isUserVerify != nil {
		log.Error(isUserVerify)
		c.JSON(500, gin.H{"error": isUserVerify.Error()})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(201, gin.H{
		"message": "user registered and verify email",
		"user" : user,
	})
}

// Create the user (insert the data query)
func (repository *Repository) CreateUser(user *Users) (int64, error) {
	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	data, err := repository.Db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	
	id, err := data.LastInsertId()
	if err != nil {
		return 0, err
	}
	user.Id = int(id) // update struct
	return id, nil
}

// Get user by email id
func (repository *Repository) GetUserByEmail(email string) (*Users, error) {
	user := Users{}
	err := repository.Db.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	return &user, err
}

// Save the email token when user create
func (repository *Repository) SaveEmailToken(userID int, tokenHash string, expiry time.Time) error {
	query := `INSERT INTO email_verifications (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := repository.Db.Exec(query, userID, tokenHash, expiry)
	return err
}

// Get email token
func (repository *Repository) GetEmailToken(tokenHash string) (*EmailVerifications, error) {
	emailVerified := EmailVerifications{}
	err := repository.Db.Get(&emailVerified, "SELECT * FROM email_verifications WHERE token = ?", tokenHash)
	return &emailVerified, err
}

// When email verification is success then make user is_verified
func (repository *Repository) VerifyUser(userID int) error {
	_, err := repository.Db.Exec("UPDATE users SET is_verified = true WHERE id = ?", userID)
	return err
}
