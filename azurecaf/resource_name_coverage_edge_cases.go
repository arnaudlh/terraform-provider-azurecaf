package azurecaf

import (
	"testing"
)

func TestGetResourceNameEdgeCasesSpecial(t *testing.T) {
	resource, _ := getResource("azurerm_resource_group")
	originalRegex := resource.ValidationRegExp
	resource.ValidationRegExp = "["  // Invalid regex pattern
	
	_, err := getResourceName("azurerm_resource_group", "-", []string{}, "test", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err == nil {
		t.Fatal("Expected error for invalid regex pattern, got nil")
	}
	
	resource.ValidationRegExp = originalRegex
	
	resource.ValidationRegExp = "^[a-z0-9]+$"  // Only allow lowercase letters and numbers
	
	_, err = getResourceName("azurerm_resource_group", "-", []string{}, "test-name", []string{}, "", "cafrandom", false, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err == nil {
		t.Fatal("Expected error for invalid name, got nil")
	}
	
	resource.ValidationRegExp = originalRegex
}

func TestGetNameResultEdgeCasesSpecial(t *testing.T) {
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
	d.Set("resource_types", []interface{}{"invalid_resource_type"})
	
	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource_type in resource_types, got nil")
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"invalid_resource_type"})
	
	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource in resource_types, got nil")
	}
}

func TestGetResultEdgeCasesSpecial(t *testing.T) {
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
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testnameverylongstring")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("max_length", 25)  // Just enough for the name plus separator
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	if len(result) > 25 {
		t.Fatalf("Expected result length <= 25, got %d", len(result))
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}
