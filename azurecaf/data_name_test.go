//go:build unit

package azurecaf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
			name: "valid resource type",
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
			wantErr:  true,
			errCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataName().Schema, tt.data)
			diags := dataNameRead(context.Background(), d, nil)

			if tt.wantErr && len(diags) != tt.errCount {
				t.Errorf("dataNameRead() error count = %v, want %v", len(diags), tt.errCount)
			}
			if !tt.wantErr && len(diags) != 0 {
				t.Errorf("dataNameRead() unexpected errors: %v", diags)
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
			name: "valid input",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"prefix"},
				"suffixes":      []interface{}{"suffix"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataName().Schema, tt.data)
			diags := getNameReadResult(d, nil)

			if tt.wantErr && len(diags) == 0 {
				t.Error("getNameReadResult() expected error, got none")
			}
			if !tt.wantErr && len(diags) > 0 {
				t.Errorf("getNameReadResult() unexpected errors: %v", diags)
			}
		})
	}
}
