package server

import "fmt"

type AuthStoreConfig struct {
	Type     string
	InMemory *InMemoryAuthStoreConfig
}

type Username = string
type Password = []byte

type AuthStore interface {
	Validate(username Username, password Password) bool
}

func NewAuthStore(config *AuthStoreConfig) (AuthStore, error) {
	switch config.Type {
	case "in-memory":
		return newInMemoryAuthStore(*config.InMemory), nil
	default:
		return nil, fmt.Errorf("unknown auth store type: %s", config.Type)
	}
}
