package server

import (
	"context"
	"encoding/base64"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAuthStore struct {
	DB map[string][]byte
}

func (m *MockAuthStore) Validate(username string, password []byte) bool {
	pass, ok := m.DB[username]
	if !ok {
		return false
	}
	if string(pass) != string(password) {
		return false
	}
	return true
}

var testStore = &MockAuthStore{
	DB: map[string][]byte{
		"admin": []byte("password"),
	},
}

func Test_BasicAuthenticatorImpl_Authenticate(t *testing.T) {
	authImpl := BasicAuthenticatorImpl{
		store: testStore,
	}

	testCases := []struct {
		name      string
		auth      string
		input     *openapi3filter.AuthenticationInput
		expectErr bool
	}{
		{
			name: "ValidCredentials",
			auth: "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:password")),
			input: &openapi3filter.AuthenticationInput{
				RequestValidationInput: &openapi3filter.RequestValidationInput{
					Request: httptest.NewRequest("GET", "http://example.com", nil),
				},
			},
			expectErr: false,
		},
		{
			name: "InvalidCredentials",
			auth: "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:badpassword")),
			input: &openapi3filter.AuthenticationInput{
				RequestValidationInput: &openapi3filter.RequestValidationInput{
					Request: httptest.NewRequest("GET", "http://example.com", nil),
				},
			},
			expectErr: true,
		},
		{
			name: "NoAuthorizationHeader",
			auth: "",
			input: &openapi3filter.AuthenticationInput{
				RequestValidationInput: &openapi3filter.RequestValidationInput{
					Request: httptest.NewRequest("GET", "http://example.com", nil),
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.RequestValidationInput.Request.Header.Set("Authorization", tc.auth)
			err := authImpl.Authenticate(context.Background(), tc.input)
			if (err != nil) != tc.expectErr {
				t.Errorf("%v: got %v, expected error? %v", tc.name, err, tc.expectErr)
			}
		})
	}
}

func TestGetPrincipal(t *testing.T) {
	store := testStore
	auth := newBasicAuthenticatorImpl(store)

	tests := []struct {
		name     string
		request  *http.Request
		expected string
	}{
		{
			name:     "Empty Authorization Header",
			request:  &http.Request{Header: make(http.Header)},
			expected: "",
		},
		{
			name: "Non-Empty Authorization Header with Valid Credentials",
			request: func() *http.Request {
				req := &http.Request{Header: make(http.Header)}
				req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testUser:testPass")))
				return req
			}(),
			expected: "testUser",
		},
		{
			name: "Non-Empty Authorization Header with Invalid Credentials",
			request: func() *http.Request {
				req := &http.Request{Header: make(http.Header)}
				req.Header.Set("Authorization", "Invalid Credentials")
				return req
			}(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			actual := auth.GetPrincipal(ctx, tt.request)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
