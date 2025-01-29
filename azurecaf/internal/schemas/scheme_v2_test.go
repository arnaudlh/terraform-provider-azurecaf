//go:build unit
// +build unit

package schemas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceNameStateUpgradeV2(t *testing.T) {
	tests := []struct {
		name    string
		state   map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "empty state",
			state: map[string]interface{}{},
			want: map[string]interface{}{
				"use_slug": true,
				"result":   nil,
			},
			wantErr: false,
		},
		{
			name: "with result",
			state: map[string]interface{}{
				"result": "test-resource",
			},
			want: map[string]interface{}{
				"result": "test-resource",
			},
			wantErr: false,
		},
		{
			name: "with random string",
			state: map[string]interface{}{
				"result":        "test-resource",
				"random_string": "abc123",
			},
			want: map[string]interface{}{
				"result":        "test-resource",
				"random_string": "abc123",
			},
			wantErr: false,
		},
		{
			name: "with all fields",
			state: map[string]interface{}{
				"result":        "test-resource",
				"random_string": "abc123",
				"random_seed":   42,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":      true,
				"passthrough":   false,
			},
			want: map[string]interface{}{
				"result":        "test-resource",
				"random_string": "abc123",
				"random_seed":   42,
				"clean_input":   true,
				"separator":     "-",
				"use_slug":      true,
				"passthrough":   false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResourceNameStateUpgradeV2(context.Background(), tt.state, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceNameStateUpgradeV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("ResourceNameStateUpgradeV2() got[%s] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

func TestV2(t *testing.T) {
	v2Resource := V2()
	if v2Resource == nil {
		t.Fatal("V2() returned nil resource")
	}
	if v2Resource.Schema == nil {
		t.Fatal("V2() schema is nil")
	}

	tests := []struct {
		name     string
		field    string
		required bool
		typ      schema.ValueType
		computed bool
		default_ interface{}
	}{
		{"name", "name", false, schema.TypeString, false, ""},
		{"resource_type", "resource_type", false, schema.TypeString, false, nil},
		{"prefixes", "prefixes", false, schema.TypeList, false, nil},
		{"suffixes", "suffixes", false, schema.TypeList, false, nil},
		{"result", "result", false, schema.TypeString, true, nil},
		{"clean_input", "clean_input", false, schema.TypeBool, false, true},
		{"separator", "separator", false, schema.TypeString, false, "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := v2Resource.Schema[tt.field]
			if !ok {
				t.Errorf("V2() schema missing field %s", tt.field)
				return
			}
			if field.Required != tt.required {
				t.Errorf("V2() field %s Required = %v, want %v", tt.field, field.Required, tt.required)
			}
			if field.Type != tt.typ {
				t.Errorf("V2() field %s Type = %v, want %v", tt.field, field.Type, tt.typ)
			}
			if field.Computed != tt.computed {
				t.Errorf("V2() field %s Computed = %v, want %v", tt.field, field.Computed, tt.computed)
			}
			if tt.default_ != nil && field.Default != tt.default_ {
				t.Errorf("V2() field %s Default = %v, want %v", tt.field, field.Default, tt.default_)
			}
		})
	}
}
