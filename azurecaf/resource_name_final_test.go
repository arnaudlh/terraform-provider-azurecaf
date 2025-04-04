package azurecaf

import (
	"testing"
)

func TestGetResourceNameFinalTest(t *testing.T) {
	resource, _ := getResource("azurerm_resource_group")
	resource.LowerCase = true

	result, err := getResourceName("azurerm_resource_group", "-", []string{}, "TestName", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with lowercase conversion returned error: %v", err)
	}
	if result != "rg-TestName" {
		t.Fatalf("Expected result 'rg-TestName', got '%s'", result)
	}

	result, err = getResourceName("azurerm_resource_group", "-", []string{}, "test-name", []string{}, "", "cafrandom", true, true, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with passthrough returned error: %v", err)
	}
	if result != "test-name" {
		t.Fatalf("Expected 'test-name', got '%s'", result)
	}
}

func TestGetNameResultFinalTest(t *testing.T) {
	provider := Provider()

	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name resource")
	}

	d := nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"azurerm_virtual_network", "azurerm_subnet"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)

	err := getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	results := d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}

func TestGetResultFinalTest(t *testing.T) {
	provider := Provider()

	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}

	d := conventionResource.TestResourceData()
	d.Set("name", "")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)

	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	d = conventionResource.TestResourceData()
	d.Set("name", "testnameverylongstring")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("max_length", 25) // Just enough for the name plus separator

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	if len(result) > 25 {
		t.Fatalf("Expected result length <= 25, got %d", len(result))
	}
}
