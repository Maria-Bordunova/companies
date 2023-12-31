package interfaces

import (
	"companies/internal/entity"
	"companies/internal/entity/event"
	"context"
	"github.com/pkg/errors"
)

var ErrStorageNonRetryable = errors.New("non-retryable storage error")
var ErrTokenNotValid = errors.New("token is not valid")

type UserAuthorizer interface {
	ByJwt(tokenString string) error
}

type CompaniesRepo interface {
	Create(ctx context.Context, createParams entity.CreateCompany) error
	FetchByUid(ctx context.Context, uid string) (*entity.Company, error)
	UpdateByUId(ctx context.Context, uid string, updateParams entity.UpdateCompany) error
	DeleteByUId(ctx context.Context, uid string) error
}

type CompanyEventProducer interface {
	Produce(ctx context.Context, uid string, eventType event.EventType) error
}
