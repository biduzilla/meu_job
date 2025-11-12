package services

import (
	"database/sql"
	"meu_job/internal/config"
	"meu_job/internal/repositories"
)

type Service struct {
	User UserServiceInterface
	Auth AuthServiceInterface
}

func New(db *sql.DB, config config.Config) *Service {
	r := repositories.New(db)
	userService := NewUserService(r.User, db)
	return &Service{
		User: userService,
		Auth: NewAuthService(userService, config),
	}
}
