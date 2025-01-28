//go:build unit

package azurecaf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataNameRead(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		wantErr  bool
		errCount int
	}{
		{
			name: "valid resource group name with special characters",
			data: map[string]interface{}{
				"name":          "rg123abc",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"dev"},
				"suffixes":      []interface{}{"prod"},
				"random_length": 0,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":  false,
			},
			wantErr: false,
		},
		{
			name: "resource group name with all allowed characters",
			data: map[string]interface{}{
				"name":          "rg123abc",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"dev"},
				"suffixes":      []interface{}{"prod"},
				"random_length": 0,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":  false,
			},
			wantErr: false,
		},
		{
			name: "invalid resource type",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "invalid_type",
			},
			wantErr:  true,
			errCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataName().Schema, tt.data)
			err := dataNameRead(context.Background(), d, nil)

			if tt.wantErr && err == nil {
				t.Error("dataNameRead() expected error, got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("dataNameRead() unexpected error: %v", err)
			}
		})
	}
}

func TestGetNameReadResult(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid resource group name with maximum length",
			data: map[string]interface{}{
				"name":          "rg123abc",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"dev"},
				"suffixes":      []interface{}{"prod"},
				"random_length": 0,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":  false,
			},
			wantErr: false,
		},
		{
			name: "resource group name with invalid character",
			data: map[string]interface{}{
				"name":          "test/rg",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"dev"},
				"suffixes":      []interface{}{"prod"},
				"random_length": 0,
				"clean_input":   false,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":  false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataName().Schema, tt.data)
			err := getNameReadResult(d, nil)

			if tt.wantErr && err == nil {
				t.Error("getNameReadResult() expected error, got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("getNameReadResult() unexpected error: %v", err)
			}
		})
	}
}
