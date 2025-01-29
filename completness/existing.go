package completness

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// The idea of this package it is to check for package completness
// To update the list of existing resources I did query
// https://registry.terraform.io/v2/provider-versions/7185?include=provider-docs
// them use the jq espression `"azurerm_\(.included[].attributes.title)"`
// followed by manual cleaning of the non resources doc links

// ResourceStructure resource definition structure
// Copied from gen.go
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
	ValidationRegExp string `json:"validation_regex,omitempty"`
	// can the resource include dashes
	Dashes bool `json:"dashes"`
	// The scope of this name where it needs to be unique
	Scope string `json:"scope,omitempty"`
}

func ValidateResourceDefinition(resources []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	sourceDefinitions, err := os.ReadFile(path.Join(wd, "resourceDefinition.json"))
	if err != nil {
		return fmt.Errorf("failed to read resource definition file: %v", err)
	}
	var data []ResourceStructure
	err = json.Unmarshal(sourceDefinitions, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal resource definitions: %v", err)
	}
	for _, name := range resources {
		if _, found := findByName(data, name); !found {
			return fmt.Errorf("resource type %s not found in the resource definition file", name)
		}
	}
	return nil
}

func GetResourceDefinition() ([]ResourceStructure, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}
	sourceDefinitions, err := os.ReadFile(path.Join(wd, "../resourceDefinition.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read resource definition file: %v", err)
	}
	var result []ResourceStructure
	err = json.Unmarshal(sourceDefinitions, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource definitions: %v", err)
	}
	return result, nil
}

func GetResourceMap() (map[string]ResourceStructure, error) {
	resourceDefs, err := GetResourceDefinition()
	if err != nil {
		return nil, err
	}
	result := make(map[string]ResourceStructure)
	for _, def := range resourceDefs {
		result[def.ResourceTypeName] = def
	}
	return result, nil
}



func findByName(slice []ResourceStructure, name string) (int, bool) {
	for i, item := range slice {
		if item.ResourceTypeName == name {
			return i, true
		}
	}
	return -1, false
}
