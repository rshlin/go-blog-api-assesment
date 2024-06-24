package server

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestNewInMemoryAuthStore(t *testing.T) {

	tests := []struct {
		name string
		salt string
		data map[string]string
		want *InMemoryAuthStore
	}{
		{
			name: "empty data",
			salt: "salt",
			data: make(map[string]string),
			want: &InMemoryAuthStore{
				store: make(map[string]string),
				salt:  []byte("salt"),
			},
		},
		{
			name: "single data",
			salt: "salt",
			data: map[string]string{
				"user1": "password1",
			},
			want: &InMemoryAuthStore{
				store: map[string]string{
					"user1": hashAndHex("password1", "salt"),
				},
				salt: []byte("salt"),
			},
		},
		{
			name: "multiple data",
			salt: "salt",
			data: map[string]string{
				"user1": "password1",
				"user2": "password2",
			},
			want: &InMemoryAuthStore{
				store: map[string]string{
					"user1": hashAndHex("password1", "salt"),
					"user2": hashAndHex("password2", "salt"),
				},
				salt: []byte("salt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := InMemoryAuthStoreConfig{Salt: tt.salt, Data: tt.data}
			got := newInMemoryAuthStore(cfg)
			if !compareAuthStores(got, tt.want) {
				t.Errorf("newInMemoryAuthStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestValidate(t *testing.T) {
	var s = "Salt"
	var config = InMemoryAuthStoreConfig{Salt: s}
	var store = newInMemoryAuthStore(config)
	var hashedPass = sha256.Sum256(append([]byte("testPassword"), s...))
	store.store["testUser"] = hex.EncodeToString(hashedPass[:])

	tests := []struct {
		name     string
		username Username
		password Password
		want     bool
	}{
		{
			name:     "CorrectUsernameAndPassword",
			username: "testUser",
			password: []byte("testPassword"),
			want:     true,
		},
		{
			name:     "IncorrectUsername",
			username: "wrongUser",
			password: []byte("testPassword"),
			want:     false,
		},
		{
			name:     "IncorrectPassword",
			username: "testUser",
			password: []byte("wrongPassword"),
			want:     false,
		},
		{
			name:     "EmptyUsername",
			username: "",
			password: []byte("testPassword"),
			want:     false,
		},
		{
			name:     "EmptyPassword",
			username: "testUser",
			password: []byte(""),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := store.Validate(tt.username, tt.password); got != tt.want {
				t.Errorf("Validate() = %v, want = %v", got, tt.want)
			}
		})
	}
}

func hashAndHex(password, salt string) string {
	s := []byte(salt)
	hash := sha256.Sum256(append([]byte(password), s...))
	return hex.EncodeToString(hash[:])
}

func compareAuthStores(got, want *InMemoryAuthStore) bool {
	if len(got.store) != len(want.store) {
		return false
	}
	for k, gv := range got.store {
		wv, ok := want.store[k]
		if !ok || gv != wv {
			return false
		}
	}
	if len(got.salt) != len(want.salt) {
		return false
	}
	for i, v := range got.salt {
		if v != want.salt[i] {
			return false
		}
	}
	return true
}
