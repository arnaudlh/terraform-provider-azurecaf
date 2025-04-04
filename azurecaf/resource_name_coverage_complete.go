package azurecaf

import (
	"testing"
)

func TestGetResourceNameCompleteCoverage(t *testing.T) {
	resource, _ := getResource("azurerm_resource_group")
	originalRegex := resource.ValidationRegExp
	resource.ValidationRegExp = "["  // Invalid regex pattern
	
	_, err := getResourceName("azurerm_resource_group", "-", []string{}, "test", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err == nil {
		t.Fatal("Expected error for invalid regex pattern, got nil")
	}
	
	resource.ValidationRegExp = originalRegex
	
	resource.LowerCase = true
	result, err := getResourceName("azurerm_resource_group", "-", []string{}, "TestName", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with lowercase conversion returned error: %v", err)
	}
	if result != "rg-testname" {
		t.Fatalf("Expected lowercase result 'rg-testname', got '%s'", result)
	}
	
	resource.ValidationRegExp = "^[a-z0-9]+$"  // Only allow lowercase letters and numbers
	_, err = getResourceName("azurerm_resource_group", "-", []string{}, "test-name", []string{}, "", "cafrandom", false, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err == nil {
		t.Fatal("Expected error for invalid name, got nil")
	}
}

func TestGetNameResultCompleteCoverage(t *testing.T) {
	provider := Provider()
	
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name resource")
	}
	
	d := nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{})
	
	err := getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for empty resource_type and resource_types, got nil")
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{"invalid_resource_type1", "invalid_resource_type2"})
	
	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource_types, got nil")
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"invalid_resource_type"})
	
	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource in resource_types, got nil")
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{"azurerm_resource_group", "azurerm_virtual_network"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}
	
	results := d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}

func TestGetResultCompleteCoverage(t *testing.T) {
	provider := Provider()
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "cafrandom")
	d.Set("resource_type", "")
	
	err := getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for empty resource_type, got nil")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "cafrandom")
	d.Set("resource_type", "invalid_resource_type")
	
	err = getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource_type, got nil")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "cafclassic")
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "cafrandom")
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "invalid_convention")
	d.Set("resource_type", "rg")
	
	err = getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid convention, got nil")
	}
}

func TestComposeNameCompleteCoverage(t *testing.T) {
	result := composeName("-", 
		[]string{"prefix1", "prefix2"}, 
		"name", 
		"slug", 
		[]string{"suffix1", "suffix2"}, 
		"random", 
		100, 
		[]string{"name", "slug", "random", "suffixes", "prefixes"})
	
	if result != "slug-name-random-suffix1-suffix2-prefix1-prefix2" {
		t.Fatalf("Expected 'slug-name-random-suffix1-suffix2-prefix1-prefix2', got '%s'", result)
	}
	
	result = composeName("-", 
		[]string{"prefix1", "prefix2"}, 
		"name", 
		"slug", 
		[]string{"suffix1", "suffix2"}, 
		"random", 
		10, 
		[]string{"name", "slug", "random", "suffixes", "prefixes"})
	
	if len(result) > 10 {
		t.Fatalf("Expected result length <= 10, got %d", len(result))
	}
	
	result = composeName("-", 
		[]string{}, 
		"", 
		"", 
		[]string{}, 
		"", 
		100, 
		[]string{"name", "slug", "random", "suffixes", "prefixes"})
	
	if result != "" {
		t.Fatalf("Expected empty string, got '%s'", result)
	}
}

func TestValidateResourceTypeCompleteCoverage(t *testing.T) {
	valid, err := validateResourceType("", []string{})
	if valid || err == nil {
		t.Fatal("Expected error for empty resource_type and resource_types, got nil")
	}
	
	valid, err = validateResourceType("azurerm_resource_group", []string{})
	if !valid || err != nil {
		t.Fatalf("validateResourceType returned error: %v", err)
	}
	
	valid, err = validateResourceType("invalid_resource_type", []string{})
	if valid || err == nil {
		t.Fatal("Expected error for invalid resource_type, got nil")
	}
	
	valid, err = validateResourceType("", []string{"azurerm_resource_group", "azurerm_virtual_network"})
	if !valid || err != nil {
		t.Fatalf("validateResourceType returned error: %v", err)
	}
	
	valid, err = validateResourceType("", []string{"invalid_resource_type1", "invalid_resource_type2"})
	if valid || err == nil {
		t.Fatal("Expected error for invalid resource_types, got nil")
	}
	
	valid, err = validateResourceType("azurerm_resource_group", []string{"azurerm_virtual_network", "azurerm_subnet"})
	if !valid || err != nil {
		t.Fatalf("validateResourceType returned error: %v", err)
	}
	
	valid, err = validateResourceType("azurerm_resource_group", []string{"invalid_resource_type"})
	if valid || err == nil {
		t.Fatal("Expected error for invalid resource_types, got nil")
	}
}
