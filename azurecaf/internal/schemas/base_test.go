package schemas

import (
	"testing"
)

func TestV4_Schema(t *testing.T) {
	s := V4_Schema()
	
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
	computedFields := []string{"result", "results", "random_string"}
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
	if s["name"].Type.String() != "TypeString" {
		t.Errorf("field name should be TypeString")
	}
	if s["resource_type"].Type.String() != "TypeString" {
		t.Errorf("field resource_type should be TypeString")
	}
	if s["prefixes"].Type.String() != "TypeList" {
		t.Errorf("field prefixes should be TypeList")
	}
	if s["suffixes"].Type.String() != "TypeList" {
		t.Errorf("field suffixes should be TypeList")
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
			newData, err := ResourceNameStateUpgradeV2(nil, tt.oldData, nil)
			if err != nil {
				t.Fatalf("error upgrading state: %s", err)
			}

			for k, v := range tt.want {
				if got, ok := newData[k]; !ok || got != v {
					t.Errorf("field %s = %v, want %v", k, got, v)
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
			newData, err := ResourceNameStateUpgradeV3(nil, tt.oldData, nil)
			if err != nil {
				t.Fatalf("error upgrading state: %s", err)
			}

			for k, v := range tt.want {
				if got, ok := newData[k]; !ok || got != v {
					t.Errorf("field %s = %v, want %v", k, got, v)
				}
			}
		})
	}
}
