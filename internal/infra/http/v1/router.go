package v1

import (
	controller "companies/internal/infra/http/v1/controller"
	"companies/internal/infra/http/v1/middleware"
	"github.com/gorilla/mux"

	"net/http"
)

func NewRouter(c controller.Controller, mr *middleware.Registry) http.Handler {
	router := mux.NewRouter()

	routerAuth := mux.NewRouter()
	routerAuth.Use(
		mr.AuthByJwt(),
	)

	routerAuth.HandleFunc("/companies", c.HandleCompanyCreate).Methods(http.MethodPost)
	router.HandleFunc("/companies/{uid}", c.HandleCompanyGetById).Methods(http.MethodGet)
	routerAuth.HandleFunc("/companies/{uid}", c.HandleCompanyUpdateById).Methods(http.MethodPatch)
	routerAuth.HandleFunc("/companies/{uid}", c.HandleCompanyDeleteById).Methods(http.MethodDelete)

	return router
}
