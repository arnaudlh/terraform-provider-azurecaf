package utils

import (
	"testing"
)

func TestRandSeq(t *testing.T) {
	tests := []struct {
		name   string
		n      int
		seed   int64
		expect int
		chars  bool
	}{
		{
			name:   "zero length",
			n:      0,
			seed:   1,
			expect: 0,
			chars:  true,
		},
		{
			name:   "positive length",
			n:      5,
			seed:   1,
			expect: 5,
			chars:  true,
		},
		{
			name:   "different seed",
			n:      5,
			seed:   2,
			expect: 5,
			chars:  true,
		},
		{
			name:   "long string",
			n:      20,
			seed:   3,
			expect: 20,
			chars:  true,
		},
		{
			name:   "negative length",
			n:      -1,
			seed:   1,
			expect: 0,
			chars:  true,
		},
		{
			name:   "max length",
			n:      100,
			seed:   1,
			expect: 100,
			chars:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandSeq(tt.n, tt.seed)
			if len(result) != tt.expect {
				t.Errorf("RandSeq() length = %v, want %v", len(result), tt.expect)
			}
			if tt.chars {
				for _, c := range result {
					if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
						t.Errorf("RandSeq() contains invalid character: %c", c)
					}
				}
			}
		})
	}
}

func TestRandSeqDeterministic(t *testing.T) {
	n := 10
	seed := int64(42)
	
	result1 := RandSeq(n, seed)
	result2 := RandSeq(n, seed)
	
	if result1 != result2 {
		t.Errorf("RandSeq() not deterministic: got %v and %v for same seed", result1, result2)
	}
}

func TestRandSeqUniqueness(t *testing.T) {
	n := 10
	results := make(map[string]bool)
	
	for seed := int64(0); seed < 10; seed++ {
		result := RandSeq(n, seed)
		if results[result] {
			t.Errorf("RandSeq() generated duplicate string with different seeds: %s", result)
		}
		results[result] = true
	}
}
