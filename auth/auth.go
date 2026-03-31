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

// Register handles user registration
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
	loc, _ := time.LoadLocation("Asia/Kolkata")
	expiry := time.Now().In(loc).Add(15 * time.Minute)

	// save token in email_verifications
	emailToken := repository.SaveEmailToken(user.Id, token, expiry)

	if emailToken != nil {
		log.Error(emailToken)
		c.JSON(500, gin.H{"error": emailToken.Error()})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// verify the email token
	err = repository.VerifyEmailToken(token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(201, gin.H{
		"message": "user registered and email-verified successfully",
		"user":    user,
	})
}

// Login handles user login
func (repository *Repository) Login(c *gin.Context) {
	req := Users{}

	// check auth login request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check the user by email id
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		log.Error(err)
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// compare the password is correct or not
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// check email is verifyed or not
	if !user.IsVerified {
		c.JSON(401, gin.H{"error": "email not verified"})
		return
	}

	// generate tokens
	access, _ := utils.GenerateAccessToken(user.Id)
	refresh, _ := utils.GenerateRefreshToken(user.Id)

	// Store new refresh token
	newHash := utils.HashToken(refresh)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	expiry := time.Now().In(loc).Add(7 * 24 * time.Hour)

	// Save the refresh token
	err = repository.SaveRefreshToken(user.Id, newHash, expiry)
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "failed to save refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user":          user,
	})

}
