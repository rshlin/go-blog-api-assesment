package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rshlin/go-blog-api-assesment/utils"
	"net/http"
	"time"

	"github.com/rshlin/go-blog-api-assesment/api"
	blogerr "github.com/rshlin/go-blog-api-assesment/blog/error"
	blogmodel "github.com/rshlin/go-blog-api-assesment/blog/model"
	blogsvc "github.com/rshlin/go-blog-api-assesment/blog/service"
)

type Server struct {
	blogService blogsvc.BlogService
	config      *Config
}

func NewServer(blogService blogsvc.BlogService, config *Config) *Server {
	return &Server{blogService: blogService, config: config}
}

func (s *Server) handleErrors(w http.ResponseWriter, err error) {
	var e *blogerr.Error
	switch {
	case errors.As(err, &e):
		switch e.ErrorType {
		case blogerr.NotFound:
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(api.Error{Msg: e.Error()})
			return
		case blogerr.Forbidden:
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(api.Error{Msg: e.Error()})
			return
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(api.Error{Msg: "Request timed out"})
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(api.Error{Msg: err.Error()})
	}
}

func (s *Server) GetPosts(w http.ResponseWriter, r *http.Request, params api.GetPostsParams) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.config.ResponseTimeoutMs)*time.Millisecond)
	defer cancel()

	posts, err := s.blogService.FindAll(ctx, params.Page, params.Size)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	var postsResponse api.PaginatedPosts
	if err = utils.MapToStruct(posts, &postsResponse); err != nil {
		s.handleErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postsResponse)
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	var createPost api.CreatePostJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&createPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var post blogmodel.Post

	err = utils.MapToStruct(createPost, &post)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.config.ResponseTimeoutMs)*time.Millisecond)
	defer cancel()

	principal := ctx.Value(AuthPrincipalContextKey)
	principalID := principal.(string)
	post.Author = principalID

	created, err := s.blogService.Create(ctx, post)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	var createdPostResponse api.Post
	err = utils.MapToStruct(created, &createdPostResponse)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdPostResponse)
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request, id int) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.config.ResponseTimeoutMs)*time.Millisecond)
	defer cancel()

	principal := ctx.Value(AuthPrincipalContextKey)
	principalID := principal.(string)

	err := s.blogService.Delete(ctx, principalID, id)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetPostById(w http.ResponseWriter, r *http.Request, id int) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.config.ResponseTimeoutMs)*time.Millisecond)
	defer cancel()

	post, err := s.blogService.FindById(ctx, id)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	var postResponse api.Post
	err = utils.MapToStruct(post, &postResponse)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(postResponse)
}

func (s *Server) UpdatePost(w http.ResponseWriter, r *http.Request, id int) {
	var updatePost api.UpdatePostJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&updatePost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(s.config.ResponseTimeoutMs)*time.Millisecond)
	defer cancel()

	var post blogmodel.Post
	err = utils.MapToStruct(updatePost, &post)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	principal := ctx.Value(AuthPrincipalContextKey)
	principalID := principal.(string)
	post.Author = principalID
	post.Id = id

	updated, err := s.blogService.Update(ctx, principalID, post)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	var updatedPostResponse api.Post
	err = utils.MapToStruct(updated, &updatedPostResponse)
	if err != nil {
		s.handleErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(updatedPostResponse)
}
