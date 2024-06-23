// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
)

const (
	BasicAuthScopes = "BasicAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Msg *string `json:"msg,omitempty"`
}

// PaginatedPosts defines model for PaginatedPosts.
type PaginatedPosts struct {
	Data *[]Post `json:"data,omitempty"`
	Meta *struct {
		Page  *int `json:"page,omitempty"`
		Size  *int `json:"size,omitempty"`
		Total *int `json:"total,omitempty"`
	} `json:"meta,omitempty"`
}

// Post defines model for Post.
type Post struct {
	Author  *string `json:"author,omitempty"`
	Content *string `json:"content,omitempty"`
	Id      *int    `json:"id,omitempty"`
	Title   *string `json:"title,omitempty"`
}

// GetPostsParams defines parameters for GetPosts.
type GetPostsParams struct {
	// Page Page number
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of items per page
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

// CreatePostJSONRequestBody defines body for CreatePost for application/json ContentType.
type CreatePostJSONRequestBody = Post

// UpdatePostJSONRequestBody defines body for UpdatePost for application/json ContentType.
type UpdatePostJSONRequestBody = Post

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Retrieve a list of all blog posts
	// (GET /posts)
	GetPosts(w http.ResponseWriter, r *http.Request, params GetPostsParams)
	// Create a new blog post
	// (POST /posts)
	CreatePost(w http.ResponseWriter, r *http.Request)
	// Delete a blog post
	// (DELETE /posts/{id})
	DeletePost(w http.ResponseWriter, r *http.Request, id int)
	// Retrieve details of a specific blog post
	// (GET /posts/{id})
	GetPostById(w http.ResponseWriter, r *http.Request, id int)
	// Update an existing blog post
	// (PUT /posts/{id})
	UpdatePost(w http.ResponseWriter, r *http.Request, id int)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetPosts operation middleware
func (siw *ServerInterfaceWrapper) GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPostsParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "size" -------------

	err = runtime.BindQueryParameter("form", true, false, "size", r.URL.Query(), &params.Size)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "size", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPosts(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreatePost operation middleware
func (siw *ServerInterfaceWrapper) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreatePost(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeletePost operation middleware
func (siw *ServerInterfaceWrapper) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", mux.Vars(r)["id"], &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeletePost(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetPostById operation middleware
func (siw *ServerInterfaceWrapper) GetPostById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", mux.Vars(r)["id"], &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPostById(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UpdatePost operation middleware
func (siw *ServerInterfaceWrapper) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", mux.Vars(r)["id"], &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdatePost(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/posts", wrapper.GetPosts).Methods("GET")

	r.HandleFunc(options.BaseURL+"/posts", wrapper.CreatePost).Methods("POST")

	r.HandleFunc(options.BaseURL+"/posts/{id}", wrapper.DeletePost).Methods("DELETE")

	r.HandleFunc(options.BaseURL+"/posts/{id}", wrapper.GetPostById).Methods("GET")

	r.HandleFunc(options.BaseURL+"/posts/{id}", wrapper.UpdatePost).Methods("PUT")

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xY0W7bNhT9FeJuj0ItuRlQ6C1utyEvhbFsT0EeaOlaZiGRLElldQ39+3BJWbJkJYhR",
	"o0EyP0kl68PDy3OOLrODTFVaSZTOQroDm22w4v71d2OUoRdtlEbjBPrhyhb0cFuNkIJ1RsgCmibaj6jV",
	"F8wcNBEseSEkd5gvlQ3oQ6ScO05P4bDyA78aXEMKv8x6TrOW0IwwoF+FG8O39O8KA8gQWvMCD1gK6bBA",
	"Q//fiu+PzDjleDk1dby3yd0SwSMivHabUMVRwSLIlHQo3eScyB/hKFyJzyo/bRWz2gi3vaUSBjoLbkV2",
	"XbtNd9b0mxWNQgexcU5DQwhCrpVfLSwLi1IVbFlyt1amYtfLG4jgAY0VSkIKybv4XUwslUbJtYAU3vuh",
	"CDR3G09gpvdaKNBvnYrFnVDyJocU/kQXxEI/MbxCh8ZCereDHG1mhHZhqSUvkMm6WqEBogkpfK3RbCEC",
	"yf2evAKiVs9et0KKqq4gTaKJEx4v8NljM7VmXp1Mo2Et5NRqXlWHqx2tcB+BQauVtOEk5nFMjwMRcK1L",
	"kflazL5YYrED/MYrXXaaJjuJMPPAyxp7F931UoNr/8ISONAYfAxvfpTUlXRigr/pyRKgKoxB5pMg8xZk",
	"PgaZA22086R3YbI3XRJ3JkviuPEK6yv2pPmHSeJ/OTyv2zrL0Np1XbJOUSTFq5PLjPvY6yrsEw8WPGcG",
	"v9ZIOfR86iFFJxgP8CL47WxEb9E8oGHYL/yDTEeAEdi6qrjZQgp/oTMCH5BxVgrryC+8LNmKciI4vYlA",
	"t8k4tPpHg9yhT02yhi/EQuXbU13RgndFeL4NjvR/giBVUEFz5Ork5/F/zManb2N84jTOMn9A+aty0dXJ",
	"9X+M6GflWDgK8Z2KcAauY0ii+/5MdP9QZiXyHOVZmB6gvbJkalse/0E8aHbu7um71OdWCB/GmcR/+7Ty",
	"EKFDme1E3vjPK5bo8Di+PvnxNr6e7lXITDef9p0DtUJ94yDyNv2EwRxSZ2o8sY24CiyPVgzEL/b9v9v3",
	"KijkDCS9rKRybK1qeZ6ajiHfZtqErGD8MGmiJy9Ai+1N/sK5Er+FRubVXwsu7v2Zt5gcHRel9dcYZjVm",
	"Yi2yoWt1PeHaf3TOX6oZOP/FKewmZ8G3A9fup1r3Hlyi9jPevOe5SsUvs6NRGv3wxiZdUwfUS3d26c4u",
	"+f6y3VkwOOOS4TdhnZDF4EbokQg6xHltyvZP5OlsVqqMlxtlXfoh/hBDc9/8FwAA//8CyBLtUBkAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
