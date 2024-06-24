package server

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/getkin/kin-openapi/openapi3filter"
	"net/http"
	"strings"
)

var InvalidHeaderError = errors.New("invalid Authorization header")
var InvalidCredentials = errors.New("invalid username or password")
var InvalidPayloadError = errors.New("invalid Authorization payload")

type BasicAuthenticatorImpl struct {
	store AuthStore
}

func newBasicAuthenticatorImpl(store AuthStore) *BasicAuthenticatorImpl {
	return &BasicAuthenticatorImpl{
		store: store,
	}
}

func parseHeader(auth string) (string, string, error) {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", InvalidHeaderError
	}
	payload, _ := base64.StdEncoding.DecodeString(parts[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", InvalidPayloadError
	}
	return pair[0], pair[1], nil
}

func (a *BasicAuthenticatorImpl) Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	r := input.RequestValidationInput.Request
	auth := r.Header.Get("Authorization")

	username, password, err := parseHeader(auth)
	if err != nil {
		return err
	}
	usernameLower := strings.ToLower(username)

	if !a.store.Validate(usernameLower, []byte(password)) {
		return InvalidCredentials
	}
	return nil
}

func (a *BasicAuthenticatorImpl) GetPrincipal(ctx context.Context, r *http.Request) interface{} {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	username, _, err := parseHeader(auth)
	if err != nil {
		return ""
	}
	return username
}
