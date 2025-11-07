package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"meu_job/internal/models"
	e "meu_job/utils/errors"
	"time"

	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

const SqlSelectUser = `
	SELECT 
		id, 
		created_at, 
		name, 
		phone, 
		email, 
		cod, 
		password_hash, 
		activated, 
		version
	FROM users
`

func NewUserRepository(
	db *sql.DB,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) getUserByQuery(query string, args ...any) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Cod,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (r *UserRepository) GetByCodAndEmail(cod int, email string) (*models.User, error) {
	query := fmt.Sprintf(`
	%s
	WHERE 
		email = $1 
		AND deleted = false 
		AND cod = $2
	`, SqlSelectUser)

	return r.getUserByQuery(query, email, cod)
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := fmt.Sprintf(`
	%s
	WHERE 
		id = $1 
		AND deleted = false
	`, SqlSelectUser)
	return r.getUserByQuery(query, id)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := fmt.Sprintf(`
	%s
	WHERE 
		email = $1 
		AND deleted = false
	`, SqlSelectUser)
	return r.getUserByQuery(query, email)
}

func (r *UserRepository) Insert(tx *sql.Tx, user *models.User) error {
	query := `
	INSERT INTO users (name, email, phone,cod, password_hash, activated,deleted)
	VALUES ($1, $2, $3, $4, $5, $6,false)
	RETURNING id, created_at, version
	`
	args := []any{
		user.Name,
		user.Email,
		user.Phone,
		user.Cod,
		user.Password.Hash,
		user.Activated,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "users_email_key":
				return e.ErrDuplicateEmail
			case "users_phone_key":
				return e.ErrDuplicatePhone
			}
		}

		return err
	}

	return nil
}

func (r *UserRepository) UpdateCodByEmail(tx *sql.Tx, user *models.User) error {
	query := `
	UPDATE users SET
	cod = $1
	WHERE id = $2 AND version = $3
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Cod, user.ID, user.Version).Scan(
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return e.ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

func (r *UserRepository) Update(tx *sql.Tx, user *models.User) error {
	query := `
	UPDATE users SET 
		name = $1,
		email = $2,
		cod = $3, 
		phone = $4, 
		password_hash = $5,
		activated = $6,
		version = version + 1
	WHERE 
		id = $7 
		AND version = $8
	RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Cod,
		user.Phone,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&user.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "users_email_key":
				return e.ErrDuplicateEmail
			case "users_phone_key":
				return e.ErrDuplicatePhone
			}
		}

		return err
	}
	return nil
}

func (r *UserRepository) Delete(tx *sql.Tx, user *models.User) error {
	query := `
	UPDATE users set
	deleted = true
	where id = $1 AND version = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, user.ID, user.Version)
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
