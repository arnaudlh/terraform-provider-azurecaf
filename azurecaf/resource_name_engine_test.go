package azurecaf

import (
	"testing"

	models "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/utils"
)

func init() {
	// Initialize test resource definitions
	models.ResourceDefinitions["azurerm_storage_account"] = models.ResourceStructure{
		ResourceTypeName: "azurerm_storage_account",
		CafPrefix:        "st",
		MinLength:        3,
		MaxLength:        24,
		RegEx:            "^[a-z0-9]{3,24}$",
		ValidationRegExp: "^[a-z0-9]{3,24}$",
		LowerCase:        true,
	}
}

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
			got := utils.RandSeq(tt.length, tt.seed)
			if len(got) != tt.want {
				t.Errorf("randSeq() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestConcatenateParameters(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		params    [][]string
		want      string
	}{
		{
			name:      "join with hyphen",
			separator: "-",
			params:    [][]string{{"prefix"}, {"name"}, {"suffix"}},
			want:      "prefix-name-suffix",
		},
		{
			name:      "empty separator",
			separator: "",
			params:    [][]string{{"prefix"}, {"name"}, {"suffix"}},
			want:      "prefixnamesuffix",
		},
		{
			name:      "single part",
			separator: "-",
			params:    [][]string{{"name"}},
			want:      "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := concatenateParameters(tt.separator, tt.params...); got != tt.want {
				t.Errorf("concatenateParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}
