//go:build unit

package main

import (
	"testing"
)

func TestValidateResourceDefinition(t *testing.T) {
	tests := []struct {
		name       string
		resources  []string
		wantErr    bool
		errMessage string
	}{
		{
			name:      "valid resources",
			resources: []string{"azurerm_resource_group", "azurerm_storage_account"},
			wantErr:   false,
		},
		{
			name:       "invalid resource",
			resources:  []string{"invalid_resource"},
			wantErr:    true,
			errMessage: "resource type invalid_resource not found in the resource definition file",
		},
		{
			name:      "empty resource list",
			resources: []string{},
			wantErr:   false,
		},
		{
			name:       "mixed valid and invalid resources",
			resources:  []string{"azurerm_resource_group", "invalid_resource"},
			wantErr:    true,
			errMessage: "resource type invalid_resource not found in the resource definition file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResourceDefinition(tt.resources)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResourceDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMessage {
				t.Errorf("ValidateResourceDefinition() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}

func TestGetResourceDefinition(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "get resource definition",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResourceDefinition()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResourceDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("GetResourceDefinition() returned nil map")
			}
		})
	}
}

func TestGetResourceMap(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "get resource map",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResourceMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResourceMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("GetResourceMap() returned nil map")
			}
		})
	}
}
