package azurecaf

import (
	"testing"
)

func TestGetResourceNameTargetedCoverage(t *testing.T) {
	resource, _ := getResource("azurerm_resource_group")
	originalRegex := resource.ValidationRegExp
	resource.ValidationRegExp = "["  // Invalid regex pattern
	
	_, err := getResourceName("azurerm_resource_group", "-", []string{}, "test", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err == nil {
		t.Fatal("Expected error for invalid regex pattern, got nil")
	}
	
	resource.ValidationRegExp = originalRegex
}

func TestGetNameResultTargetedCoverage(t *testing.T) {
	provider := Provider()
	
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name resource")
	}
	
	d := nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{"invalid_resource_type"})
	d.Set("random_length", 5)
	
	err := getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource_type in resource_types, got nil")
	}
}

func TestGetResultTargetedCoverage(t *testing.T) {
	provider := Provider()
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "invalid_convention")
	d.Set("resource_type", "rg")
	
	err := getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid convention, got nil")
	}
}
