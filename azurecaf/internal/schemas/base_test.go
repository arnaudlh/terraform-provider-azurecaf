package schemas

import (
	"context"
	"reflect"
	"testing"
)

func TestV4_Schema(t *testing.T) {
	resource := V4_Schema()
	if resource == nil {
		t.Fatal("V4_Schema() returned nil")
	}
	if resource.Schema == nil {
		t.Fatal("V4_Schema().Schema is nil")
	}
	s := resource.Schema
	
	// Test required fields
	requiredFields := []string{"name", "resource_type"}
	for _, field := range requiredFields {
		if v, ok := s[field]; !ok {
			t.Errorf("missing required field %s", field)
		} else if !v.Required {
			t.Errorf("field %s should be required", field)
		}
	}

	// Test computed fields
	computedFields := []string{"result", "results"}
	for _, field := range computedFields {
		if v, ok := s[field]; !ok {
			t.Errorf("missing computed field %s", field)
		} else if !v.Computed {
			t.Errorf("field %s should be computed", field)
		}
	}

	// Test optional fields
	optionalFields := []string{"prefixes", "suffixes", "random_length", "random_seed", "separator", "clean_input", "passthrough", "use_slug"}
	for _, field := range optionalFields {
		if _, ok := s[field]; !ok {
			t.Errorf("missing optional field %s", field)
		}
	}

	// Test field types
	fieldTypes := map[string]string{
		"name":          "TypeString",
		"resource_type": "TypeString",
		"prefixes":      "TypeList",
		"suffixes":      "TypeList",
	}
	for field, expectedType := range fieldTypes {
		if v, ok := s[field]; !ok {
			t.Errorf("missing field %s", field)
		} else if v.Type.String() != expectedType {
			t.Errorf("field %s should be %s, got %s", field, expectedType, v.Type.String())
		}
	}
}

func TestResourceNameStateUpgradeV2(t *testing.T) {
	tests := []struct {
		name    string
		oldData map[string]interface{}
		want    map[string]interface{}
	}{
		{
			name: "basic upgrade",
			oldData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"prefix1", "prefix2"},
				"suffixes":      []interface{}{"suffix1", "suffix2"},
				"random_length": 5,
				"random_seed":   42,
			},
			want: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"prefix1", "prefix2"},
				"suffixes":      []interface{}{"suffix1", "suffix2"},
				"random_length": 5,
				"random_seed":   42,
			},
		},
		{
			name: "empty fields",
			oldData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
			},
			want: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
			},
		},
		{
			name: "with computed fields",
			oldData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"result":        "test-rg",
				"results": map[string]interface{}{
					"azurerm_resource_group": "test-rg",
				},
			},
			want: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"result":        "test-rg",
				"results": map[string]interface{}{
					"azurerm_resource_group": "test-rg",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newData, err := ResourceNameStateUpgradeV2(context.TODO(), tt.oldData, nil)
			if err != nil {
				t.Fatalf("error upgrading state: %s", err)
			}

			for k, v := range tt.want {
				got, ok := newData[k]
				if !ok {
					t.Errorf("field %s missing from upgraded state", k)
					continue
				}
				switch val := v.(type) {
				case []interface{}:
					gotSlice, ok := got.([]interface{})
					if !ok {
						t.Errorf("field %s: expected slice but got %T", k, got)
						continue
					}
					if !reflect.DeepEqual(gotSlice, val) {
						t.Errorf("field %s = %v, want %v", k, gotSlice, val)
					}
				case map[string]interface{}:
					gotMap, ok := got.(map[string]interface{})
					if !ok {
						t.Errorf("field %s: expected map but got %T", k, got)
						continue
					}
					if !reflect.DeepEqual(gotMap, val) {
						t.Errorf("field %s = %v, want %v", k, gotMap, val)
					}
				default:
					if !reflect.DeepEqual(got, v) {
						t.Errorf("field %s = %v, want %v", k, got, v)
					}
				}
			}
		})
	}
}

func TestResourceNameStateUpgradeV3(t *testing.T) {
	tests := []struct {
		name    string
		oldData map[string]interface{}
		want    map[string]interface{}
	}{
		{
			name: "basic upgrade",
			oldData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"prefix1", "prefix2"},
				"suffixes":      []interface{}{"suffix1", "suffix2"},
				"random_length": 5,
				"random_seed":   42,
			},
			want: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"prefixes":      []interface{}{"prefix1", "prefix2"},
				"suffixes":      []interface{}{"suffix1", "suffix2"},
				"random_length": 5,
				"random_seed":   42,
			},
		},
		{
			name: "empty fields",
			oldData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
			},
			want: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
			},
		},
		{
			name: "with use_slug",
			oldData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"use_slug":      true,
			},
			want: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_resource_group",
				"use_slug":      true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newData, err := ResourceNameStateUpgradeV3(context.TODO(), tt.oldData, nil)
			if err != nil {
				t.Fatalf("error upgrading state: %s", err)
			}

			for k, v := range tt.want {
				got, ok := newData[k]
				if !ok {
					t.Errorf("field %s missing from upgraded state", k)
					continue
				}
				switch val := v.(type) {
				case []interface{}:
					gotSlice, ok := got.([]interface{})
					if !ok {
						t.Errorf("field %s: expected slice but got %T", k, got)
						continue
					}
					if !reflect.DeepEqual(gotSlice, val) {
						t.Errorf("field %s = %v, want %v", k, gotSlice, val)
					}
				case map[string]interface{}:
					gotMap, ok := got.(map[string]interface{})
					if !ok {
						t.Errorf("field %s: expected map but got %T", k, got)
						continue
					}
					if !reflect.DeepEqual(gotMap, val) {
						t.Errorf("field %s = %v, want %v", k, gotMap, val)
					}
				default:
					if !reflect.DeepEqual(got, v) {
						t.Errorf("field %s = %v, want %v", k, got, v)
					}
				}
			}
		})
	}
}
