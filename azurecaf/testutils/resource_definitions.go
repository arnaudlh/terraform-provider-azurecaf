package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func loadResourceDefinitions() map[string]ResourceDefinition {
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
