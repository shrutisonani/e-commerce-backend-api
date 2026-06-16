package routes

type Repository struct {
	Db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		DB: db,
	}
}