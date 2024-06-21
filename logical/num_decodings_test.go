package logical

import (
	"testing"
)

func TestNumDecodings(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{
			name: "Test 1",
			s:    "12",
			want: 2,
		},
		{
			name: "Test 2",
			s:    "226",
			want: 3,
		},
		{
			name: "Test 3",
			s:    "0",
			want: 0,
		},
		{
			name: "Test 4",
			s:    "06",
			want: 0,
		},
		{
			name: "Test 5",
			s:    "2101",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumDecodings(tt.s); got != tt.want {
				t.Errorf("numDecodings() = %v, want %v", got, tt.want)
			}
		})
	}
}
