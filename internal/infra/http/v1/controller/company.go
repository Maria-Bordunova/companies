package controller

import (
	"companies/internal/ctx"
	"companies/internal/domain/interfaces"
	"companies/internal/entity"
	"companies/pkg/gen/oapi"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

func (c *Controller) HandleCompanyCreate(rw http.ResponseWriter, r *http.Request) {
	log := ctx.Logger(r.Context())

	params := &oapi.CompanyInput{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		RespondWithError(rw, err.Error(), oapi.ValidationError, http.StatusBadRequest)
		return
	}
	company := entity.CreateCompany{
		UId:         params.Id,
		Name:        params.Name,
		Description: params.Description,
		Employees:   params.Employees,
		Registered:  params.Registered,
		Type:        (entity.CompanyType)(params.Type), // TODO better implementation with strict check
	}

	err = c.companiesRepo.Create(r.Context(), company)
	if err != nil {
		log.
			With("company uid", company.UId).
			With("company name", company.Name).
			With(zap.Error(err)).
			Error("error occurred when creating company")

		if errors.Is(err, interfaces.ErrStorageNonRetryable) {
			RespondWithError(rw, "Company update error: "+err.Error(), oapi.CompanyNotFound, http.StatusNotFound)
			return
		}

		RespondWithError(rw, "Company create error: "+err.Error(), oapi.UnknownStorageError, http.StatusInternalServerError)
		return
	}

	createdCompany, err := c.companiesRepo.FetchByUid(r.Context(), company.UId) // TODO move it to domain layer
	if err != nil {
		log.
			With("company uid", company.UId).
			With(zap.Error(err)).
			Error("error occurred when fetching company by uid")

		RespondWithError(rw, "Company fetch error: "+err.Error(), oapi.UnknownStorageError, http.StatusInternalServerError)
		return
	}
	if createdCompany == nil {
		RespondWithError(rw, "Company not found", oapi.CompanyNotFound, http.StatusInternalServerError)
		return
	}

	RespondWithData(rw, convertCompanyToView(*createdCompany))
}

func (c *Controller) HandleCompanyGetById(rw http.ResponseWriter, r *http.Request) {
	log := ctx.Logger(r.Context())

	uid, err := parseUid(r)
	if err != nil {
		RespondWithError(rw, "Bad params: "+err.Error(), oapi.ValidationError, http.StatusBadRequest)
		return
	}

	company, err := c.companiesRepo.FetchByUid(r.Context(), uid)
	if err != nil {
		log.
			With("company uid", uid).
			With(zap.Error(err)).
			Error("error occurred when fetching company by uid")

		RespondWithError(rw, "Company fetch error: "+err.Error(), oapi.UnknownStorageError, http.StatusInternalServerError)
		return
	}
	if company == nil {
		RespondWithError(rw, "Company not found", oapi.CompanyNotFound, http.StatusNotFound)
		return
	}

	RespondWithData(rw, convertCompanyToView(*company))
}

func (c *Controller) HandleCompanyDeleteById(rw http.ResponseWriter, r *http.Request) {
	log := ctx.Logger(r.Context())

	uid, err := parseUid(r)
	if err != nil {
		RespondWithError(rw, "Bad params: "+err.Error(), oapi.ValidationError, http.StatusBadRequest)
		return
	}

	err = c.companiesRepo.DeleteByUId(r.Context(), uid)
	if err != nil {
		log.
			With("company uid", uid).
			With(zap.Error(err)).
			Error("error occurred when deleting company by uid")

		if errors.Is(err, interfaces.ErrStorageNonRetryable) {
			RespondWithError(rw, "Company delete error: "+err.Error(), oapi.CompanyNotFound, http.StatusNotFound)
			return
		}

		RespondWithError(rw, "Company delete error: "+err.Error(), oapi.UnknownStorageError, http.StatusInternalServerError)
		return
	}

	c.RespondWithCode(rw, http.StatusNoContent)
}

func (c *Controller) HandleCompanyUpdateById(rw http.ResponseWriter, r *http.Request) {
	log := ctx.Logger(r.Context())

	uid, err := parseUid(r)
	if err != nil {
		RespondWithError(rw, "Bad params: "+err.Error(), oapi.ValidationError, http.StatusBadRequest)
		return
	}

	params := &oapi.CompanyUpdate{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(params)
	if err != nil {
		RespondWithError(rw, err.Error(), oapi.ValidationError, http.StatusBadRequest)
		return
	}
	company := entity.UpdateCompany{
		Name:        params.Name,
		Description: params.Description,
		Employees:   params.Employees,
		Registered:  params.Registered,
		Type:        (*entity.CompanyType)(params.Type), // TODO better implementation with strict check
	}

	err = c.companiesRepo.UpdateByUId(r.Context(), uid, company)
	if err != nil {
		log.
			With("company uid", uid).
			With(zap.Error(err)).
			Error("error occurred when updating company by uid")

		if errors.Is(err, interfaces.ErrStorageNonRetryable) {
			RespondWithError(rw, "Company update error: "+err.Error(), oapi.CompanyNotFound, http.StatusNotFound)
			return
		}

		RespondWithError(rw, "Company update error: "+err.Error(), oapi.UnknownStorageError, http.StatusInternalServerError)
		return
	}
	updatedCompany, err := c.companiesRepo.FetchByUid(r.Context(), uid) // TODO move it to domain layer
	if err != nil {
		log.
			With("company uid", uid).
			With(zap.Error(err)).
			Error("error occurred when fetching company by uid")

		RespondWithError(rw, "Company fetch error: "+err.Error(), oapi.UnknownStorageError, http.StatusInternalServerError)
		return
	}
	if updatedCompany == nil {
		RespondWithError(rw, "Company not found", oapi.CompanyNotFound, http.StatusNotFound)
		return
	}

	RespondWithData(rw, convertCompanyToView(*updatedCompany))
}

func parseUid(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	uid, ok := vars["uid"]
	if !ok || uid == "" {
		return "", errors.New("user Id required")
	}
	return uid, nil
}

func convertCompanyToView(company entity.Company) oapi.Company {
	view := oapi.Company{
		Description: company.Description,
		Employees:   company.Employees,
		Id:          company.UId,
		Name:        company.Name,
		Registered:  company.Registered,
		Type:        oapi.CompanyType(company.Type), // TODO better implementation with strict check
	}

	return view
}
