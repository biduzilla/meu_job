package models

import (
	"errors"
	"meu_job/utils/validator"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type Role int8

const (
	USER Role = iota + 1
	BUSINESS
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  password
	Phone     string
	Activated bool
	Cod       int
	Role
	BaseModel
}

type UserDTO struct {
	ID    int64  `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserSaveDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type password struct {
	Plaintext *string
	Hash      []byte
}

func (r Role) String() string {
	switch r {
	case USER:
		return "USER"
	case BUSINESS:
		return "BUSINESS"
	default:
		return "UNKNOWN"
	}
}

func ParseRole(s string) Role {
	switch strings.ToUpper(s) {
	case "USER":
		return USER
	case "BUSINESS":
		return BUSINESS
	default:
		return 0
	}
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}

func (u *UserDTO) ToModel() *User {
	return &User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}

func (u *UserSaveDTO) ToModel() (*User, error) {
	user := &User{
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}

	err := user.Password.Set(u.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func (m *User) ValidateUser(v *validator.Validator) {
	v.Check(m.Name != "", "name", "must be provided")
	v.Check(len(m.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(m.Phone != "", "phone", "must be provided")

	ValidateEmail(v, m.Email)

	if m.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *m.Password.Plaintext)
	}

	if m.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
