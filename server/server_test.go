package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rshlin/go-blog-api-assesment/api"
	blogerr "github.com/rshlin/go-blog-api-assesment/blog/error"
	"github.com/rshlin/go-blog-api-assesment/blog/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockBlogService struct {
	mock.Mock
}

func (m *MockBlogService) FindById(ctx context.Context, id int) (*model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockBlogService) Create(ctx context.Context, post model.Post) (*model.Post, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockBlogService) Update(ctx context.Context, author model.Author, post model.Post) (*model.Post, error) {
	args := m.Called(ctx, author, post)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockBlogService) Delete(ctx context.Context, author model.Author, id int) error {
	args := m.Called(ctx, author, id)
	return args.Error(0)
}

func (m *MockBlogService) FindAll(ctx context.Context, page, size int) (*model.PaginatedPosts, error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).(*model.PaginatedPosts), args.Error(1)
}

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name                 string
		page                 int
		prepare              func(mock *MockBlogService)
		expectedHTTPStatus   int
		expectedResponseBody string
	}{
		{
			name: "Success",
			page: 1,
			prepare: func(m *MockBlogService) {
				m.On("FindAll", mock.Anything, 1, mock.Anything).Return(&model.PaginatedPosts{Data: make([]model.Post, 0)}, nil)
			},
			expectedHTTPStatus:   http.StatusOK,
			expectedResponseBody: `{"data":[],"pageMetadata":{"page":0,"size":0,"total":0}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockBlogService)
			tt.prepare(m)
			s := Server{
				blogService: m,
				config: &Config{
					ResponseTimeoutMs: 500,
				},
			}
			request := httptest.NewRequest("GET", "/posts", nil)
			responseRecorder := httptest.NewRecorder()

			s.GetPosts(responseRecorder, request, api.GetPostsParams{Page: tt.page, Size: 5})

			response := responseRecorder.Result()
			assert.Equal(t, tt.expectedHTTPStatus, response.StatusCode)
			assert.JSONEq(t, tt.expectedResponseBody, responseToString(response))
		})
	}
}

func TestServer_CreatePost(t *testing.T) {
	mockService := new(MockBlogService)
	ctx := context.Background()
	handler := &Server{
		blogService: mockService,
		config:      &Config{ResponseTimeoutMs: 1000},
	}

	tests := []struct {
		name          string
		mockRequest   api.CreatePostJSONRequestBody
		setup         func(m *MockBlogService)
		checkResponse func(recorder *httptest.ResponseRecorder)
		expectedErr   error
	}{
		{
			name: "valid case",
			mockRequest: api.CreatePostJSONRequestBody{
				Content: "Test content",
				Title:   "Test title",
				Author:  "Test author",
			},
			setup: func(m *MockBlogService) {
				m.On("Create", mock.Anything, mock.Anything).Return(&model.Post{}, nil).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "invalid case - service error",
			mockRequest: api.CreatePostJSONRequestBody{
				Content: "Test content",
				Title:   "Test title",
				Author:  "Test author",
			},
			setup: func(m *MockBlogService) {
				m.On("Create", mock.Anything, mock.Anything).Return(&model.Post{}, errors.New("unknown error")).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "invalid case - not found error",
			mockRequest: api.CreatePostJSONRequestBody{
				Content: "Test content",
				Title:   "Test title",
				Author:  "Test author",
			},
			setup: func(m *MockBlogService) {
				m.On("Create", mock.Anything, mock.Anything).Return(&model.Post{}, &blogerr.Error{ErrorType: blogerr.NotFound}).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "invalid case - forbidden error",
			mockRequest: api.CreatePostJSONRequestBody{
				Content: "Test content",
				Title:   "Test title",
				Author:  "Test author",
			},
			setup: func(m *MockBlogService) {
				m.On("Create", mock.Anything, mock.Anything).Return(&model.Post{}, &blogerr.Error{ErrorType: blogerr.Forbidden}).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer mockService.AssertExpectations(t)
			test.setup(mockService)

			recorder := httptest.NewRecorder()
			body, _ := json.Marshal(test.mockRequest)
			request, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(body))
			ctx = context.WithValue(ctx, AuthPrincipalContextKey, "testUser")

			handler.CreatePost(recorder, request.WithContext(ctx))

			test.checkResponse(recorder)
		})
	}
}

func TestDeletePostMethod(t *testing.T) {
	tests := []struct {
		name                 string
		prepare              func(*MockBlogService)
		id                   int
		expectedResponseBody string
		expectedStatus       int
	}{
		{
			name: "Successful deletion",
			prepare: func(service *MockBlogService) {
				service.On("Delete", mock.Anything, mock.Anything, 1).Return(nil)
			},
			id:                   1,
			expectedResponseBody: "",
			expectedStatus:       http.StatusNoContent,
		},
		{
			name: "Non-existent post",
			prepare: func(service *MockBlogService) {
				service.On("Delete", mock.Anything, mock.Anything, 2).Return(blogerr.NewError(blogerr.NotFound, "not found", nil))
			},
			id:                   2,
			expectedResponseBody: `{"msg":"ErrorType: Not Found, message: not found"}`,
			expectedStatus:       http.StatusNotFound,
		},
		{
			name: "Forbidden deletion",
			prepare: func(service *MockBlogService) {
				service.On("Delete", mock.Anything, mock.Anything, 3).Return(blogerr.NewError(blogerr.Forbidden, "forbidden", nil))
			},
			id:                   3,
			expectedResponseBody: `{"msg":"ErrorType: Forbidden, message: forbidden"}`,
			expectedStatus:       http.StatusForbidden,
		},
		{
			name: "Unknown server error",
			prepare: func(service *MockBlogService) {
				service.On("Delete", mock.Anything, mock.Anything, 4).Return(errors.New("unknown error"))
			},
			id:                   4,
			expectedResponseBody: `{"msg":"unknown error"}`,
			expectedStatus:       http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := new(MockBlogService)
			tc.prepare(service)

			s := NewServer(service, &Config{ResponseTimeoutMs: 1000})

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), AuthPrincipalContextKey, "user_id"))
			rr := httptest.NewRecorder()

			s.DeletePost(rr, req, tc.id)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if len(tc.expectedResponseBody) != 0 {
				assert.JSONEq(t, tc.expectedResponseBody, strings.TrimSuffix(rr.Body.String(), "\n"))
			}
			service.AssertExpectations(t)
		})
	}
}

func TestServer_GetPostById(t *testing.T) {
	testCases := []struct {
		name                string
		requestId           int
		blogServiceFunction func(mockSvc *MockBlogService)
		expectedRespStatus  int
	}{
		{
			name:      "Valid Post",
			requestId: 1,
			blogServiceFunction: func(mockSvc *MockBlogService) {
				mockSvc.On("FindById", mock.Anything, mock.Anything).
					Return(&model.Post{Id: 1, Title: "title", Author: "author", Content: "content"}, nil).Once()
			},
			expectedRespStatus: http.StatusOK,
		},
		{
			name:      "Post Not Found",
			requestId: 1,
			blogServiceFunction: func(mockSvc *MockBlogService) {
				mockSvc.On("FindById", mock.Anything, mock.Anything).Return(&model.Post{}, blogerr.NewError(blogerr.NotFound, "not found", nil)).Once()
			},
			expectedRespStatus: http.StatusNotFound,
		},
		{
			name:      "Database Error",
			requestId: 1,
			blogServiceFunction: func(mockSvc *MockBlogService) {
				mockSvc.On("FindById", mock.Anything, mock.Anything).Return(&model.Post{}, errors.New("database error")).Once()
			},
			expectedRespStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(MockBlogService)
			tc.blogServiceFunction(mockSvc)

			s := &Server{
				blogService: mockSvc,
				config:      &Config{ResponseTimeoutMs: 5000},
			}

			request, _ := http.NewRequest("GET", fmt.Sprintf("/posts/%d", tc.requestId), bytes.NewBuffer(nil))
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				s.GetPostById(w, r, tc.requestId)
			})

			handler.ServeHTTP(response, request)

			mockSvc.AssertExpectations(t)
			assert.Equal(t, tc.expectedRespStatus, response.Code)

			if tc.expectedRespStatus == http.StatusOK {
				var resp api.Post
				_ = json.Unmarshal(response.Body.Bytes(), &resp)
				assert.Equal(t, tc.requestId, resp.Id)
			}
		})
	}
}

func TestServer_UpdatePost(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		prepare    func(m *MockBlogService)
		reqBody    api.Post
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Successful update",
			id:   1,
			prepare: func(m *MockBlogService) {
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&model.Post{}, nil)
			},
			reqBody:    api.Post{Title: "New Post", Content: "New Content"},
			wantStatus: http.StatusOK,
		},
		{
			name: "Error mapping request body to struct",
			id:   1,
			prepare: func(m *MockBlogService) {
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&model.Post{}, errors.New("error mapping struct"))
			},
			reqBody:    api.Post{Title: "New Post", Content: "New Content"},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Error on update",
			id:   1,
			prepare: func(m *MockBlogService) {
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&model.Post{}, errors.New("update error"))
			},
			reqBody:    api.Post{Title: "New Post", Content: "New Content"},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Not found error on update",
			id:   1,
			prepare: func(m *MockBlogService) {
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&model.Post{}, blogerr.NewError(blogerr.NotFound, "not found", nil))
			},
			reqBody:    api.Post{Title: "New Post", Content: "New Content"},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Forbidden error on update",
			id:   1,
			prepare: func(m *MockBlogService) {
				m.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&model.Post{}, blogerr.NewError(blogerr.Forbidden, "forbidden", nil))
			},
			reqBody:    api.Post{Title: "New Post", Content: "New Content"},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.reqBody)

			req, err := http.NewRequest("POST", "/update", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": "1"})

			config := &Config{ResponseTimeoutMs: 1000}
			m := &MockBlogService{}
			tt.prepare(m)

			ctx := context.WithValue(context.Background(), AuthPrincipalContextKey, "test_user_id")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			s := NewServer(m, config)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				s.UpdatePost(w, r, tt.id)
			})
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

// Converts http.Response to string.
func responseToString(response *http.Response) string {
	var responseBody bytes.Buffer
	_, _ = responseBody.ReadFrom(response.Body)
	_ = response.Body.Close()
	return responseBody.String()
}
