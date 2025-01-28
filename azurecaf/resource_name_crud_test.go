//go:build unit

package azurecaf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceNameCreate(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid resource",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"prefix"},
				"suffixes":      []interface{}{"suffix"},
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
			d := schema.TestResourceDataRaw(t, resourceName().Schema, tt.data)
			diags := resourceNameCreate(context.Background(), d, nil)

			if tt.wantErr && len(diags) == 0 {
				t.Error("resourceNameCreate() expected error, got none")
			}
			if !tt.wantErr && len(diags) > 0 {
				t.Errorf("resourceNameCreate() unexpected errors: %v", diags)
			}
		})
	}
}

func TestResourceNameUpdate(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid update",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"new-prefix"},
				"suffixes":      []interface{}{"new-suffix"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceName().Schema, tt.data)
			diags := resourceNameUpdate(context.Background(), d, nil)

			if tt.wantErr && len(diags) == 0 {
				t.Error("resourceNameUpdate() expected error, got none")
			}
			if !tt.wantErr && len(diags) > 0 {
				t.Errorf("resourceNameUpdate() unexpected errors: %v", diags)
			}
		})
	}
}

func TestResourceNameRead(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid read",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"result":        "prefix-test-suffix",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceName().Schema, tt.data)
			diags := resourceNameRead(context.Background(), d, nil)

			if tt.wantErr && len(diags) == 0 {
				t.Error("resourceNameRead() expected error, got none")
			}
			if !tt.wantErr && len(diags) > 0 {
				t.Errorf("resourceNameRead() unexpected errors: %v", diags)
			}
		})
	}
}

func TestGetDifference(t *testing.T) {
	tests := []struct {
		name     string
		old      interface{}
		new      interface{}
		expected bool
	}{
		{
			name:     "different values",
			old:      "old",
			new:      "new",
			expected: true,
		},
		{
			name:     "same values",
			old:      "same",
			new:      "same",
			expected: false,
		},
		{
			name:     "nil old value",
			old:      nil,
			new:      "new",
			expected: true,
		},
		{
			name:     "nil new value",
			old:      "old",
			new:      nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDifference(tt.old, tt.new)
			if result != tt.expected {
				t.Errorf("getDifference() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertInterfaceToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string input",
			input:    "test",
			expected: "test",
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "integer input",
			input:    123,
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertInterfaceToString(tt.input)
			if result != tt.expected {
				t.Errorf("convertInterfaceToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}
