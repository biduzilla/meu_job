package errors

import (
	"errors"
	"fmt"
	"meu_job/internal/jsonlog"
	"meu_job/utils"
	"meu_job/utils/validator"
	"net/http"
)

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrEditConflict          = errors.New("edit conflict")
	ErrDuplicateEmail        = errors.New("duplicate email")
	ErrDuplicateName         = errors.New("duplicate name")
	ErrDuplicateCNPJ         = errors.New("duplicate CNPJ")
	ErrDuplicatePhone        = errors.New("duplicate phone")
	ErrInvalidData           = errors.New("invalid data")
	ErrInvalidCredentials    = errors.New("invalid authentication credentials")
	ErrInactiveAccount       = errors.New("your user account must be activated to access this resource")
	ErrStartDateAfterEndDate = errors.New("start date must be before end date")
)

type errorResponse struct {
	logger *jsonlog.Logger
}

type ErrorResponseInterface interface {
	NotPermittedResponse(w http.ResponseWriter, r *http.Request)
	AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request)
	InactiveAccountResponse(w http.ResponseWriter, r *http.Request)
	InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request)
	InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request)
	RateLimitExceededResponse(w http.ResponseWriter, r *http.Request)
	ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
	NotFoundResponse(w http.ResponseWriter, r *http.Request)
	MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request)
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
	FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string)
	EditConflictResponse(w http.ResponseWriter, r *http.Request)
	HandlerErrorResponse(w http.ResponseWriter, r *http.Request, err error, v *validator.Validator)
}

func NewErrorResponse(logger *jsonlog.Logger) *errorResponse {
	return &errorResponse{logger: logger}
}

func (e *errorResponse) HandlerErrorResponse(w http.ResponseWriter, r *http.Request, err error, v *validator.Validator) {
	switch {
	case errors.Is(err, ErrInvalidData):
		e.FailedValidationResponse(w, r, v.Errors)

	case errors.Is(err, ErrRecordNotFound):
		e.NotFoundResponse(w, r)

	case errors.Is(err, ErrDuplicateEmail) && v != nil:
		v.AddError("email", "a register with this email address already exists")
		e.FailedValidationResponse(w, r, v.Errors)

	case errors.Is(err, ErrDuplicateName) && v != nil:
		v.AddError("name", "a register with this name already exists")
		e.FailedValidationResponse(w, r, v.Errors)

	case errors.Is(err, ErrDuplicateCNPJ) && v != nil:
		v.AddError("cnpj", "a register with this cnpj already exists")
		e.FailedValidationResponse(w, r, v.Errors)

	case errors.Is(err, ErrDuplicatePhone) && v != nil:
		v.AddError("phone", "a register with this phone number already exists")
		e.FailedValidationResponse(w, r, v.Errors)

	case errors.Is(err, ErrEditConflict):
		e.EditConflictResponse(w, r)

	case errors.Is(err, ErrInactiveAccount):
		e.InactiveAccountResponse(w, r)

	default:
		e.ServerErrorResponse(w, r, err)
	}
}

func (e *errorResponse) NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	e.errorResponse(w, r, http.StatusForbidden, message)
}

func (e *errorResponse) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (e *errorResponse) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	e.errorResponse(w, r, http.StatusForbidden, message)
}

func (e *errorResponse) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *errorResponse) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *errorResponse) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceed"
	e.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (e *errorResponse) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	e.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (e *errorResponse) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	e.errorResponse(w, r, http.StatusNotFound, message)
}

func (e *errorResponse) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	e.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (e *errorResponse) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (e *errorResponse) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	e.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (e *errorResponse) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	e.errorResponse(w, r, http.StatusConflict, message)
}

func (e *errorResponse) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := utils.Envelope{"error": message}
	err := utils.WriteJSON(w, status, env, nil)
	if err != nil {
		e.logError(r, err)
		w.WriteHeader(500)
	}
}

func (e *errorResponse) logError(r *http.Request, err error) {
	e.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}
