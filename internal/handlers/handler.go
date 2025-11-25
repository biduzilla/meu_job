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
		Service: s,
		User:    NewUserHandler(s.User, errRsp),
		Auth:    NewAuthHandler(s.Auth, errRsp),
	}
}

func parseID(
	w http.ResponseWriter,
	r *http.Request,
	errRsp errors.ErrorResponseInterface,
) (int64, bool) {
	id, err := utils.ReadIntPathVariable(r, "id")
	if err != nil {
		errRsp.BadRequestResponse(w, r, err)
		return 0, false
	}
	return id, true
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
