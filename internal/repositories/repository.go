package repositories

import "database/sql"

type Repository struct {
	User UserRepositoryInterface
}

func New(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
