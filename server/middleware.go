package server

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rshlin/go-blog-api-assesment/api"
	"log"
	"net/http"
)

func CreateMiddleware(authenticator Authenticator, cfg *Config) []mux.MiddlewareFunc {
	swagger, err := openapi3.NewLoader().LoadFromData(api.Schema)
	swagger.Servers = nil
	if err != nil {
		log.Fatal("error loading swagger spec", err)
	}

	options := middleware.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody:          cfg.ValidationExcludeRequestBody,
			ExcludeRequestQueryParams:   cfg.ValidationExcludeRequestQueryParams,
			ExcludeResponseBody:         cfg.ValidationExcludeResponseBody,
			ExcludeReadOnlyValidations:  cfg.ValidationExcludeReadOnlyValidations,
			ExcludeWriteOnlyValidations: cfg.ValidationExcludeWriteOnlyValidations,
			IncludeResponseStatus:       cfg.ValidationIncludeResponseStatus,
			AuthenticationFunc:          authMiddleware(authenticator),
			SkipSettingDefaults:         false,
		},
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			// todo logger

			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(api.Error{Msg: message})
		},
	}
	oapiValidator := middleware.OapiRequestValidatorWithOptions(swagger, &options)

	principalMw := func(next http.Handler) http.Handler {
		return principalMiddleware(authenticator, next)
	}

	return []mux.MiddlewareFunc{
		oapiValidator,
		principalMw,
	}
}
