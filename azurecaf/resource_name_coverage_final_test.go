package azurecaf

import (
	"testing"
)

func TestGetResourceNameFinalCoverageExtra(t *testing.T) {
	result, err := getResourceName("azurerm_resource_group", "-", []string{}, "test-name", []string{}, "", "cafrandom", true, true, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with passthrough returned error: %v", err)
	}
	if result != "test-name" {
		t.Fatalf("Expected 'test-name', got '%s'", result)
	}

	_, err = getResourceName("invalid_resource_type", "-", []string{}, "test", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}

	result, err = getResourceName("azurerm_storage_account", "-", []string{}, "TestName", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with lowercase conversion returned error: %v", err)
	}
	if result != "stestame" {
		t.Fatalf("Expected 'stestame', got '%s'", result)
	}

	_, err = getResourceName("azurerm_storage_account", "-", []string{}, "test_name_with_invalid_chars!", []string{}, "", "cafrandom", false, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err == nil {
		t.Fatal("Expected error for invalid name, got nil")
	}

	result, err = getResourceName("azurerm_resource_group", "-", []string{}, "test", []string{}, "", "cafrandom", true, false, false, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with useSlug=false returned error: %v", err)
	}
	if result != "test" {
		t.Fatalf("Expected 'test', got '%s'", result)
	}
}

func TestGetNameResultFinalCoverageExtra(t *testing.T) {
	provider := Provider()
	
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name resource")
	}
	
	d := nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{"azurerm_resource_group", "azurerm_virtual_network"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err := getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}
	
	results := d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "invalid_resource_type")
	
	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource_type, got nil")
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"azurerm_virtual_network", "azurerm_subnet"})
	d.Set("random_length", 5)
	
	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	results = d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}

func TestGetResultFinalCoverageExtra(t *testing.T) {
	provider := Provider()
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "random")
	d.Set("resource_type", "rg")
	
	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "passthrough")
	d.Set("resource_type", "rg")
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result != "testname" {
		t.Fatalf("Expected result 'testname', got %s", result)
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
}
