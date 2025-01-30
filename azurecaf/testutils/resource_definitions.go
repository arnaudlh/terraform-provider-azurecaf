package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

type ResourceDefinition struct {
	ResourceTypeName string `json:"name"`
	CafPrefix       string `json:"slug,omitempty"`
	MinLength      int    `json:"min_length"`
	MaxLength      int    `json:"max_length"`
	LowerCase      bool   `json:"lowercase,omitempty"`
	RegEx          string `json:"regex,omitempty"`
	ValidationRegExp string `json:"validation_regexp,omitempty"`
	Scope          string `json:"scope,omitempty"`
}

// GetResourceDefinitions loads resource definitions directly from JSON
func GetResourceDefinitions() map[string]ResourceDefinition {
	jsonPath := filepath.Join("..", "..", "resourceDefinition.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}

	var definitions map[string]ResourceDefinition
	if err := json.Unmarshal(data, &definitions); err != nil {
		panic(err)
	}

	return definitions
}

// GetResourceByType returns a specific resource definition by its type
func GetResourceByType(resourceType string) (*ResourceTestData, bool) {
	defs := GetResourceDefinitions()
	if def, ok := defs[resourceType]; ok {
		return &ResourceTestData{
			ResourceType:    resourceType,
			Slug:           def.CafPrefix,
			ValidationRegex: def.ValidationRegExp,
			MinLength:      def.MinLength,
			MaxLength:      def.MaxLength,
			LowerCase:      def.LowerCase,
			Scope:          def.Scope,
		}, true
	}
	return nil, false
}
