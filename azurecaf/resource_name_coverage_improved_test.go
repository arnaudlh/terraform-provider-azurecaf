package azurecaf

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestGetResourceNameImproved tests the getResourceName function with various inputs
func TestGetResourceNameImproved(t *testing.T) {
	// Test with standard parameters
	resourceTypeName := "azurerm_resource_group"
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

	result, err := getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned unexpected error: %v", err)
	}
	
	// Verify result contains expected components
	if !strings.Contains(result, "rg") {
		t.Fatalf("Expected result to contain resource group slug 'rg', got %s", result)
	}
	
	// Test with invalid resource type
	_, err = getResourceName("invalid_resource_type", separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
	
	// Test with passthrough enabled
	passthrough = true
	result, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName with passthrough returned error: %v", err)
	}
	if result != name {
		t.Fatalf("With passthrough, expected result to be %s, got %s", name, result)
	}
	
	// Test with special characters and clean_input disabled
	passthrough = false
	cleanInput = false
	name = "test!@#name"
	_, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err == nil {
		t.Fatal("Expected error for special characters with clean_input disabled but got none")
	}
	
	// Test with useSlug disabled
	cleanInput = true
	useSlug = false
	result, err = getResourceName(resourceTypeName, separator, prefixes, "testname", suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName with useSlug=false returned error: %v", err)
	}
	if strings.Contains(result, "rg") {
		t.Fatalf("With useSlug=false, expected result not to contain slug 'rg', got %s", result)
	}
}

// TestGetNameResultImproved tests the getNameResult function with various inputs
func TestGetNameResultImproved(t *testing.T) {
	// Create a schema resource for testing
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

	// Test with invalid resource type
	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "invalid_resource_type")
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}

	// Test with multiple resource types
	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{"azurerm_resource_group", "azurerm_virtual_machine"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult with multiple resource types returned error: %v", err)
	}

	results := d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	
	// Test with single resource type
	d = r.TestResourceData()
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

	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult with single resource type returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	// Test with no resource type specified
	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for no resource type but got none")
	}
}

// TestGetResultImproved tests the getResult function with various inputs
func TestGetResultImproved(t *testing.T) {
	// Create a schema resource for testing
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

	// Test with passthrough convention
	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionPassThrough)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with passthrough convention returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result != "prefix-testname-postfix" {
		t.Fatalf("With passthrough convention, expected result to be prefix-testname-postfix, got %s", result)
	}

	// Test with CAF random convention
	d = r.TestResourceData()
	d.Set("name", "")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with CAF random convention returned error: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result with CAF random convention")
	}
	
	// Test with max length constraint
	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 10) // Less than default max length

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with max length constraint returned error: %v", err)
	}

	result = d.Get("result").(string)
	if len(result) > 10 {
		t.Fatalf("Expected result length <= 10, got %d", len(result))
	}
	
	// Test with invalid resource type
	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "invalid_resource_type")
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
}

// TestHelperFunctionsImproved tests the helper functions
func TestHelperFunctionsImproved(t *testing.T) {
	// Test cleanString
	resource, _ := getResource("azurerm_resource_group")
	result := cleanString("test!@#name", resource)
	if result != "testname" {
		t.Fatalf("cleanString expected to remove special characters, got %s", result)
	}
	
	// Test cleanSlice
	slice := []string{"test!@#1", "test!@#2"}
	cleanedSlice := cleanSlice(slice, resource)
	if cleanedSlice[0] != "test1" || cleanedSlice[1] != "test2" {
		t.Fatalf("cleanSlice expected to remove special characters, got %v", cleanedSlice)
	}
	
	// Test concatenateParameters
	result = concatenateParameters("-", []string{"prefix"}, []string{"name"}, []string{"suffix"})
	if result != "prefix-name-suffix" {
		t.Fatalf("concatenateParameters expected to join with separator, got %s", result)
	}
	
	// Test getResource
	resource, err := getResource("azurerm_resource_group")
	if err != nil {
		t.Fatalf("getResource returned error: %v", err)
	}
	if resource.ResourceTypeName != "azurerm_resource_group" {
		t.Fatalf("getResource expected to return resource with correct name, got %s", resource.ResourceTypeName)
	}
	
	// Test getSlug
	slug := getSlug("azurerm_resource_group", ConventionCafClassic)
	if slug != "rg" {
		t.Fatalf("getSlug expected to return 'rg', got %s", slug)
	}
	
	// Test trimResourceName
	result = trimResourceName("abcdefghijklmnopqrstuvwxyz", 10)
	if result != "abcdefghij" {
		t.Fatalf("trimResourceName expected to trim to 10 characters, got %s", result)
	}
	
	// Test convertInterfaceToString
	interfaceSlice := []interface{}{"test1", "test2"}
	stringSlice := convertInterfaceToString(interfaceSlice)
	if stringSlice[0] != "test1" || stringSlice[1] != "test2" {
		t.Fatalf("convertInterfaceToString expected to convert to string slice, got %v", stringSlice)
	}
	
	// Test validateResourceType
	valid, err := validateResourceType("azurerm_resource_group", []string{})
	if !valid || err != nil {
		t.Fatalf("validateResourceType expected to validate single resource type, got valid=%v, err=%v", valid, err)
	}
	
	valid, err = validateResourceType("", []string{"azurerm_resource_group", "azurerm_virtual_machine"})
	if !valid || err != nil {
		t.Fatalf("validateResourceType expected to validate multiple resource types, got valid=%v, err=%v", valid, err)
	}
	
	valid, err = validateResourceType("", []string{})
	if valid || err == nil {
		t.Fatal("validateResourceType expected to fail with no resource types")
	}
}

// TestComposeNameImproved tests the composeName function
func TestComposeNameImproved(t *testing.T) {
	// Test with all components
	result := composeName("-", []string{"prefix"}, "name", "slug", []string{"suffix"}, "random", 100, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if !strings.Contains(result, "name") || !strings.Contains(result, "slug") || !strings.Contains(result, "random") || !strings.Contains(result, "prefix") || !strings.Contains(result, "suffix") {
		t.Fatalf("composeName expected to include all components, got %s", result)
	}
	
	// Test with max length constraint
	result = composeName("-", []string{"prefix"}, "name", "slug", []string{"suffix"}, "random", 10, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if len(result) > 10 {
		t.Fatalf("composeName expected to respect max length, got %s with length %d", result, len(result))
	}
	
	// Test with empty components
	result = composeName("-", []string{}, "", "", []string{}, "", 100, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if result != "" {
		t.Fatalf("composeName expected to return empty string with empty components, got %s", result)
	}
	
	// Test with multiple prefixes and suffixes
	result = composeName("-", []string{"prefix1", "prefix2"}, "name", "slug", []string{"suffix1", "suffix2"}, "random", 100, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if !strings.Contains(result, "prefix1") || !strings.Contains(result, "prefix2") || !strings.Contains(result, "suffix1") || !strings.Contains(result, "suffix2") {
		t.Fatalf("composeName expected to include multiple prefixes and suffixes, got %s", result)
	}
}
