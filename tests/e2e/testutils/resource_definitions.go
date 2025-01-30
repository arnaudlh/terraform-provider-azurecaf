package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ResourceDefinition struct {
	Name            string `json:"name"`
	MinLength       int    `json:"min_length"`
	MaxLength       int    `json:"max_length"`
	ValidationRegex string `json:"validation_regex"`
	Scope          string `json:"scope"`
	Slug           string `json:"slug"`
	Dashes         bool   `json:"dashes"`
	LowerCase      bool   `json:"lowercase"`
	Regex          string `json:"regex"`
}

var resourceDefinitions []ResourceDefinition

func GetResourceDefinitions() map[string]ResourceDefinition {
	if len(resourceDefinitions) > 0 {
		result := make(map[string]ResourceDefinition)
		for _, def := range resourceDefinitions {
			result[def.Name] = def
		}
		return result
	}

	jsonPath := filepath.Join("..", "..", "resourceDefinition.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &resourceDefinitions); err != nil {
		panic(err)
	}

	result := make(map[string]ResourceDefinition)
	for _, def := range resourceDefinitions {
		// Clean up regex patterns by removing quotes
		def.ValidationRegex = strings.Trim(def.ValidationRegex, "\"")
		def.Regex = strings.Trim(def.Regex, "\"")
		result[def.Name] = def
	}

	return result
}

func GetResourceByType(resourceType string) (*ResourceTestData, bool) {
	if def, ok := GetResourceDefinitions()[resourceType]; ok {
		return &ResourceTestData{
			ResourceType:    resourceType,
			MinLength:      def.MinLength,
			MaxLength:      def.MaxLength,
			ValidationRegex: def.ValidationRegex,
			Scope:          def.Scope,
			Slug:           def.Slug,
			Dashes:         def.Dashes,
			LowerCase:      def.LowerCase,
			Regex:          def.Regex,
		}, true
	}
	return nil, false
}
