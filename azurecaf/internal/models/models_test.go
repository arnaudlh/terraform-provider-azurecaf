package models

import (
	"testing"
)

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
			_, err := ValidateResourceType(tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResourceType() error = %v, wantErr %v", err, tt.wantErr)
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
			_, err := GetResourceStructure(tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResourceStructure() error = %v, wantErr %v", err, tt.wantErr)
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
			err := ValidateLength(tt.input, tt.minLength, tt.maxLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLength() error = %v, wantErr %v", err, tt.wantErr)
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
			err := ValidateRegex(tt.input, tt.regex)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
