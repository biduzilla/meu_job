package repositories

import "database/sql"

type Repository struct {
	User     UserRepositoryInterface
	Business BusinessRepositoryInterface
}

func New(db *sql.DB) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Business: NewBusinessRepository(db),
	}
}
