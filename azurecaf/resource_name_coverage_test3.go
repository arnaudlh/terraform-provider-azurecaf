package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCleanStringCoverage2(t *testing.T) {
	name := "test!@#$%^&*()_+name"
	resource := ResourceDefinitions["azurerm_resource_group"]
	
	result := cleanString(name, &resource)
	
	if result == name {
		t.Fatalf("Expected cleaned string, got '%s'", result)
	}
	
	name = "testname"
	result = cleanString(name, &resource)
	
	if result != name {
		t.Fatalf("Expected '%s', got '%s'", name, result)
	}
}

func TestConvertInterfaceToStringCoverage2(t *testing.T) {
	source := []interface{}{"test1", 2, true, 3.14}
	
	result := convertInterfaceToString(source)
	
	if len(result) != len(source) {
		t.Fatalf("Expected length %d, got %d", len(source), len(result))
	}
	
	if result[0] != "test1" {
		t.Fatalf("Expected 'test1', got '%s'", result[0])
	}
	
	if result[1] != "2" {
		t.Fatalf("Expected '2', got '%s'", result[1])
	}
	
	if result[2] != "true" {
		t.Fatalf("Expected 'true', got '%s'", result[2])
	}
	
	if result[3] != "3.14" {
		t.Fatalf("Expected '3.14', got '%s'", result[3])
	}
	
	source = []interface{}{}
	result = convertInterfaceToString(source)
	
	if len(result) != 0 {
		t.Fatalf("Expected empty slice, got length %d", len(result))
	}
}

func TestRandSeqCoverage2(t *testing.T) {
	length := 10
	seed := int64(123)
	
	result := randSeq(length, &seed)
	
	if len(result) != length {
		t.Fatalf("Expected length %d, got %d", length, len(result))
	}
	
	result2 := randSeq(length, &seed)
	
	if result == result2 {
		t.Fatalf("Expected different results with same seed pointer")
	}
	
	result3 := randSeq(length, nil)
	
	if len(result3) != length {
		t.Fatalf("Expected length %d, got %d", length, len(result3))
	}
	
	result4 := randSeq(0, nil)
	
	if result4 != "" {
		t.Fatalf("Expected empty string, got '%s'", result4)
	}
}

func TestGetResourceCoverage2(t *testing.T) {
	resourceType := "rg"
	resource, err := getResource(resourceType)
	
	if err != nil {
		t.Fatalf("getResource returned error for valid mapped resource: %v", err)
	}
	
	if resource == nil {
		t.Fatal("Expected resource to be returned for valid mapped resource")
	}
	
	resourceType = "azurerm_resource_group"
	resource, err = getResource(resourceType)
	
	if err != nil {
		t.Fatalf("getResource returned error for valid resource: %v", err)
	}
	
	if resource == nil {
		t.Fatal("Expected resource to be returned for valid resource")
	}
	
	resourceType = "nonexistent_resource"
	_, err = getResource(resourceType)
	
	if err == nil {
		t.Fatal("Expected error for nonexistent resource but got none")
	}
}

func TestTrimResourceNameCoverage2(t *testing.T) {
	resourceName := "this-is-a-very-long-resource-name-that-needs-to-be-trimmed"
	maxLength := 20
	
	result := trimResourceName(resourceName, maxLength)
	
	if len(result) != maxLength {
		t.Fatalf("Expected length %d, got %d", maxLength, len(result))
	}
	
	if result != resourceName[:maxLength] {
		t.Fatalf("Expected '%s', got '%s'", resourceName[:maxLength], result)
	}
	
	resourceName = "short-name"
	maxLength = 20
	
	result = trimResourceName(resourceName, maxLength)
	
	if result != resourceName {
		t.Fatalf("Expected '%s', got '%s'", resourceName, result)
	}
	
	resourceName = ""
	maxLength = 10
	
	result = trimResourceName(resourceName, maxLength)
	
	if result != "" {
		t.Fatalf("Expected empty string, got '%s'", result)
	}
	
	resourceName = "test"
	maxLength = 0
	
	result = trimResourceName(resourceName, maxLength)
	
	if result != "" {
		t.Fatalf("Expected empty string, got '%s'", result)
	}
}

func TestValidateResourceTypeCoverage2(t *testing.T) {
	resourceType := "azurerm_resource_group"
	resourceTypes := []string{}
	valid, err := validateResourceType(resourceType, resourceTypes)
	
	if !valid {
		t.Fatal("Expected valid resource type to be valid")
	}
	
	if err != nil {
		t.Fatalf("validateResourceType returned error for valid resource type: %v", err)
	}
	
	resourceType = ""
	resourceTypes = []string{"azurerm_resource_group", "azurerm_storage_account"}
	valid, err = validateResourceType(resourceType, resourceTypes)
	
	if !valid {
		t.Fatal("Expected valid resource types to be valid")
	}
	
	if err != nil {
		t.Fatalf("validateResourceType returned error for valid resource types: %v", err)
	}
	
	resourceType = "azurerm_resource_group"
	resourceTypes = []string{"azurerm_storage_account"}
	valid, err = validateResourceType(resourceType, resourceTypes)
	
	if !valid {
		t.Fatal("Expected valid resource type and valid resource types to be valid")
	}
	
	if err != nil {
		t.Fatalf("validateResourceType returned error for valid resource type and valid resource types: %v", err)
	}
	
	resourceType = "invalid_resource_type"
	resourceTypes = []string{}
	valid, err = validateResourceType(resourceType, resourceTypes)
	
	if valid {
		t.Fatal("Expected invalid resource type to be invalid")
	}
	
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
	
	resourceType = ""
	resourceTypes = []string{"invalid_resource_type"}
	valid, err = validateResourceType(resourceType, resourceTypes)
	
	if valid {
		t.Fatal("Expected invalid resource types to be invalid")
	}
	
	if err == nil {
		t.Fatal("Expected error for invalid resource types but got none")
	}
	
	resourceType = "invalid_resource_type"
	resourceTypes = []string{"another_invalid_resource_type"}
	valid, err = validateResourceType(resourceType, resourceTypes)
	
	if valid {
		t.Fatal("Expected invalid resource type and invalid resource types to be invalid")
	}
	
	if err == nil {
		t.Fatal("Expected error for invalid resource type and invalid resource types but got none")
	}
	
	resourceType = ""
	resourceTypes = []string{}
	valid, err = validateResourceType(resourceType, resourceTypes)
	
	if valid {
		t.Fatal("Expected empty resource type and empty resource types to be invalid")
	}
	
	if err == nil {
		t.Fatal("Expected error for empty resource type and empty resource types but got none")
	}
}

func TestGetResourceNameCoverage2(t *testing.T) {
	resourceType := "azurerm_resource_group"
	separator := "-"
	prefixes := []string{"prefix"}
	name := "testname"
	suffixes := []string{"suffix"}
	randomSuffix := "random"
	convention := ConventionCafClassic
	cleanInput := true
	passthrough := false
	useSlug := true
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	
	if err != nil {
		t.Fatalf("getResourceName returned error for valid parameters: %v", err)
	}
	
	if result == "" {
		t.Fatal("Expected non-empty result for valid parameters")
	}
	
	resourceType = "azurerm_resource_group"
	separator = "-"
	prefixes = []string{"prefix"}
	name = "testname"
	suffixes = []string{"suffix"}
	randomSuffix = "random"
	convention = ConventionCafClassic
	cleanInput = true
	passthrough = true
	useSlug = true
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result, err = getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	
	if err != nil {
		t.Fatalf("getResourceName returned error for passthrough = true: %v", err)
	}
	
	if result != name {
		t.Fatalf("Expected result to be name with passthrough = true, got '%s'", result)
	}
	
	resourceType = "invalid_resource_type"
	separator = "-"
	prefixes = []string{"prefix"}
	name = "testname"
	suffixes = []string{"suffix"}
	randomSuffix = "random"
	convention = ConventionCafClassic
	cleanInput = true
	passthrough = false
	useSlug = true
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	_, err = getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
}

func TestComposeNameCoverage2(t *testing.T) {
	separator := "-"
	prefixes := []string{"prefix1", "prefix2"}
	name := "testname"
	slug := "rg"
	suffixes := []string{"suffix1", "suffix2"}
	randomSuffix := "random"
	maxlength := 63
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result := composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	separator = "-"
	prefixes = []string{}
	name = ""
	slug = ""
	suffixes = []string{}
	randomSuffix = ""
	maxlength = 63
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result != "" {
		t.Fatalf("Expected empty result, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{"prefix1", "prefix2"}
	name = "testname"
	slug = "rg"
	suffixes = []string{"suffix1", "suffix2"}
	randomSuffix = "random"
	maxlength = 0
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result != "" {
		t.Fatalf("Expected empty result with maxlength = 0, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{"prefix1", "prefix2"}
	name = "testname"
	slug = "rg"
	suffixes = []string{"suffix1", "suffix2"}
	randomSuffix = "random"
	maxlength = 5
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if len(result) > maxlength {
		t.Fatalf("Expected length <= %d, got %d", maxlength, len(result))
	}
}

func TestGetSlugCoverage2(t *testing.T) {
	resourceType := "azurerm_resource_group"
	convention := ConventionCafClassic
	
	result := getSlug(resourceType, convention)
	
	if result != "rg" {
		t.Fatalf("Expected 'rg', got '%s'", result)
	}
	
	resourceType = "azurerm_resource_group"
	convention = ConventionCafRandom
	
	result = getSlug(resourceType, convention)
	
	if result != "rg" {
		t.Fatalf("Expected 'rg', got '%s'", result)
	}
	
	resourceType = "azurerm_resource_group"
	convention = ConventionPassThrough
	
	result = getSlug(resourceType, convention)
	
	if result != "" {
		t.Fatalf("Expected empty string, got '%s'", result)
	}
	
	resourceType = "invalid_resource_type"
	convention = ConventionCafClassic
	
	result = getSlug(resourceType, convention)
	
	if result != "" {
		t.Fatalf("Expected empty string, got '%s'", result)
	}
}

func TestGetNameResultCoverage2(t *testing.T) {
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
			"results": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
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
		t.Fatalf("getNameResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{"azurerm_resource_group", "azurerm_storage_account"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}

	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "invalid_resource_type")
	d.Set("resource_types", []interface{}{})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource_type but got none")
	}
}
