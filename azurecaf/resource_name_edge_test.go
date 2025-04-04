package azurecaf

import (
	"strings"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResourceName_EdgeCases(t *testing.T) {
	testCases := []struct {
		name           string
		resourceType   string
		separator      string
		prefixes       []string
		inputName      string
		suffixes       []string
		randomSuffix   string
		convention     string
		cleanInput     bool
		passthrough    bool
		useSlug        bool
		namePrecedence []string
		expectError    bool
		errorContains  string
	}{
		{
			name:           "Invalid resource type",
			resourceType:   "invalid_type",
			separator:      "-",
			prefixes:       []string{},
			inputName:      "test",
			suffixes:       []string{},
			randomSuffix:   "",
			convention:     ConventionCafClassic,
			cleanInput:     true,
			passthrough:    false,
			useSlug:        true,
			namePrecedence: []string{"name", "slug", "random", "suffixes", "prefixes"},
			expectError:    true,
			errorContains:  "invalid resource type",
		},
		{
			name:           "Invalid regex pattern",
			resourceType:   "azurerm_resource_group",
			separator:      "-",
			prefixes:       []string{"!@#"},
			inputName:      "test!@#",
			suffixes:       []string{"!@#"},
			randomSuffix:   "!@#",
			convention:     ConventionCafClassic,
			cleanInput:     false,
			passthrough:    false,
			useSlug:        true,
			namePrecedence: []string{"name", "slug", "random", "suffixes", "prefixes"},
			expectError:    true,
			errorContains:  "doesn't match",
		},
		{
			name:           "Empty name with passthrough",
			resourceType:   "azurerm_resource_group",
			separator:      "-",
			prefixes:       []string{},
			inputName:      "",
			suffixes:       []string{},
			randomSuffix:   "",
			convention:     ConventionCafClassic,
			cleanInput:     true,
			passthrough:    true,
			useSlug:        true,
			namePrecedence: []string{"name", "slug", "random", "suffixes", "prefixes"},
			expectError:    true,
			errorContains:  "doesn't match",
		},
		{
			name:           "Max length exactly at limit",
			resourceType:   "azurerm_resource_group",
			separator:      "-",
			prefixes:       []string{"prefix"},
			inputName:      "exactlength",
			suffixes:       []string{"suffix"},
			randomSuffix:   "random",
			convention:     ConventionCafClassic,
			cleanInput:     true,
			passthrough:    false,
			useSlug:        true,
			namePrecedence: []string{"name", "slug", "random", "suffixes", "prefixes"},
			expectError:    false,
		},
		{
			name:           "Special characters with clean input",
			resourceType:   "azurerm_resource_group",
			separator:      "-",
			prefixes:       []string{"prefix!@#"},
			inputName:      "test!@#name",
			suffixes:       []string{"suffix!@#"},
			randomSuffix:   "random!@#",
			convention:     ConventionCafClassic,
			cleanInput:     true,
			passthrough:    false,
			useSlug:        true,
			namePrecedence: []string{"name", "slug", "random", "suffixes", "prefixes"},
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getResourceName(
				tc.resourceType,
				tc.separator,
				tc.prefixes,
				tc.inputName,
				tc.suffixes,
				tc.randomSuffix,
				tc.convention,
				tc.cleanInput,
				tc.passthrough,
				tc.useSlug,
				tc.namePrecedence,
			)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got no error", tc.errorContains)
				} else if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == "" {
					t.Error("Expected non-empty result")
				}
			}
		})
	}
}

func TestGetNameResult_EdgeCases(t *testing.T) {
	testCases := []struct {
		name          string
		resourceType  string
		resourceTypes []string
		expectError   bool
		errorContains string
		setup         func(*schema.ResourceData)
	}{
		{
			name:          "Empty resource type and types",
			resourceType:  "",
			resourceTypes: []string{},
			expectError:   true,
			errorContains: "empty",
			setup: func(d *schema.ResourceData) {
				d.Set("name", "test")
				d.Set("prefixes", []interface{}{})
				d.Set("suffixes", []interface{}{})
				d.Set("separator", "-")
				d.Set("clean_input", true)
				d.Set("passthrough", false)
				d.Set("use_slug", true)
				d.Set("random_length", 5)
				d.Set("random_seed", 123)
			},
		},
		{
			name:          "Invalid resource in resource_types",
			resourceType:  "",
			resourceTypes: []string{"invalid_type"},
			expectError:   true,
			errorContains: "invalid resource type",
			setup: func(d *schema.ResourceData) {
				d.Set("name", "test")
				d.Set("prefixes", []interface{}{})
				d.Set("suffixes", []interface{}{})
				d.Set("separator", "-")
				d.Set("clean_input", true)
				d.Set("passthrough", false)
				d.Set("use_slug", true)
				d.Set("random_length", 5)
				d.Set("random_seed", 123)
			},
		},
		{
			name:          "Multiple valid resource types",
			resourceType:  "",
			resourceTypes: []string{"azurerm_resource_group", "azurerm_container_registry", "azurerm_application_gateway"},
			expectError:   false,
			setup: func(d *schema.ResourceData) {
				d.Set("name", "test")
				d.Set("prefixes", []interface{}{})
				d.Set("suffixes", []interface{}{})
				d.Set("separator", "-")
				d.Set("clean_input", true)
				d.Set("passthrough", false)
				d.Set("use_slug", true)
				d.Set("random_length", 5)
				d.Set("random_seed", 123)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := schema.Resource{Schema: resourceName().Schema}
			d := r.TestResourceData()
			tc.setup(d)
			d.Set("resource_type", tc.resourceType)
			d.Set("resource_types", tc.resourceTypes)

			err := getNameResult(d, nil)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got no error", tc.errorContains)
				} else if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tc.resourceType != "" {
					if d.Get("result").(string) == "" {
						t.Error("Expected non-empty result")
					}
				}
				if len(tc.resourceTypes) > 0 {
					results := d.Get("results").(map[string]interface{})
					if len(results) != len(tc.resourceTypes) {
						t.Errorf("Expected %d results, got %d", len(tc.resourceTypes), len(results))
					}
				}
			}
		})
	}
}
