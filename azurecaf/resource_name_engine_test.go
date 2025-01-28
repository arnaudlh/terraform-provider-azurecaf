//go:build unit

package azurecaf

import (
	"testing"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
)

func TestCleanSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "slice with empty strings",
			input:    []string{"", "test", ""},
			expected: []string{"test"},
		},
		{
			name:     "slice with spaces",
			input:    []string{" ", "test", " space "},
			expected: []string{"test", "space"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := &models.ResourceStructure{
				MinLength: 1,
				MaxLength: 63,
			}
			result := cleanSlice(tt.input, resource)
			if len(result) != len(tt.expected) {
				t.Errorf("cleanSlice() got %v, want %v", result, tt.expected)
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("cleanSlice() got %v, want %v", result, tt.expected)
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
		{"string with spaces", " test ", "test"},
		{"string with special chars", "test@#$%", "test"},
		{"mixed case", "TestString", "teststring"},
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

func TestConcatenateParameters(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		params    []string
		expected  string
	}{
		{"empty params", "", []string{}, ""},
		{"single param", "test", []string{"-"}, "test"},
		{"multiple params", "a", []string{"-", "b", "-", "c"}, "a-b-c"},
		{"custom separator", "x", []string{"_", "y"}, "x_y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := concatenateParameters(tt.input, tt.params)
			if result != tt.expected {
				t.Errorf("concatenateParameters() = %v, want %v", result, tt.expected)
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
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := getResource(tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
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
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple string", "test", "test"},
		{"mixed case", "TestString", "teststring"},
		{"with spaces", "test string", "test-string"},
		{"with special chars", "test@#$%string", "test-string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSlug(tt.input)
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
			err, _ := validateResourceType(tt.resourceType, tt.resourceTypes)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateResourceType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
