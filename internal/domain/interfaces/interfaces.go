package interfaces

import (
	"companies/internal/entity"
	"context"
	"github.com/pkg/errors"
)

var ErrStorageNonRetryable = errors.New("non-retryable storage error")

type UserAuthorizer interface {
	AuthByJwt(tokenString string) error
}

type CompaniesRepo interface {
	Create(ctx context.Context, createParams entity.CreateCompany) error
	FetchByUid(ctx context.Context, uid string) (*entity.Company, error)
	UpdateByUId(ctx context.Context, uid string, updateParams entity.UpdateCompany) error
	DeleteByUId(ctx context.Context, uid string) error
}
