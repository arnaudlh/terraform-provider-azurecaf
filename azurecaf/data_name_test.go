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
			name: "valid resource type",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"rg"},
				"suffixes":      []interface{}{"prod"},
				"random_length": 5,
				"clean_input":   true,
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
			name: "valid input",
			data: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"rg"},
				"suffixes":      []interface{}{"prod"},
				"random_length": 5,
				"clean_input":   true,
			},
			wantErr: false,
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
