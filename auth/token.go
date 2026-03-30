package auth

func NewAuthToken (db *sqlx.DB) *Repository {
	return &Repository{Db:db}
}

func (repository *Repository) SaveRefreshToken(userID int, token string, expiry time.Time) error {
    query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`
    _, err := repository.Db.Exec(query, userID, token, expiry)
    return err
}

func (repository *Repository) GetRefreshToken(token string) (*RefreshToken, error) {
    var rt RefreshToken

    query := `
        SELECT * FROM refresh_tokens 
        WHERE token = ? AND expires_at > UTC_TIMESTAMP()
    `
    err := r.Db.Get(&rt, query, token)
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

    c.ShouldBindJSON(&req)

    // access, refresh, err := h.Service.RefreshToken(req.RefreshToken)
    // if err != nil {
    //     c.JSON(401, gin.H{"error": err.Error()})
    //     return
    // }

	tokenHash := HashToken(req.RefreshToken)

	refresh, err := repository.GetRefreshToken(tokenHash)
    if err != nil {
		log.Error(err)
        return errors.New("invalid refresh token")
    }

	// Rotate token (delete old)
    repository.DeleteRefreshToken(refresh.Id)

	// Generate new tokens
    newAccess, _ := utils.GenerateAccessToken(refresh.UserId)
    newRefresh, _ := utils.GenerateRefreshToken(refresh.UserId)

	// Store new refresh token
    newHash := HashToken(newRefresh)
	loc, _ := time.LoadLocation("Asia/Kolkata")
    expiry := time.Now().In(loc).Add(7 * 24 * time.Hour)

	repository.SaveRefreshToken(refresh.UserId, newHash, expiry)

    c.JSON(200, gin.H{
        "access_token":  access,
        "refresh_token": refresh,
    })
}