package handlers

import (
	"meu_job/internal/contexts"
	"meu_job/internal/models"
	"meu_job/internal/models/filters"
	"meu_job/internal/services"
	"meu_job/utils"
	e "meu_job/utils/errors"
	"meu_job/utils/validator"
	"net/http"
)

type businessHandler struct {
	GenericHandlerInterface[models.Business, models.BusinessDTO]
	business services.BusinessServiceInterface
	errRsp   e.ErrorResponseInterface
}

type BusinessHandlerInterface interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	AddUserInBusiness(w http.ResponseWriter, r *http.Request)
	GenericHandlerInterface[
		models.Business,
		models.BusinessDTO,
	]
}

func NewBusinessHandler(
	business services.BusinessServiceInterface,
	errRsp e.ErrorResponseInterface,
) *businessHandler {
	return &businessHandler{
		GenericHandlerInterface: NewGenericHandler(business, errRsp),
		business:                business,
		errRsp:                  errRsp,
	}
}

func (h *businessHandler) AddUserInBusiness(w http.ResponseWriter, r *http.Request) {
	businessID, err := utils.ReadIntPathVariable(r, "businessID")
	if err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	userID, err := utils.ReadIntPathVariable(r, "userID")
	if err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	userLogado := contexts.ContextGetUser(r)

	err = h.business.AddUserInBusiness(businessID, userID, userLogado.ID)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	respond(
		w,
		r,
		http.StatusNoContent,
		nil,
		nil,
		h.errRsp,
	)
}

func (h *businessHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	var input struct {
		name, cnpj, email string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.name = utils.ReadString(qs, "name", "")
	input.cnpj = utils.ReadString(qs, "cnpj", "")
	input.email = utils.ReadString(qs, "email", "")
	input.Filters.Page = utils.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = utils.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		h.errRsp.HandlerErrorResponse(w, r, e.ErrInvalidData, v)
		return
	}

	user := contexts.ContextGetUser(r)
	business, metadata, err := h.business.FindAll(
		input.name,
		input.email,
		input.cnpj,
		user.ID,
		input.Filters,
	)

	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"business": business, "metadata": metadata}, nil, h.errRsp)
}

// func (h *businessHandler) FindByID(w http.ResponseWriter, r *http.Request) {
// 	id, ok := parseID(w, r, h.errRsp)
// 	if !ok {
// 		return
// 	}

// 	user := contexts.ContextGetUser(r)
// 	business, err := h.business.FindByID(id, user.ID)
// 	if err != nil {
// 		h.errRsp.HandlerErrorResponse(w, r, err, nil)
// 		return
// 	}
// 	respond(w, r, http.StatusOK, utils.Envelope{"business": business.ToDTO()}, nil, h.errRsp)
// }

// func (h *businessHandler) Save(w http.ResponseWriter, r *http.Request) {
// 	var dto models.BusinessDTO
// 	if err := utils.ReadJSON(w, r, &dto); err != nil {
// 		h.errRsp.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	v := validator.New()
// 	user := contexts.ContextGetUser(r)
// 	model := dto.ToModel()
// 	if model.User == nil {
// 		model.User = user
// 	}

// 	if err := h.business.Save(model, v); err != nil {
// 		h.errRsp.HandlerErrorResponse(w, r, err, v)
// 		return
// 	}

// 	respond(w, r, http.StatusCreated, utils.Envelope{"business": model.ToDTO()}, nil, h.errRsp)
// }

// func (h *businessHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	var dto models.BusinessDTO
// 	if err := utils.ReadJSON(w, r, &dto); err != nil {
// 		h.errRsp.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	v := validator.New()
// 	user := contexts.ContextGetUser(r)
// 	model := dto.ToModel()
// 	if model.User == nil {
// 		model.User = user
// 	}

// 	if err := h.business.Update(model, v); err != nil {
// 		h.errRsp.HandlerErrorResponse(w, r, err, v)
// 		return
// 	}

// 	respond(w, r, http.StatusOK, utils.Envelope{"business": model}, nil, h.errRsp)
// }

// func (h *businessHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	id, ok := parseID(w, r, h.errRsp)
// 	if !ok {
// 		return
// 	}

// 	user := contexts.ContextGetUser(r)
// 	if err := h.business.Delete(id, user.ID); err != nil {
// 		h.errRsp.HandlerErrorResponse(w, r, err, nil)
// 		return
// 	}

// 	respond(
// 		w,
// 		r,
// 		http.StatusNoContent,
// 		nil,
// 		nil,
// 		h.errRsp,
// 	)
// }
