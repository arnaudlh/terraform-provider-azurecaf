//go:build unit

package schemas

import (
	"context"
	"testing"
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
				"result": "",
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
	v2Schema := V2()
	if v2Schema == nil {
		t.Fatal("V2() returned nil schema")
	}

	tests := []struct {
		name     string
		field    string
		required bool
		typ      string
		computed bool
	}{
		{"name", "name", true, "TypeString", false},
		{"resource_type", "resource_type", true, "TypeString", false},
		{"prefixes", "prefixes", false, "TypeList", false},
		{"suffixes", "suffixes", false, "TypeList", false},
		{"result", "result", false, "TypeString", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := v2Schema[tt.field]
			if !ok {
				t.Errorf("V2() schema missing field %s", tt.field)
				return
			}
			if field.Required != tt.required {
				t.Errorf("V2() field %s Required = %v, want %v", tt.field, field.Required, tt.required)
			}
			if field.Type.String() != tt.typ {
				t.Errorf("V2() field %s Type = %v, want %v", tt.field, field.Type, tt.typ)
			}
			if field.Computed != tt.computed {
				t.Errorf("V2() field %s Computed = %v, want %v", tt.field, field.Computed, tt.computed)
			}
		})
	}
}
