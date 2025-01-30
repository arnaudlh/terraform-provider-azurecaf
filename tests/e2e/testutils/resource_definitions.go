package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/testutils"
)

func GetResourceDefinitions() map[string]*models.ResourceStructure {
	defs := make(map[string]*models.ResourceStructure)
	
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

func GetResourceByType(resourceType string) (*testutils.ResourceTestData, bool) {
	return testutils.GetResourceByType(resourceType)
}
