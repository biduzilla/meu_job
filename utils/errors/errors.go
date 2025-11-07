package errors

import "errors"

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrEditConflict          = errors.New("edit conflict")
	ErrDuplicateEmail        = errors.New("duplicate email")
	ErrDuplicateName         = errors.New("duplicate name")
	ErrDuplicatePhone        = errors.New("duplicate phone")
	ErrInvalidData           = errors.New("invalid data")
	ErrInvalidCredentials    = errors.New("invalid authentication credentials")
	ErrInactiveAccount       = errors.New("your user account must be activated to access this resource")
	ErrStartDateAfterEndDate = errors.New("start date must be before end date")
)
