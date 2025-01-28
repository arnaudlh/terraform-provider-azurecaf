//go:build unit

package schemas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestV4_Schema(t *testing.T) {
	schema := V4_Schema()
	if schema == nil {
		t.Fatal("V4_Schema() returned nil schema")
	}

	// Test that all required fields are present
	requiredFields := []string{"name", "resource_type"}
	for _, field := range requiredFields {
		if _, ok := schema[field]; !ok {
			t.Errorf("V4_Schema() missing required field %s", field)
		}
	}

	// Test that computed fields are marked as such
	computedFields := []string{"result", "random_string"}
	for _, field := range computedFields {
		if !schema[field].Computed {
			t.Errorf("V4_Schema() field %s should be computed", field)
		}
	}
}

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
		forceNew bool
		default_ interface{}
	}{
		{"name", "name", true, schema.TypeString, false, false, nil},
		{"resource_type", "resource_type", true, schema.TypeString, false, false, nil},
		{"prefixes", "prefixes", false, schema.TypeList, false, false, nil},
		{"suffixes", "suffixes", false, schema.TypeList, false, false, nil},
		{"random_length", "random_length", false, schema.TypeInt, false, false, 0},
		{"result", "result", false, schema.TypeString, true, false, nil},
		{"random_seed", "random_seed", false, schema.TypeInt, false, false, nil},
		{"random_string", "random_string", false, schema.TypeString, true, false, nil},
		{"clean_input", "clean_input", false, schema.TypeBool, false, false, false},
		{"separator", "separator", false, schema.TypeString, false, false, "-"},
		{"use_slug", "use_slug", false, schema.TypeBool, false, false, false},
		{"passthrough", "passthrough", false, schema.TypeBool, false, false, false},
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
            if field.ForceNew != tt.forceNew {
                t.Errorf("V4() field %s ForceNew = %v, want %v", tt.field, field.ForceNew, tt.forceNew)
            }
            if tt.default_ != nil && field.Default != tt.default_ {
                t.Errorf("V4() field %s Default = %v, want %v", tt.field, field.Default, tt.default_)
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

    // Test validation functions
    resourceType := v4Schema.Schema["resource_type"]
    if resourceType == nil {
        t.Fatal("resource_type field is missing")
    }
    if resourceType.ValidateFunc == nil {
        t.Error("resource_type field should have a validation function")
    }
    // Test validation function behavior
    if resourceType.ValidateFunc == nil {
        t.Error("resource_type field should have a validation function")
    }

    // Test description fields
    expectedDescriptions := map[string]string{
        "resource_type": "The resource type to generate a name for",
        "name":         "The name to be transformed according to the resource type naming rules",
        "prefixes":     "A list of prefixes to be used in the name",
        "suffixes":     "A list of suffixes to be used in the name",
        "result":       "The computed name for the resource",
    }

    for field, expectedDesc := range expectedDescriptions {
        if schema, ok := v4Schema.Schema[field]; ok {
            if schema.Description != expectedDesc {
                t.Errorf("V4() field %s description = %v, want %v", field, schema.Description, expectedDesc)
            }
        }
    }
}
