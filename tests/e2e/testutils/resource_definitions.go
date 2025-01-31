package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ResourceDefinition struct {
	ResourceTypeName  string `json:"name"`
	CafPrefix        string `json:"slug,omitempty"`
	MinLength        int    `json:"min_length"`
	MaxLength        int    `json:"max_length"`
	LowerCase        bool   `json:"lowercase,omitempty"`
	RegEx            string `json:"regex,omitempty"`
	ValidationRegExp string `json:"validation_regexp,omitempty"`
	Scope           string `json:"scope,omitempty"`
}

func GetResourceDefinitions() map[string]interface{} {
	jsonPath := filepath.Join("..", "..", "resourceDefinition.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}

	var definitionsArray []map[string]interface{}
	if err := json.Unmarshal(data, &definitionsArray); err != nil {
		panic(err)
	}

	// Convert array to map using resource name as key
	definitions := make(map[string]interface{})
	for _, def := range definitionsArray {
		if name, ok := def["name"].(string); ok {
			definitions[name] = def
		}
	}

	return definitions
}
