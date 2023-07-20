package middleware

import (
	"companies/internal/domain/interfaces"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type Registry struct {
	authByJwt mux.MiddlewareFunc
}

func (r *Registry) AuthByJwt() mux.MiddlewareFunc {
	return r.authByJwt
}

func NewRegistry(authorizer interfaces.UserAuthorizer) *Registry {
	return &Registry{
		authByJwt: newAuthByJwt(authorizer),
	}
}

func newAuthByJwt(authorizer interfaces.UserAuthorizer) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			split := strings.Split(header, "Bearer")
			if len(split) != 2 {
				log.Println("failed to authorize project by accessToken: bad authorization header")

				//v1.RespondWithError(w, http.StatusUnauthorized, oapi.Error{
				//	Code:    v1.ControlUnauthorized,
				//	Message: "Bad authorization header",
				//	Type:    oapi.Request,
				//})
				return
			}
			//secretAccessToken := removeTestPrefix(strings.TrimSpace(split[1]))
			//
			//project, err := authorizer.BySecretKey(r.Context(), secretAccessToken)
			//
			//if err != nil {
			//	log.Errorw("failed to authorize project by accessToken", zap.Error(err))
			//	handleErr(w, err)
			//	return
			//}
			//
			//if project == nil {
			//	handleProjectNotFound(w)
			//	return
			//}
			//
			//r = r.WithContext(
			//	ctx.WithProject(r.Context(), *project),
			//)

			next.ServeHTTP(w, r)
		})
	}
}
