package v1

import (
	controller "companies/internal/infra/http/v1/handler"
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

	router.HandleFunc("/companies/{companyId}", c.HandleCompanyGetById).Methods(http.MethodGet)
	routerAuth.HandleFunc("/companies/{companyId}", c.HandleCompanyCreate).Methods(http.MethodPost)
	routerAuth.HandleFunc("/companies/{companyId}", c.HandleCompanyUpdateById).Methods(http.MethodPatch)
	routerAuth.HandleFunc("/companies/{companyId}", c.HandleCompanyDeleteById).Methods(http.MethodDelete)

	return router
}
