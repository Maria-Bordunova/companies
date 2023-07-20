package controller

import (
	"companies/internal/domain/interfaces"
)

type Error struct {
	Error string `json:"error"`
}

type Controller struct {
	companiesRepo interfaces.CompaniesRepo
}

func NewController(
	companiesRepo interfaces.CompaniesRepo,
) Controller {
	return Controller{
		companiesRepo: companiesRepo,
	}
}
