package azurecaf

import (
	"testing"
)

func TestGetResource_AllCases(t *testing.T) {
	resourceType := "azurerm_resource_group"
	resource, err := getResource(resourceType)
	
	if err != nil {
		t.Fatalf("getResource returned error for valid resource: %v", err)
	}
	
	if resource == nil {
		t.Fatal("Expected resource to be returned for valid resource")
	}
	
	resourceType = "rg"
	resource, err = getResource(resourceType)
	
	if err != nil {
		t.Fatalf("getResource returned error for valid mapped resource: %v", err)
	}
	
	if resource == nil {
		t.Fatal("Expected resource to be returned for valid mapped resource")
	}
}

func TestTrimResourceName_AllCases(t *testing.T) {
	resourceName := "short-name"
	maxLength := 20
	
	result := trimResourceName(resourceName, maxLength)
	
	if result != resourceName {
		t.Fatalf("Expected %s, got %s", resourceName, result)
	}
	
	resourceName = "this-is-a-very-long-resource-name-that-needs-to-be-trimmed"
	maxLength = 20
	
	result = trimResourceName(resourceName, maxLength)
	
	if len(result) != maxLength {
		t.Fatalf("Expected length %d, got %d", maxLength, len(result))
	}
	
	if result != resourceName[:maxLength] {
		t.Fatalf("Expected %s, got %s", resourceName[:maxLength], result)
	}
	
	resourceName = "test"
	maxLength = 0
	
	result = trimResourceName(resourceName, maxLength)
	
	if result != "" {
		t.Fatalf("Expected empty string, got %s", result)
	}
	
	resourceName = ""
	maxLength = 10
	
	result = trimResourceName(resourceName, maxLength)
	
	if result != "" {
		t.Fatalf("Expected empty string, got %s", result)
	}
}

func TestValidateResourceType_AllCases(t *testing.T) {
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
