package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAuthenticator(t *testing.T) {
	store := &MockAuthStore{}

	// Test basic authentication type
	cfg := &AuthConfig{Type: "basic"}
	authenticator, err := NewAuthenticator(cfg, store)
	assert.NoError(t, err)
	assert.IsType(t, &BasicAuthenticatorImpl{}, authenticator)

	// Test unknown authentication type
	cfg = &AuthConfig{Type: "unknown"}
	authenticator, err = NewAuthenticator(cfg, store)
	assert.Error(t, err)
	assert.Nil(t, authenticator)
}
