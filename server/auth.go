package server

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	"net/http"
)

const AuthPrincipalContextKey = "auth-principal"

type AuthConfig struct {
	Type string
}

type Authenticator interface {
	Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error
	GetPrincipal(ctx context.Context, r *http.Request) interface{}
}

func NewAuthenticator(cfg *AuthConfig, store AuthStore) (Authenticator, error) {
	switch cfg.Type {
	case "basic":
		return newBasicAuthenticatorImpl(store), nil
	default:
		return nil, fmt.Errorf("unknown authentication type: %s", cfg.Type)
	}
}

func authMiddleware(a Authenticator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return a.Authenticate(ctx, input)
	}
}

func principalMiddleware(a Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal := a.GetPrincipal(r.Context(), r)

		ctx := context.WithValue(r.Context(), AuthPrincipalContextKey, principal)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
