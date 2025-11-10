package services

import (
	"database/sql"
	"meu_job/internal/repositories"
)

type Service struct {
	User UserServiceInterface
}

func New(db *sql.DB) *Service {
	r := repositories.New(db)
	return &Service{
		User: NewUserService(r.User, db),
	}
}
