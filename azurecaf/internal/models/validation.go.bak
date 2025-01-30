package models

import (
	"fmt"
	"regexp"
)

// GetResourceStructure returns a resource structure by its type
func GetResourceStructure(resourceType string) (*ResourceStructure, error) {
	if resourceKey, existing := ResourceMaps[resourceType]; existing {
		resourceType = resourceKey
	}
	if resource, resourceFound := ResourceDefinitions[resourceType]; resourceFound {
		return &resource, nil
	}
	return nil, fmt.Errorf("invalid resource type %s", resourceType)
}

// ValidateResourceType validates that a resource type exists and has valid configuration
func ValidateResourceType(resourceType string) (*ResourceStructure, error) {
	resource, err := GetResourceStructure(resourceType)
	if err != nil {
		return nil, err
	}

	if resource.ValidationRegExp != "" {
		if _, err := regexp.Compile(resource.ValidationRegExp); err != nil {
			return nil, fmt.Errorf("invalid validation regex pattern for resource type %s: %v", resourceType, err)
		}
	}

	return resource, nil
}
