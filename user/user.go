package user

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

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
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// PASSWORD RESETS TABLE
type PasswordResets struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// EMAIL VERIFICATIONS TABLE
type EmailVerifications struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserRepository struct {
	Db *sqlx.DB
}

func NewUser(db *sqlx.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (repository *UserRepository) UserRegister(c *gin.Context) {
	user := Users{}

	// binf input data
	if err := c.BindJSON(&user); err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// password
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(user.Password), 10)

	user.Password = string(hash)

	// insert data
	query := `INSERT INTO users (name,email,password) VALUES (?,?,?)`

	_, err := repository.Db.Exec(query, user.Name, user.Email, user.Password)

	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (repository *UserRepository) Users(c *gin.Context) {

	userList, err := repository.GetUsers()

	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userList)

}

func (repository *UserRepository) GetUsers() ([]Users, error) {

	users := []Users{}

	err := repository.Db.Select(&users, "SELECT * FROM users")

	return users, err
}

func (repository *UserRepository) UserById(c *gin.Context) {

	user := Users{}

	id := c.Param("id")

	err := repository.Db.Get(&user, "SELECT * FROM users WHERE id= '"+id+"' and is_verified = true and is_active = true")

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "user is not verified or actived"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (repository *UserRepository) UpdateUser(c *gin.Context) {

	user := Users{}

	id := c.Param("id")

	// fetch existing data
	err := repository.Db.Get(&user, "SELECT * FROM users WHERE id= '"+id+"' and is_verified = true and is_active = true")

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "user is not verified or actived"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// bind input
	input := Users{}
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// update new password
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(input.Password), 10)

	input.Password = string(hash)

	// merge (manual partial update logic)

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		user.Password = input.Password
	}

	// update data
	_, err = repository.Db.Exec(
		"UPDATE users SET name = ?, email = ?, password = ? WHERE id = ? and is_verified = true and is_active = true",
		user.Name,
		user.Email,
		user.Password,
		id,
	)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)

}
