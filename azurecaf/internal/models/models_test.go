package models

import (
	"testing"
	"regexp"
	"fmt"
)

func init() {
	// Initialize test resource definitions
	ResourceDefinitions["azurerm_storage_account"] = ResourceStructure{
		ResourceTypeName:  "azurerm_storage_account",
		CafPrefix:        "st",
		MinLength:        3,
		MaxLength:        24,
		RegEx:           "^[a-z0-9]{3,24}$",
		ValidationRegExp: "^[a-z0-9]{3,24}$",
		LowerCase:       true,
	}
}

func TestValidateResourceType(t *testing.T) {
	tests := []struct {
		name        string
		resourceType string
		wantErr     bool
	}{
		{
			name:        "valid resource type",
			resourceType: "azurerm_storage_account",
			wantErr:     false,
		},
		{
			name:        "invalid resource type",
			resourceType: "invalid_resource",
			wantErr:     true,
		},
		{
			name:        "empty resource type",
			resourceType: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateResourceType(tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateResourceType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetResourceStructure(t *testing.T) {
	tests := []struct {
		name        string
		resourceType string
		wantErr     bool
	}{
		{
			name:        "existing resource type",
			resourceType: "azurerm_storage_account",
			wantErr:     false,
		},
		{
			name:        "non-existent resource type",
			resourceType: "nonexistent_type",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getResourceStructure(tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResourceStructure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		minLength   int
		maxLength   int
		wantErr     bool
	}{
		{
			name:      "valid length",
			input:     "test123",
			minLength: 3,
			maxLength: 10,
			wantErr:   false,
		},
		{
			name:      "too short",
			input:     "ab",
			minLength: 3,
			maxLength: 10,
			wantErr:   true,
		},
		{
			name:      "too long",
			input:     "verylongstring",
			minLength: 3,
			maxLength: 10,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLength(tt.input, tt.minLength, tt.maxLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRegex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		regex       string
		wantErr     bool
	}{
		{
			name:    "valid storage account name",
			input:   "teststorage123",
			regex:   "^[a-z0-9]{3,24}$",
			wantErr: false,
		},
		{
			name:    "invalid characters",
			input:   "test-storage",
			regex:   "^[a-z0-9]{3,24}$",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			regex:   "^[a-z0-9]{3,24}$",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegex(tt.input, tt.regex)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRegex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
