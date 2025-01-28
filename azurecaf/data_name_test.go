//go:build unit

package azurecaf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataNameRead(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid resource group name",
			data: map[string]interface{}{
				"name":          "rg-test-123_(prod)",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{},
				"suffixes":      []interface{}{},
				"random_length": 0,
				"clean_input":   false,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":  true,
			},
			wantErr: false,
		},
		{
			name: "resource group name with special characters",
			data: map[string]interface{}{
				"name":          "rg-test-123.(prod)",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{},
				"suffixes":      []interface{}{},
				"random_length": 0,
				"clean_input":   false,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":  true,
			},
			wantErr: false,
		},
		{
			name: "invalid resource type",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "invalid_type",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataName().Schema, tt.data)
			diags := dataNameRead(context.Background(), d, nil)

			if tt.wantErr && len(diags) == 0 {
				t.Error("dataNameRead() expected error, got none")
			}
			if !tt.wantErr && len(diags) > 0 {
				t.Errorf("dataNameRead() unexpected error: %v", diags)
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
				"name":          "rg-test-123_(prod)_dev",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{},
				"suffixes":      []interface{}{},
				"random_length": 0,
				"clean_input":   false,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":  true,
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
