package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAuthStore(t *testing.T) {
	tests := []struct {
		name    string
		config  *AuthStoreConfig
		wantErr bool
	}{
		{
			name:    "valid in-memory type",
			config:  &AuthStoreConfig{Type: "in-memory", InMemory: &InMemoryAuthStoreConfig{}},
			wantErr: false,
		},
		{
			name:    "unknown type",
			config:  &AuthStoreConfig{Type: "unknown", InMemory: nil},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAuthStore(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, fmt.Sprintf("unknown auth store type: %s", tt.config.Type), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
