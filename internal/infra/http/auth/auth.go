package auth

import (
	"companies/internal/config"
	"companies/internal/domain/interfaces"
	"github.com/golang-jwt/jwt/v4"
)

type JwtValidator struct {
	authConfig config.Auth
}

func NewJwtValidator(authConfig config.Auth) *JwtValidator {
	return &JwtValidator{authConfig: authConfig}
}

func (o *JwtValidator) ByJwt(tokenStr string) error {
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(o.authConfig.JwtKey), nil
		},
	)
	if err != nil {
		return err
	}
	if !token.Valid {
		return interfaces.ErrTokenNotValid
	}
	return nil
}
