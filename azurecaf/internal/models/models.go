package models

import (
	"fmt"
	"regexp"
)

func init() {
	// Initialize ResourceMaps after ResourceDefinitions is populated
	for k, v := range ResourceDefinitions {
		ResourceMaps[k] = k
		ResourceMaps[v.ResourceTypeName] = k
	}
}

// ValidateResourceType validates if a resource type exists
func ValidateResourceType(resourceType string) (bool, error) {
	if _, exists := ResourceDefinitions[resourceType]; !exists {
		return false, fmt.Errorf("invalid resource type: %s", resourceType)
	}
	return true, nil
}

// GetResourceStructure returns the resource structure for a given type
func GetResourceStructure(resourceType string) (*ResourceStructure, error) {
	if resource, exists := ResourceDefinitions[resourceType]; exists {
		return &resource, nil
	}
	return nil, fmt.Errorf("resource type not found: %s", resourceType)
}

// ValidateLength checks if a string meets length requirements
func ValidateLength(input string, minLength, maxLength int) error {
	length := len(input)
	if length < minLength || length > maxLength {
		return fmt.Errorf("length must be between %d and %d, got %d", minLength, maxLength, length)
	}
	return nil
}

// ValidateRegex checks if a string matches a regex pattern
func ValidateRegex(input, pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %s", err.Error())
	}
	if !regex.MatchString(input) {
		return fmt.Errorf("input does not match pattern %s", pattern)
	}
	return nil
}

// ResourceStructure stores the CafPrefix and the MaxLength of an azure resource
type ResourceStructure struct {
	// Resource type name
	ResourceTypeName string `json:"name"`
	// Resource prefix as defined in the Azure Cloud Adoption Framework
	CafPrefix string `json:"slug,omitempty"`
	// MaxLength attribute define the maximum length of the name
	MinLength int `json:"min_length"`
	// MaxLength attribute define the maximum length of the name
	MaxLength int `json:"max_length"`
	// enforce lowercase
	LowerCase bool `json:"lowercase,omitempty"`
	// Regular expression to apply to the resource type
	RegEx string `json:"regex,omitempty"`
	// the Regular expression to validate the generated string
	ValidationRegExp string `json:"validatation_regex,omitempty"`
	// can the resource include dashes
	Dashes bool `json:"dashes"`
	// The scope of this name where it needs to be unique
	Scope string `json:"scope,omitempty"`
}
