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

func TestConvertInterfaceToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []string
	}{
		{
			name:     "string slice input",
			input:    []interface{}{"test", "test2"},
			expected: []string{"test", "test2"},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []string{},
		},
		{
			name:     "empty slice input",
			input:    []interface{}{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertInterfaceToString(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("convertInterfaceToString() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("convertInterfaceToString() at index %d = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}
