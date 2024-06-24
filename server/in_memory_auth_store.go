package server

import (
	"crypto/sha256"
	"encoding/hex"
)

type InMemoryAuthStoreConfig struct {
	Salt string
	Data map[string]string
}

type salt = []byte

type hashedPassword = string

type InMemoryAuthStore struct {
	store map[Username]hashedPassword
	salt  salt
}

func newInMemoryAuthStore(cfg InMemoryAuthStoreConfig) *InMemoryAuthStore {
	s := []byte(cfg.Salt)
	store := make(map[string]string)
	for username, password := range cfg.Data {
		hash := sha256.Sum256(append([]byte(password), s...))
		store[username] = hex.EncodeToString(hash[:])
	}
	return &InMemoryAuthStore{
		store: store,
		salt:  s,
	}
}

func (s *InMemoryAuthStore) Validate(username Username, password Password) bool {
	hash := sha256.Sum256(append(password, s.salt...))
	expectedHash, ok := s.store[username]
	if !ok {
		return false
	}
	return hex.EncodeToString(hash[:]) == expectedHash
}
