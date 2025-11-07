package services

import (
	"database/sql"
	"errors"
	"meu_job/internal/models"
	"meu_job/internal/repositories"
	"meu_job/utils"
	e "meu_job/utils/errors"
	"meu_job/utils/validator"
)

type UserService struct {
	user repositories.UserRepositoryInterface
	db   *sql.DB
}

type UserServiceInterface interface {
	GetUserByEmail(email string, v *validator.Validator) (*models.User, error)
	ActivateUser(cod int, email string, v *validator.Validator) (*models.User, error)
	Update(user *models.User) error
	GetUserByCodAndEmail(cod int, email string, v *validator.Validator) (*models.User, error)
	RegisterUserHandler(user *models.User, v *validator.Validator) error
	Insert(user *models.User, v *validator.Validator) error
}

func NewUserService(
	userRepository repositories.UserRepositoryInterface,
	db *sql.DB,
) *UserService {
	return &UserService{
		user: userRepository,
		db:   db,
	}
}

func (s *UserService) GetUserByEmail(email string, v *validator.Validator) (*models.User, error) {
	user, err := s.user.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ActivateUser(cod int, email string, v *validator.Validator) (*models.User, error) {
	if models.ValidateEmail(v, email); !v.Valid() {
		return nil, e.ErrInvalidData
	}

	user, err := s.user.GetByCodAndEmail(cod, email)

	if err != nil {
		return nil, err
	}

	user.Activated = true
	user.Cod = 0

	if err = s.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(user *models.User) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		err := s.user.Update(tx, user)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *UserService) GetUserByCodAndEmail(cod int, email string, v *validator.Validator) (*models.User, error) {
	user, err := s.user.GetByCodAndEmail(cod, email)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrRecordNotFound):
			v.AddError("code", "invalid validation code or email")
			return nil, e.ErrInvalidData
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) RegisterUserHandler(user *models.User, v *validator.Validator) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		user.Cod = utils.GenerateRandomCode()
		return s.Insert(user, v)
	})
}

func (s *UserService) Insert(user *models.User, v *validator.Validator) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		if user.ValidateUser(v); !v.Valid() {
			return e.ErrInvalidData
		}

		return s.user.Insert(tx, user)
	})
}

func (s *UserService) Delete(idUser int64) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		return s.user.Delete(tx, idUser)
	})
}
