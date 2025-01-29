//go:build unit

package azurecaf

import (
	"testing"

	"github.com/arnaudlh/terraform-provider-azurecaf/azurecaf/internal/models"
)

func TestCleanSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		regex    string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			regex:    "^[a-z0-9]+$",
			expected: []string{},
		},
		{
			name:     "slice with empty strings",
			input:    []string{"", "test123", ""},
			regex:    "^[a-z0-9]+$",
			expected: []string{"", "test123", ""},
		},
		{
			name:     "slice with invalid chars",
			input:    []string{"test@123", "test-123", "test_123"},
			regex:    "^[a-z0-9-]+$",
			expected: []string{"test123", "test-123", "test123"},
		},
		{
			name:     "slice with mixed content",
			input:    []string{"", "test123", "test@456", ""},
			regex:    "^[a-z0-9]+$",
			expected: []string{"", "test123", "test456", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := &models.ResourceStructure{
				MinLength:        1,
				MaxLength:        63,
				RegEx:           tt.regex,
				ValidationRegExp: tt.regex,
			}
			result := cleanSlice(tt.input, resource)
			if len(result) != len(tt.expected) {
				t.Errorf("cleanSlice() got %v, want %v", result, tt.expected)
				return
			}
			for i := range tt.expected {
				if result[i] != tt.expected[i] {
					t.Errorf("cleanSlice() index %d got %v, want %v", i, result[i], tt.expected[i])
					return
				}
			}
		})
	}
}

func TestCleanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"string with spaces", " test ", " test "},
		{"string with special chars", "test@#$%", "test@#$%"},
		{"mixed case", "TestString", "TestString"},
	}

	resource := &models.ResourceStructure{
		MinLength: 1,
		MaxLength: 63,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanString(tt.input, resource)
			if result != tt.expected {
				t.Errorf("cleanString() = %v, want %v", result, tt.expected)
			}
		})
	}
}



func TestGetResource(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		wantErr      bool
		wantMinLen   int
		wantMaxLen   int
	}{
		{
			name:         "valid resource",
			resourceType: "azurerm_resource_group",
			wantErr:     false,
			wantMinLen:  1,
			wantMaxLen:  90,
		},
		{
			name:         "invalid resource",
			resourceType: "invalid_resource",
			wantErr:     true,
		},
		{
			name:         "empty resource",
			resourceType: "",
			wantErr:     false,
			wantMinLen:  1,
			wantMaxLen:  250,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := getResource(tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resource != nil {
				if resource.MinLength != tt.wantMinLen {
					t.Errorf("getResource() MinLength = %v, want %v", resource.MinLength, tt.wantMinLen)
				}
				if resource.MaxLength != tt.wantMaxLen {
					t.Errorf("getResource() MaxLength = %v, want %v", resource.MaxLength, tt.wantMaxLen)
				}
			}
		})
	}
}

func TestGetSlug(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expected     string
	}{
		{
			name:         "empty resource type",
			resourceType: "",
			expected:     "",
		},
		{
			name:         "valid resource type",
			resourceType: "azurerm_resource_group",
			expected:     "rg",
		},
		{
			name:         "invalid resource type",
			resourceType: "invalid_resource",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSlug(tt.resourceType)
			if result != tt.expected {
				t.Errorf("getSlug() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrimResourceName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{"short name", "test", 10, "test"},
		{"exact length", "test", 4, "test"},
		{"long name", "teststring", 4, "test"},
		{"empty string", "", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimResourceName(tt.input, tt.maxLength)
			if result != tt.expected {
				t.Errorf("trimResourceName() = %v, want %v", result, tt.expected)
			}
			if len(result) > tt.maxLength {
				t.Errorf("trimResourceName() returned string longer than maxLength")
			}
		})
	}
}

func TestValidateResourceType(t *testing.T) {
	tests := []struct {
		name          string
		resourceType  string
		resourceTypes []string
		wantErr       bool
	}{
		{
			name:          "valid resource",
			resourceType:  "azurerm_resource_group",
			resourceTypes: []string{"azurerm_resource_group"},
			wantErr:      false,
		},
		{
			name:          "invalid resource",
			resourceType:  "invalid_resource",
			resourceTypes: []string{"azurerm_resource_group"},
			wantErr:      true,
		},
		{
			name:          "empty resource",
			resourceType:  "",
			resourceTypes: []string{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, _ := validateResourceType(tt.resourceType, tt.resourceTypes)
			if isValid == tt.wantErr {
				t.Errorf("validateResourceType() got = %v, want = %v", isValid, !tt.wantErr)
			}
		})
	}
}
