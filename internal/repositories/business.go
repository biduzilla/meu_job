package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"meu_job/internal/models"
	"meu_job/internal/models/filters"
	e "meu_job/utils/errors"
	"time"

	"github.com/lib/pq"
)

type BusinessRepository struct {
	db *sql.DB
}

func NewBusinessRepository(db *sql.DB) *BusinessRepository {
	return &BusinessRepository{
		db: db,
	}
}

const SQLSelectDataBusiness = `
		b.id,
		b.name,
		b.cnpj,
		b.email,
		b.phone,
		b.version,
		b.deleted,
		b.created_by,
		b.created_at,
		b.updated_by,
		b.updated_at
	`

func ScanBusinessPage(r *sql.Rows, totalRecords *int, business *models.Business) error {
	return r.Scan(
		&totalRecords,
		&business.ID,
		&business.Name,
		&business.CNPJ,
		&business.Email,
		&business.Phone,
		&business.Version,
		&business.Deleted,
		&business.CreatedBy,
		&business.CreatedAt,
		&business.UpdatedBy,
		&business.UpdatedAt,
		&business.User.ID,
		&business.User.Name,
		&business.User.Phone,
		&business.User.Email,
		&business.User.Cod,
		&business.User.Password.Hash,
		&business.User.Activated,
		&business.User.Version,
		&business.User.CreatedBy,
		&business.User.CreatedAt,
		&business.User.UpdatedBy,
		&business.User.UpdatedAt,
	)
}

func ScanBusiness(r *sql.Row, business *models.Business) error {
	err := r.Scan(
		&business.ID,
		&business.Name,
		&business.CNPJ,
		&business.Email,
		&business.Phone,
		&business.Version,
		&business.Deleted,
		&business.CreatedBy,
		&business.CreatedAt,
		&business.UpdatedBy,
		&business.UpdatedAt,
		&business.User.ID,
		&business.User.Name,
		&business.User.Phone,
		&business.User.Email,
		&business.User.Cod,
		&business.User.Password.Hash,
		&business.User.Activated,
		&business.User.Version,
		&business.User.CreatedBy,
		&business.User.CreatedAt,
		&business.User.UpdatedBy,
		&business.User.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return e.ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (r *BusinessRepository) GetByID(id int64, userID int64) (*models.Business, error) {
	query := fmt.Sprintf(`
	select 
		%s,
		%s
	from business b
	left join users u on b.user_id = u.id
	where 
		b.id = $1
		and b.user_id = $2
		and b.deleted = false
	`, SQLSelectDataBusiness,
		SQLSelectDataUser,
	)

	business := models.Business{
		User: &models.User{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id, userID)
	if err := ScanBusiness(row, &business); err != nil {
		return nil, err
	}

	return &business, nil
}

func (r *BusinessRepository) GetAll(
	name,
	email,
	cnpj string,
	userID int64,
	f filters.Filters,
) ([]*models.Business, filters.Metadata, error) {
	query := fmt.Sprintf(`
	select
		count(b) over,
		%s,
		%s
		from business b
		left join users u on b.user_id = u.id
		where 
			(to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
			and (to_tsvector('simple', email) @@ plainto_tsquery('simple', $2) OR $2 = '')
			and (to_tsvector('simple', cnpj) @@ plainto_tsquery('simple', $3) OR $3 = '')
			and b.user_id = $4
			and b.deleted = false
		order by 
			%s %s,
			id asc
		limit $5 offset $6			
	`,
		SQLSelectDataBusiness,
		SQLSelectDataUser,
		f.SortColumn(),
		f.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, email, cnpj, f.Limit(), f.Offset()}

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	businessList := []*models.Business{}

	for rows.Next() {
		business := models.Business{
			User: &models.User{},
		}

		err := ScanBusinessPage(rows, &totalRecords, &business)
		if err != nil {
			return nil, filters.Metadata{}, err
		}

		businessList = append(businessList, &business)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metaData := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)
	return businessList, metaData, nil
}

func (r *BusinessRepository) Insert(business *models.Business, userID int64, tx *sql.Tx) error {
	query := `
	insert into(
		name,
		cnpj,
		email,
		phone,
		user_id,
		created_by	
	)
	values ($1,$2,$3,$4,$5,$6)
	returning 
		id, 
		created_at,
		version
	`

	args := []any{
		business.Name,
		business.CNPJ,
		business.Email,
		business.Phone,
		userID,
		userID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&business.ID,
		&business.CreatedAt,
		&business.Version,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "unique_user_business_name":
				return e.ErrDuplicateName
			case "unique_user_business_cnpj":
				return e.ErrDuplicateCNPJ
			case "unique_user_business_email":
				return e.ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (r *BusinessRepository) Update(business *models.Business, userID int64) error {
	query := `
	update business
	set
		name = $1,
		cnpj = $2,
		email = $3,
		phone = $4,
		updated_by = $5,	
		updated_at = $6,
		version = version + 1	
	)
	where
		id = $7
		and user_id = $8
		and deleted = false
		and version = $9
	returning version
	`

	args := []any{
		business.Name,
		business.CNPJ,
		business.Email,
		business.Phone,
		userID,
		time.Now(),
		business.ID,
		userID,
		business.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&business.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		} else if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "unique_user_business_name":
				return e.ErrDuplicateName
			case "unique_user_business_cnpj":
				return e.ErrDuplicateCNPJ
			case "unique_user_business_email":
				return e.ErrDuplicateEmail
			}
		}

		return err
	}
	return nil
}

func (r *BusinessRepository) Delete(id, userID int64) error {
	query := `
	UPDATE business
	SET
		deleted = true
	WHERE
		id = $1
		AND user_id = $2
		AND deleted = false
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return e.ErrRecordNotFound
	}

	return nil
}
