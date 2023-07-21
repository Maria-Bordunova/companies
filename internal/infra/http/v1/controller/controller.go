package controller

import (
	"companies/internal/domain/interfaces"
)

type Error struct {
	Error string `json:"error"`
}

type Controller struct {
	companiesRepo interfaces.CompaniesRepo
	eventProducer interfaces.CompanyEventProducer
}

func NewController(
	companiesRepo interfaces.CompaniesRepo,
	eventProducer interfaces.CompanyEventProducer,
) Controller {
	return Controller{
		companiesRepo: companiesRepo,
		eventProducer: eventProducer,
	}
}
