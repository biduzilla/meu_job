package handlers

import (
	"database/sql"
	"meu_job/internal/config"
	"meu_job/internal/services"
	"meu_job/utils"
	"meu_job/utils/errors"
	"net/http"
)

type Handler struct {
	User    UserHandlerInterface
	Auth    AuthHandlerInterface
	Service *services.Service
}

func NewHandler(
	db *sql.DB,
	errRsp errors.ErrorResponseInterface,
	config config.Config,
) *Handler {
	s := services.New(db, config)

	return &Handler{
		User: NewUserHandler(s.User),
		Auth: NewAuthHandler(s.Auth, errRsp),
	}
}

func respond(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data utils.Envelope,
	headers http.Header,
	errRsp errors.ErrorResponseInterface,
) {
	err := utils.WriteJSON(w, status, data, headers)
	if err != nil {
		errRsp.ServerErrorResponse(w, r, err)
	}
}
