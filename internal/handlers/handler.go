package handlers

import (
	"meu_job/internal/services"
	"meu_job/utils"
	"meu_job/utils/errors"
	"net/http"
)

type Handler struct {
	User    UserHandlerInterface
	Service *services.Service
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
