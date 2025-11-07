package repositories

import "database/sql"

type Repository struct {
	user UserRepositoryInterface
}

func New(db *sql.DB) *Repository {
	return &Repository{
		user: NewUserRepository(db),
	}
}
