//go:build unit

package schemas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestV4(t *testing.T) {
	v4Schema := V4()
	if v4Schema == nil {
		t.Fatal("V4() returned nil schema")
	}

	tests := []struct {
		name     string
		field    string
		required bool
		typ      schema.ValueType
		computed bool
	}{
		{"name", "name", true, schema.TypeString, false},
		{"resource_type", "resource_type", true, schema.TypeString, false},
		{"prefixes", "prefixes", false, schema.TypeList, false},
		{"suffixes", "suffixes", false, schema.TypeList, false},
		{"random_length", "random_length", false, schema.TypeInt, false},
		{"result", "result", false, schema.TypeString, true},
		{"random_seed", "random_seed", false, schema.TypeInt, false},
		{"random_string", "random_string", false, schema.TypeString, true},
		{"clean_input", "clean_input", false, schema.TypeBool, false},
		{"separator", "separator", false, schema.TypeString, false},
		{"use_slug", "use_slug", false, schema.TypeBool, false},
		{"passthrough", "passthrough", false, schema.TypeBool, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := v4Schema.Schema[tt.field]
			if !ok {
				t.Errorf("V4() schema missing field %s", tt.field)
				return
			}
			if field.Required != tt.required {
				t.Errorf("V4() field %s Required = %v, want %v", tt.field, field.Required, tt.required)
			}
			if field.Type != tt.typ {
				t.Errorf("V4() field %s Type = %v, want %v", tt.field, field.Type, tt.typ)
			}
			if field.Computed != tt.computed {
				t.Errorf("V4() field %s Computed = %v, want %v", tt.field, field.Computed, tt.computed)
			}
		})
	}

	// Test element types for list fields
	listFields := map[string]schema.ValueType{
		"prefixes": schema.TypeString,
		"suffixes": schema.TypeString,
	}

	for field, elemType := range listFields {
		if v4Schema.Schema[field].Type != schema.TypeList {
			t.Errorf("V4() field %s should be TypeList", field)
			continue
		}
		elem, ok := v4Schema.Schema[field].Elem.(*schema.Schema)
		if !ok {
			t.Errorf("V4() field %s Elem should be *schema.Schema", field)
			continue
		}
		if elem.Type != elemType {
			t.Errorf("V4() field %s Elem.Type = %v, want %v", field, elem.Type, elemType)
		}
	}
}
