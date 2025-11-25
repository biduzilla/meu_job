package services

import (
	"database/sql"
	"meu_job/internal/models"
	"meu_job/internal/models/filters"
	"meu_job/internal/repositories"
	"meu_job/utils"
	"meu_job/utils/errors"
	"meu_job/utils/validator"
)

type BusinessService struct {
	business repositories.BusinessRepositoryInterface
	db       *sql.DB
}

type BusinessServiceInterface interface {
	Save(b *models.Business, v *validator.Validator) error
	FindByID(id, userID int64) (*models.Business, error)
	FindAll(
		name,
		email,
		cnpj string,
		userID int64,
		f filters.Filters,
	) ([]*models.Business, filters.Metadata, error)
	Update(b *models.Business, v *validator.Validator) error
	Delete(id, userID int64) error
}

func NewBusinessService(
	businessRepository repositories.BusinessRepositoryInterface,
	db *sql.DB,
) *BusinessService {
	return &BusinessService{
		business: businessRepository,
		db:       db,
	}
}

func (s *BusinessService) Save(b *models.Business, v *validator.Validator) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		b.ValidateBusiness(v)
		if !v.Valid() {
			return errors.ErrInvalidData
		}

		return s.business.Insert(b, b.User.ID, tx)
	})
}

func (s *BusinessService) FindByID(id, userID int64) (*models.Business, error) {
	return s.business.GetByID(id, userID)
}

func (s *BusinessService) FindAll(
	name,
	email,
	cnpj string,
	userID int64,
	f filters.Filters,
) ([]*models.Business, filters.Metadata, error) {
	return s.business.GetAll(name, email, cnpj, userID, f)
}

func (s *BusinessService) Update(b *models.Business, v *validator.Validator) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		if b.ValidateBusiness(v); !v.Valid() {
			return errors.ErrInvalidData
		}

		return s.business.Update(b, b.User.ID)
	})
}

func (s *BusinessService) Delete(id, userID int64) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		return s.business.Delete(id, userID)
	})
}
