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

type businessService struct {
	business repositories.BusinessRepositoryInterface
	db       *sql.DB
}

type BusinessServiceInterface interface {
	FindAll(
		name,
		email,
		cnpj string,
		userID int64,
		f filters.Filters,
	) ([]*models.Business, filters.Metadata, error)
	Save(b *models.Business, userID int64, v *validator.Validator) error
	FindByID(id, userID int64) (*models.Business, error)
	Update(b *models.Business, userID int64, v *validator.Validator) error
	Delete(id, userID int64) error
	AddUserInBusiness(businessID, userID, userLogadoID int64) error
}

func NewBusinessService(
	businessRepository repositories.BusinessRepositoryInterface,
	db *sql.DB,
) *businessService {
	return &businessService{
		business: businessRepository,
		db:       db,
	}
}

func (s *businessService) FindAll(
	name,
	email,
	cnpj string,
	userID int64,
	f filters.Filters,
) ([]*models.Business, filters.Metadata, error) {
	return s.business.GetAll(name, email, cnpj, userID, f)
}

func (s *businessService) AddUserInBusiness(businessID, userID, userLogadoID int64) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		return s.business.AddUserInBusiness(businessID, userID, userLogadoID, tx)
	})
}

func (s *businessService) Save(b *models.Business, userID int64, v *validator.Validator) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		b.ValidateBusiness(v)
		if !v.Valid() {
			return errors.ErrInvalidData
		}

		return s.business.Insert(b, userID, tx)
	})
}

func (s *businessService) FindByID(id, userID int64) (*models.Business, error) {
	return s.business.GetByID(id, userID)
}

func (s *businessService) Update(b *models.Business, userID int64, v *validator.Validator) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		if b.ValidateBusiness(v); !v.Valid() {
			return errors.ErrInvalidData
		}

		return s.business.Update(b, userID, tx)
	})
}

func (s *businessService) Delete(id, userID int64) error {
	return utils.RunInTx(s.db, func(tx *sql.Tx) error {
		return s.business.Delete(id, userID, tx)
	})
}
