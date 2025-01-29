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
		expected string
	}{
		{
			name: "valid resource group name",
			data: map[string]interface{}{
				"name":          "rg-test-123",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"dev"},
				"suffixes":      []interface{}{"001"},
				"random_length": 0,
				"random_seed":   12345,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":     true,
				"passthrough":   false,
			},
			wantErr:  false,
			expected: "dev-rg-test-123-001",
		},
		{
			name: "storage account with random string",
			data: map[string]interface{}{
				"name":          "data",
				"resource_type": "azurerm_storage_account",
				"prefixes":      []interface{}{"dev"},
				"random_length": 5,
				"random_seed":   12345,
				"clean_input":   true,
				"separator":     "",
				"use_slug":     true,
				"passthrough":   false,
			},
			wantErr:  false,
			expected: "devstdataisjlq",
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
			if !tt.wantErr {
				result := d.Get("result").(string)
				if result != tt.expected {
					t.Errorf("dataNameRead() got result = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestGetNameReadResult(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		wantErr  bool
		expected string
	}{
		{
			name: "key vault with random string",
			data: map[string]interface{}{
				"name":          "secrets",
				"resource_type": "azurerm_key_vault",
				"prefixes":      []interface{}{"dev"},
				"random_length": 5,
				"random_seed":   12345,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":   false,
			},
			wantErr:  false,
			expected: "dev-secrets-isjlq",
		},
		{
			name: "function app without random string",
			data: map[string]interface{}{
				"name":          "api",
				"resource_type": "azurerm_function_app",
				"prefixes":      []interface{}{"dev"},
				"suffixes":      []interface{}{"func"},
				"random_length": 0,
				"random_seed":   12345,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":     false,
				"passthrough":   false,
			},
			wantErr:  false,
			expected: "dev-api-func",
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
			if !tt.wantErr {
				result := d.Get("result").(string)
				if result != tt.expected {
					t.Errorf("getNameReadResult() got result = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
