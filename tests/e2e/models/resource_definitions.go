package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ResourceDefinition struct {
	CafPrefix       string `json:"caf_prefix"`
	ValidationRegExp string `json:"validation_regexp"`
	MinLength      int    `json:"min_length"`
	MaxLength      int    `json:"max_length"`
	LowerCase      bool   `json:"lowercase"`
	Scope          string `json:"scope"`
}

var resourceDefinitions map[string]ResourceDefinition

func GetResourceDefinitions() map[string]ResourceDefinition {
	if resourceDefinitions != nil {
		return resourceDefinitions
	}

	jsonPath := filepath.Join("..", "..", "resourceDefinition.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &resourceDefinitions); err != nil {
		panic(err)
	}

	return resourceDefinitions
}
