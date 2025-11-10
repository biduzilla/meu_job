package handlers

import (
	"database/sql"
	"meu_job/internal/services"
	"meu_job/utils"
	"meu_job/utils/errors"
	"net/http"
)

type Handler struct {
	User    UserHandlerInterface
	Service *services.Service
}

func NewHandler(
	db *sql.DB,
	errRsp errors.ErrorResponseInterface,
) *Handler {
	s := services.New(db)

	return &Handler{
		User: NewUserHandler(s.User),
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
