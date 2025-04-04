package azurecaf

import (
	"testing"
)

func TestGetResourceEdgeCases(t *testing.T) {
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

func TestTrimResourceNameEdgeCases(t *testing.T) {
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

func TestValidateResourceTypeEdgeCases(t *testing.T) {
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

func TestGetResourceNameEdgeCases(t *testing.T) {
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
