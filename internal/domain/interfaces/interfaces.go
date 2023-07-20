package interfaces

import (
	"companies/internal/entity"
	"context"
)

type UserAuthorizer interface {
	AuthByJwt(tokenString string) error
}

type CompaniesRepo interface {
	Create(ctx context.Context, createParams entity.CreateCompany) (*entity.Company, error)
	GetById(ctx context.Context, uid string) (*entity.Company, error)
	UpdateById(ctx context.Context, uid string, updateParams entity.UpdateCompany) error
	DeleteById(ctx context.Context, uid string) error
}
