package controller

import (
	"companies/pkg/gen/oapi"
	"encoding/json"
	"net/http"
)

func RespondWithData(rw http.ResponseWriter, response interface{}) {
	RespondWithDataAndCode(rw, response, http.StatusOK)
}

func RespondWithError(rw http.ResponseWriter, errMsg string, errCode oapi.ErrorCode, statusCode int) {
	RespondWithDataAndCode(rw, oapi.Error{Error: errMsg, Code: errCode}, statusCode)
}

func RespondWithDataAndCode(rw http.ResponseWriter, response interface{}, code int) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(response)
}

func (c Controller) RespondWithCode(rw http.ResponseWriter, code int) {
	rw.WriteHeader(code)
}

func (c Controller) RespondWithCodeOK(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
}
