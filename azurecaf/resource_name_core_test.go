package azurecaf

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResourceName_AllComponents(t *testing.T) {
	result, err := getResourceName(
		"azurerm_resource_group",
		"-",
		[]string{"prefix1", "prefix2"},
		"testname",
		[]string{"suffix1", "suffix2"},
		"random123",
		ConventionCafClassic,
		true,
		false,
		true,
		[]string{"name", "slug", "random", "suffixes", "prefixes"},
	)

	if err != nil {
		t.Fatalf("getResourceName returned unexpected error: %v", err)
	}

	expectedComponents := []string{"prefix1", "prefix2", "rg", "testname", "random123", "suffix1", "suffix2"}
	for _, component := range expectedComponents {
		if !strings.Contains(result, component) {
			t.Errorf("Expected result to contain %s, got %s", component, result)
		}
	}

	validationRegex := ResourceDefinitions["azurerm_resource_group"].ValidationRegExp
	regex := regexp.MustCompile(validationRegex)
	if !regex.MatchString(result) {
		t.Errorf("Result %s does not match validation regex %s", result, validationRegex)
	}
}

func TestGetResourceName_MaxLength(t *testing.T) {
	longName := strings.Repeat("a", 100)
	result, err := getResourceName(
		"azurerm_resource_group",
		"-",
		[]string{"prefix"},
		longName,
		[]string{"suffix"},
		"random",
		ConventionCafClassic,
		true,
		false,
		true,
		[]string{"name", "slug", "random", "suffixes", "prefixes"},
	)

	if err != nil {
		t.Fatalf("getResourceName returned unexpected error: %v", err)
	}

	maxLength := ResourceDefinitions["azurerm_resource_group"].MaxLength
	if len(result) > maxLength {
		t.Errorf("Expected result length <= %d, got %d", maxLength, len(result))
	}
}

func TestGetResourceName_LowercaseEnforcement(t *testing.T) {
	result, err := getResourceName(
		"azurerm_storage_account",
		"-",
		[]string{"PREFIX"},
		"TestName",
		[]string{"SUFFIX"},
		"RANDOM",
		ConventionCafClassic,
		true,
		false,
		true,
		[]string{"name", "slug", "random", "suffixes", "prefixes"},
	)

	if err != nil {
		t.Fatalf("getResourceName returned unexpected error: %v", err)
	}

	if strings.ToLower(result) != result {
		t.Errorf("Expected lowercase result for storage account, got %s", result)
	}
}

func TestGetNameResult_SingleResource(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Default:  []interface{}{},
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Default:  []interface{}{},
			},
			"separator": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"clean_input": {
				Type: schema.TypeBool,
			},
			"passthrough": {
				Type: schema.TypeBool,
			},
			"use_slug": {
				Type: schema.TypeBool,
			},
			"random_length": {
				Type: schema.TypeInt,
			},
			"random_seed": {
				Type: schema.TypeInt,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"results": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{})
	d.Set("suffixes", []interface{}{})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned unexpected error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	validationRegex := ResourceDefinitions["azurerm_resource_group"].ValidationRegExp
	regex := regexp.MustCompile(validationRegex)
	if !regex.MatchString(result) {
		t.Errorf("Result %s does not match validation regex %s", result, validationRegex)
	}
}

func TestComposeName_PrecedenceOrder(t *testing.T) {
	result := composeName(
		"-",
		[]string{"prefix1", "prefix2"},
		"name",
		"slug",
		[]string{"suffix1", "suffix2"},
		"random",
		100,
		[]string{"name", "slug", "random", "suffixes", "prefixes"},
	)

	components := strings.Split(result, "-")

	expectedOrder := []string{"prefix1", "prefix2", "slug", "name", "random", "suffix1", "suffix2"}
	if !reflect.DeepEqual(components, expectedOrder) {
		t.Errorf("Expected component order %v, got %v", expectedOrder, components)
	}
}

func TestComposeName_MaxLengthTruncation(t *testing.T) {
	maxLength := 20
	result := composeName(
		"-",
		[]string{"prefix1", "prefix2"},
		"verylongname",
		"slug",
		[]string{"suffix1", "suffix2"},
		"random",
		maxLength,
		[]string{"name", "slug", "random", "suffixes", "prefixes"},
	)

	if len(result) > maxLength {
		t.Errorf("Expected result length <= %d, got %d", maxLength, len(result))
	}

	if !strings.Contains(result, "name") {
		t.Error("Expected truncated result to contain high-priority 'name' component")
	}
}

func TestCleanString_SpecialCharacters(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		resource string
		expected string
	}{
		{
			name:     "Resource Group Special Characters",
			input:    "test!@#$%^&*()name",
			resource: "azurerm_resource_group",
			expected: "test()name",
		},
		{
			name:     "Storage Account Special Characters",
			input:    "test!@#$%^&*()name",
			resource: "azurerm_storage_account",
			expected: "testname",
		},
		{
			name:     "Key Vault Special Characters",
			input:    "test!@#$%^&*()name",
			resource: "azurerm_key_vault",
			expected: "testname",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource := ResourceDefinitions[tc.resource]
			result := cleanString(tc.input, &resource)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestGetSlug_AllResourceTypes(t *testing.T) {
	for resourceType, resource := range ResourceDefinitions {
		slug := getSlug(resourceType, ConventionCafClassic)
		if resource.CafPrefix != "" && slug != resource.CafPrefix {
			t.Errorf("For resource type %s, expected slug %s, got %s", resourceType, resource.CafPrefix, slug)
		}
	}
}
