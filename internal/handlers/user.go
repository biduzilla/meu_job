package handlers

import "meu_job/internal/services"

type UserHandler struct {
	user services.UserServiceInterface
}

type UserHandlerInterface interface {
	ActivateUserHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(
	user services.UserServiceInterface,
) *UserHandler {
	return &UserHandler{
		user: user,
	}
}

func (h* UserHandler) 
