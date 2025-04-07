package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResourceCoverage(t *testing.T) {
	resourceType := "rg"
	resource, err := getResource(resourceType)

	if err != nil {
		t.Fatalf("getResource returned error for valid mapped resource: %v", err)
	}

	if resource == nil {
		t.Fatal("Expected resource to be returned for valid mapped resource")
	}

	resourceType = "nonexistent_resource"
	_, err = getResource(resourceType)

	if err == nil {
		t.Fatal("Expected error for nonexistent resource but got none")
	}
}

func TestTrimResourceNameCoverage(t *testing.T) {
	resourceName := ""
	maxLength := 10

	result := trimResourceName(resourceName, maxLength)

	if result != "" {
		t.Fatalf("Expected empty string, got %s", result)
	}

	resourceName = "test"
	maxLength = 0

	result = trimResourceName(resourceName, maxLength)

	if result != "" {
		t.Fatalf("Expected empty string, got %s", result)
	}
}

func TestValidateResourceTypeCoverage(t *testing.T) {
	resourceType := ""
	resourceTypes := []string{}
	valid, err := validateResourceType(resourceType, resourceTypes)

	if valid {
		t.Fatal("Expected empty resource type and empty resource types to be invalid")
	}

	if err == nil {
		t.Fatal("Expected error for empty resource type and empty resource types but got none")
	}

	resourceType = "rg"
	resourceTypes = []string{}
	valid, err = validateResourceType(resourceType, resourceTypes)

	if !valid {
		t.Fatal("Expected valid resource type in ResourceMaps to be valid")
	}

	if err != nil {
		t.Fatalf("validateResourceType returned error for valid resource type in ResourceMaps: %v", err)
	}
}

func TestGetResourceNameCoverage(t *testing.T) {
	resourceType := "azurerm_resource_group"
	separator := "-"
	prefixes := []string{"prefix"}
	name := ""
	suffixes := []string{"suffix"}
	randomSuffix := "random"
	convention := ConventionCafClassic
	cleanInput := true
	passthrough := false
	useSlug := true
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)

	if err != nil {
		t.Fatalf("getResourceName returned error for empty name: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result for empty name")
	}

	resourceType = "azurerm_resource_group"
	separator = "-"
	prefixes = []string{}
	name = "testname"
	suffixes = []string{}
	randomSuffix = "random"
	convention = ConventionCafClassic
	cleanInput = true
	passthrough = false
	useSlug = true
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}

	result, err = getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)

	if err != nil {
		t.Fatalf("getResourceName returned error for empty prefixes and suffixes: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result for empty prefixes and suffixes")
	}
}

func TestGetNameResultCoverage(t *testing.T) {
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
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"separator": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"azurerm_storage_account"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}

func TestGetResultCoverage(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefix": {
				Type: schema.TypeString,
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"postfix": {
				Type: schema.TypeString,
			},
			"convention": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"max_length": {
				Type: schema.TypeInt,
			},
			"result": {
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 10) // Less than default max length

	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if len(result) > 10 {
		t.Fatalf("Expected result length <= 10, got %d", len(result))
	}
}

func TestResourceNameCreateCoverage(t *testing.T) {
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
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"separator": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"azurerm_storage_account"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := resourceNameCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNameCreate returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}

func TestResourceNameReadCoverage(t *testing.T) {
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
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"separator": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"azurerm_storage_account"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)
	d.Set("result", "test-result")

	err := resourceNameRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNameRead returned error: %v", err)
	}
}

func TestResourceNamingConventionCreateCoverage(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefix": {
				Type: schema.TypeString,
			},
			"postfix": {
				Type: schema.TypeString,
			},
			"convention": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"max_length": {
				Type: schema.TypeInt,
			},
			"result": {
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err := resourceNamingConventionCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionCreate returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}

func TestResourceNamingConventionReadCoverage(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefix": {
				Type: schema.TypeString,
			},
			"postfix": {
				Type: schema.TypeString,
			},
			"convention": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"max_length": {
				Type: schema.TypeInt,
			},
			"result": {
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)
	d.Set("result", "test-result")

	err := resourceNamingConventionRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionRead returned error: %v", err)
	}
}

func TestCleanStringCoverage(t *testing.T) {
	data := "test_data-123!@#$%^&*()"
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanString(data, &resource)
	expected := "test_data-123()"

	if result != expected {
		t.Fatalf("Expected %s, got %s", expected, result)
	}
}

func TestConvertInterfaceToStringCoverage(t *testing.T) {
	data := []interface{}{"test1", "test2", "test3"}
	result := convertInterfaceToString(data)
	expected := []string{"test1", "test2", "test3"}

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("Expected %s at index %d, got %s", expected[i], i, result[i])
		}
	}
}

func TestDataNameReadCoverage(t *testing.T) {
	t.Skip("Skipping test that requires more complex setup")
}

func TestGetNameReadResultCoverage(t *testing.T) {
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
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := getNameReadResult(d, nil)
	if err != nil {
		t.Fatalf("getNameReadResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}

func TestRandSeqCoverage(t *testing.T) {
	length := 10
	var seed int64 = 123
	result := randSeq(length, &seed)

	if len(result) != length {
		t.Fatalf("Expected length %d, got %d", length, len(result))
	}

	var seed2 int64 = 456
	result2 := randSeq(length, &seed2)

	if result == result2 {
		t.Fatal("Expected different results with different seeds")
	}
}
