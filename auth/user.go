package auth

import (
	"database/sql"
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewUser(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}

func (repository *Repository) Users(c *gin.Context) {

	userList, err := repository.GetUsers()

	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userList)

}

func (repository *Repository) GetUsers() ([]Users, error) {

	users := []Users{}

	err := repository.Db.Select(&users, "SELECT * FROM users")

	return users, err
}

func (repository *Repository) UserById(c *gin.Context) {

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

func (repository *Repository) UpdateUser(c *gin.Context) {

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
	hashed, _ := utils.HashPassword(input.Password, 10)
	input.Password = hashed

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
