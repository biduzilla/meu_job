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

type businessRepository struct {
	db *sql.DB
}

func NewBusinessRepository(db *sql.DB) *businessRepository {
	return &businessRepository{
		db: db,
	}
}

type BusinessRepositoryInterface interface {
	GetByID(id int64, userID int64) (*models.Business, error)
	GetAll(
		name,
		email,
		cnpj string,
		userID int64,
		f filters.Filters,
	) ([]*models.Business, filters.Metadata, error)
	Insert(business *models.Business, userID int64, tx *sql.Tx) error
	Update(business *models.Business, userID int64, tx *sql.Tx) error
	Delete(id, userID int64, tx *sql.Tx) error
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

func scanBusinessPage(r *sql.Rows, totalRecords *int, business *models.Business) error {
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
	)
}

func scanBusiness(r *sql.Row, business *models.Business) error {
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

func (r *businessRepository) GetByID(id int64, userID int64) (*models.Business, error) {
	query := fmt.Sprintf(`
	select 
		%s
	from business b
	left join users u on b.user_id = u.id
	where 
		b.id = $1
		and b.deleted = false
		and exists (
			select 1
			from business_users bu
			where 
				bu.business_id = b.id
				and bu.user_id = $2
		)
	`, SQLSelectDataBusiness)

	business := models.Business{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id, userID)
	if err := scanBusiness(row, &business); err != nil {
		return nil, err
	}

	return &business, nil
}

func (r *businessRepository) GetAll(
	name,
	email,
	cnpj string,
	userID int64,
	f filters.Filters,
) ([]*models.Business, filters.Metadata, error) {
	query := fmt.Sprintf(`
		select
			count(*) over(),
			%s
		from business b
		join business_users bu on bu.business_id = b.id
		where
			bu.user_id = $4
			and (to_tsvector('simple', b.name) @@ plainto_tsquery('simple', $1) OR $1 = '')
			and (to_tsvector('simple', b.email) @@ plainto_tsquery('simple', $2) OR $2 = '')
			and (to_tsvector('simple', b.cnpj) @@ plainto_tsquery('simple', $3) OR $3 = '')
			and b.deleted = false
		order by %s %s, b.id
		limit $5 offset $6
	`, SQLSelectDataBusiness, f.SortColumn(), f.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, email, cnpj, userID, f.Limit(), f.Offset()}

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	businessList := []*models.Business{}

	for rows.Next() {
		business := models.Business{}

		err := scanBusinessPage(rows, &totalRecords, &business)
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

func (r *businessRepository) Insert(business *models.Business, userID int64, tx *sql.Tx) error {
	query := `
	insert into business (
		name,
		cnpj,
		email,
		phone,
		created_by
	)
	values ($1,$2,$3,$4,$5)
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
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&business.ID,
		&business.CreatedAt,
		&business.Version,
	)

	if err != nil {
		return r.uniqueErrors(err)
	}

	_, err = tx.ExecContext(ctx, `
    INSERT INTO business_users (business_id, user_id)
    VALUES ($1, $2)
	`, business.ID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *businessRepository) Update(business *models.Business, userID int64, tx *sql.Tx) error {
	query := `
	update business
	set
		name = $1,
		cnpj = $2,
		email = $3,
		phone = $4,
		updated_by = $5,	
		updated_at=now(),
		version = version + 1	
	where
		id = $6
		and exists(
			select 1 
			from business_users bu 
			where 
				bu.business_id=$6 
				and bu.user_id=$5
		)
		and deleted = false
		and version = $7
	returning version
	`

	args := []any{
		business.Name,
		business.CNPJ,
		business.Email,
		business.Phone,
		userID,
		business.ID,
		business.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&business.Version,
	)

	return r.uniqueErrors(err)
}

func (r *businessRepository) Delete(id, userID int64, tx *sql.Tx) error {
	query := `
		update business
		set 
			deleted = true,
			updated_by = $2,
			updated_at = NOW(),
			version = version + 1
		where id = $1
		and exists (
			select 1 from business_users 
			where business_id = $1 and user_id = $2
		)
		and deleted = false
		returning id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var returnedID int64

	err := tx.QueryRowContext(ctx, query, id, userID).Scan(&returnedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (r *businessRepository) uniqueErrors(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Constraint {
		case "unique_business_name_deleted":
			return e.ErrDuplicateName
		case "unique_business_cnpj_deleted":
			return e.ErrDuplicateCNPJ
		case "unique_business_email_deleted":
			return e.ErrDuplicateEmail
		}
	}
	return err
}
