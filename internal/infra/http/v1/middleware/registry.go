package middleware

import (
	"companies/internal/company_ctx"
	"companies/internal/domain/interfaces"
	"companies/internal/infra/http/v1/controller"
	"companies/pkg/gen/oapi"
	"companies/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
)

type Registry struct {
	authByJwt mux.MiddlewareFunc
	logging   mux.MiddlewareFunc
}

func (r *Registry) AuthByJwt() mux.MiddlewareFunc {
	return r.authByJwt
}

func newAuthByJwt(authorizer interfaces.UserAuthorizer) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			split := strings.Split(header, "Bearer")
			if len(split) != 2 {
				log.Println("failed to authorize user by secretToken: bad authorization header")
				controller.RespondWithError(rw, "bad authorization header", oapi.Unauthorized, http.StatusUnauthorized)
				return
			}
			token := strings.TrimSpace(split[1])

			err := authorizer.ByJwt(token)
			if err != nil {
				log.Println("failed to authorize user by secretToken", zap.Error(err))

				if errors.Is(err, interfaces.ErrTokenNotValid) {
					controller.RespondWithError(rw, "Token is not valid error: "+err.Error(), oapi.ForbiddenAccess, http.StatusForbidden)
					return
				}

				controller.RespondWithError(rw, "Unknown error occurred while authorization: "+err.Error(), oapi.UnknownAuthError, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func (r *Registry) Logging() mux.MiddlewareFunc {
	return r.logging
}

func newLogging(logger *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger
			xRequestId := company_ctx.RequestId(r.Context())
			if xRequestId != "" {
				log = logger.With(zap.String("xRequestId", xRequestId))
			}

			ctx := company_ctx.WithLogger(r.Context(), log)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func NewRegistry(authorizer interfaces.UserAuthorizer, logger *logger.Logger) *Registry {
	return &Registry{
		authByJwt: newAuthByJwt(authorizer),
		logging:   newLogging(logger),
	}
}
