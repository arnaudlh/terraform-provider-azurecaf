package azurecaf

import (
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
		seed   int64
		want   int
	}{
		{
			name:   "generate 5 char string",
			length: 5,
			seed:   123,
			want:   5,
		},
		{
			name:   "generate 10 char string",
			length: 10,
			seed:   456,
			want:   10,
		},
		{
			name:   "zero length",
			length: 0,
			seed:   789,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randSeq(tt.length, &tt.seed)
			if len(got) != tt.want {
				t.Errorf("randSeq() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestCleanInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "remove special characters",
			input: "Test@123!",
			want:  "Test123",
		},
		{
			name:  "remove spaces",
			input: "Test Input",
			want:  "TestInput",
		},
		{
			name:  "already clean",
			input: "Test123",
			want:  "Test123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanInput(tt.input); got != tt.want {
				t.Errorf("cleanInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildResourceName(t *testing.T) {
	tests := []struct {
		name       string
		parts      []string
		separator  string
		want       string
	}{
		{
			name:      "join with hyphen",
			parts:     []string{"prefix", "name", "suffix"},
			separator: "-",
			want:      "prefix-name-suffix",
		},
		{
			name:      "empty separator",
			parts:     []string{"prefix", "name", "suffix"},
			separator: "",
			want:      "prefixnamesuffix",
		},
		{
			name:      "single part",
			parts:     []string{"name"},
			separator: "-",
			want:      "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildResourceName(tt.parts, tt.separator); got != tt.want {
				t.Errorf("buildResourceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
