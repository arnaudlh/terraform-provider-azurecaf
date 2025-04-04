package azurecaf

import (
	"context"
	"testing"
)

func TestConcatenateParameters(t *testing.T) {
	result := concatenateParameters("-", []string{})
	if result != "" {
		t.Fatalf("Expected empty string, got %s", result)
	}

	result = concatenateParameters("-", []string{"test"})
	if result != "test" {
		t.Fatalf("Expected 'test', got %s", result)
	}

	result = concatenateParameters("-", []string{"test1"}, []string{"test2"}, []string{"test3"})
	if result != "test1-test2-test3" {
		t.Fatalf("Expected 'test1-test2-test3', got %s", result)
	}
}

func TestGetResource(t *testing.T) {
	resource, err := getResource("azurerm_resource_group")
	if err != nil {
		t.Fatalf("getResource returned error: %v", err)
	}
	if resource.ResourceTypeName != "azurerm_resource_group" {
		t.Fatalf("Expected resource type 'azurerm_resource_group', got %s", resource.ResourceTypeName)
	}

	_, err = getResource("invalid_resource_type")
	if err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}
}

func TestTrimResourceName(t *testing.T) {
	result := trimResourceName("test", 10)
	if result != "test" {
		t.Fatalf("Expected 'test', got %s", result)
	}

	result = trimResourceName("testtesttest", 5)
	if result != "testt" {
		t.Fatalf("Expected 'testt', got %s", result)
	}
}

func TestValidateResourceType(t *testing.T) {
	valid, err := validateResourceType("azurerm_resource_group", []string{})
	if !valid || err != nil {
		t.Fatalf("validateResourceType returned error: %v", err)
	}

	valid, err = validateResourceType("invalid_resource_type", []string{})
	if valid || err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}
}

func TestGetResourceName(t *testing.T) {
	prefixes := []string{"prefix"}
	suffixes := []string{"suffix"}
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	name, err := getResourceName("azurerm_resource_group", "-", prefixes, "test", suffixes, "random", "cafrandom", true, false, true, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned error: %v", err)
	}
	if name == "" {
		t.Fatal("Expected non-empty name, got empty string")
	}

	_, err = getResourceName("invalid_resource_type", "-", prefixes, "test", suffixes, "random", "cafrandom", true, false, true, namePrecedence)
	if err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}
}

func TestResourceNameStateUpgradeV2(t *testing.T) {
	state := map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_resource_group",
		"prefixes":      []interface{}{"prefix"},
		"suffixes":      []interface{}{"suffix"},
	}

	newState, err := resourceNameStateUpgradeV2(context.Background(), state, nil)
	if err != nil {
		t.Fatalf("resourceNameStateUpgradeV2 returned error: %v", err)
	}

	if newState["name"] != "test" {
		t.Fatalf("Expected name 'test', got %s", newState["name"])
	}
	if newState["resource_type"] != "azurerm_resource_group" {
		t.Fatalf("Expected resource_type 'azurerm_resource_group', got %s", newState["resource_type"])
	}
}

func TestGetResult(t *testing.T) {
	t.Skip("Skipping test that requires more complex ResourceData setup")
}
