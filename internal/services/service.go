package services

import (
	"database/sql"
	"meu_job/internal/repositories"
)

type Service struct {
	user UserServiceInterface
}

func New(db *sql.DB) *Service {
	r := repositories.New(db)
	return &Service{
		user: NewUserService(r.User, db),
	}
}
