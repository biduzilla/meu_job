package services

import (
	"database/sql"
	"meu_job/internal/config"
	"meu_job/internal/models"
	"meu_job/internal/repositories"
	"meu_job/utils/validator"
)

type Service struct {
	User     UserServiceInterface
	Auth     AuthServiceInterface
	Business BusinessServiceInterface
}

type GenericServiceInterface[
	T models.ModelInterface[D],
	D any,
] interface {
	Save(entity *T, userID int64, v *validator.Validator) error
	FindByID(id, userID int64) (*T, error)
	Update(entity *T, userID int64, v *validator.Validator) error
	Delete(id, userID int64) error
}

func New(db *sql.DB, config config.Config) *Service {
	r := repositories.New(db)
	userService := NewUserService(r.User, db)
	return &Service{
		User:     userService,
		Auth:     NewAuthService(userService, config),
		Business: NewBusinessService(r.Business, db),
	}
}
