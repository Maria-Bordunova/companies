package v1

import (
	"companies/internal/infra/http/v1/controller"
	"companies/internal/infra/http/v1/middleware"
	"github.com/gorilla/mux"

	"net/http"
)

func NewRouter(c controller.Controller, mr *middleware.Registry) http.Handler {
	router := mux.NewRouter()

	router.Use(
		mr.Logging(),
	)
	authRouter := router.Methods(http.MethodPost, http.MethodPatch, http.MethodDelete).Subrouter()
	authRouter.Use(
		mr.Logging(),
		mr.AuthByJwt(),
	)

	router.HandleFunc("/companies/{uid}", c.HandleCompanyGetById).Methods(http.MethodGet)
	authRouter.HandleFunc("/companies", c.HandleCompanyCreate).Methods(http.MethodPost)
	authRouter.HandleFunc("/companies/{uid}", c.HandleCompanyUpdateById).Methods(http.MethodPatch)
	authRouter.HandleFunc("/companies/{uid}", c.HandleCompanyDeleteById).Methods(http.MethodDelete)

	return router
}
