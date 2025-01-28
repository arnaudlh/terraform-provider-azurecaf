//go:build unit

package azurecaf

import (
	"testing"
)

func TestRandSeq(t *testing.T) {
	tests := []struct {
		name   string
		n      int
		seed   int64
		expect int
	}{
		{
			name:   "zero length",
			n:      0,
			seed:   123,
			expect: 0,
		},
		{
			name:   "positive length",
			n:      5,
			seed:   123,
			expect: 5,
		},
		{
			name:   "large length",
			n:      20,
			seed:   123,
			expect: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := randSeq(tt.n, tt.seed)
			if len(result) != tt.expect {
				t.Errorf("randSeq() length = %v, want %v", len(result), tt.expect)
			}
			// Verify it only contains allowed characters
			for _, c := range result {
				if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
					t.Errorf("randSeq() contains invalid character: %c", c)
				}
			}
		})
	}
}
