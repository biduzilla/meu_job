package handlers

import (
	"meu_job/internal/models"
	"meu_job/internal/services"
	"meu_job/utils"
	e "meu_job/utils/errors"
	"meu_job/utils/validator"
	"net/http"
)

type UserHandler struct {
	user   services.UserServiceInterface
	errRsp e.ErrorResponseInterface
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

func (h *UserHandler) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Cod   int    `json:"cod"`
		Email string `json:"email"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	user, err := h.user.ActivateUser(
		input.Cod,
		input.Email,
		v,
	)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(
		w,
		r,
		http.StatusOK,
		utils.Envelope{"user": user.ToDTO()},
		nil,
		h.errRsp,
	)
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO models.UserSaveDTO
	if err := utils.ReadJSON(w, r, &userDTO); err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	user, err := userDTO.ToModel()
	if err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	err = h.user.RegisterUserHandler(user, v)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(
		w,
		r,
		http.StatusCreated,
		utils.Envelope{"user": user.ToDTO()},
		nil,
		h.errRsp,
	)
}
