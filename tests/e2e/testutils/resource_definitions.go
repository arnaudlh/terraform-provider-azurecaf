package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type ResourceDefinition struct {
	Name            string `json:"name"`
	CafPrefix       string `json:"caf_prefix"`
	ValidationRegExp string `json:"validation_regexp"`
	MinLength      int    `json:"min_length"`
	MaxLength      int    `json:"max_length"`
	LowerCase      bool   `json:"lowercase"`
	Scope          string `json:"scope"`
}

func GetResourceDefinitions() map[string]*ResourceDefinition {
	defs := make(map[string]*ResourceDefinition)
	
	jsonPath := filepath.Join("..", "..", "resourceDefinition.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &defs); err != nil {
		panic(err)
	}

	// Clean up regex patterns
	for _, def := range defs {
		def.ValidationRegExp = strings.Trim(def.ValidationRegExp, "\"")
	}

	return defs
}

func GetResourceByType(resourceType string) (*ResourceDefinition, bool) {
	defs := GetResourceDefinitions()
	def, ok := defs[resourceType]
	return def, ok
}
