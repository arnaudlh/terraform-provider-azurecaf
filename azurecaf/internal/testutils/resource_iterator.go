package testutils

import (
	"sort"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
)

func GetAllResourceTestData() []*models.ResourceStructure {
	resources := make([]*models.ResourceStructure, 0, len(models.ResourceDefinitions))
	
	for resourceType, def := range models.ResourceDefinitions {
		resources = append(resources, &models.ResourceStructure{
			Name:            resourceType,
			CafPrefix:       def.CafPrefix,
			ValidationRegExp: def.ValidationRegExp,
			MinLength:      def.MinLength,
			MaxLength:      def.MaxLength,
			LowerCase:      def.LowerCase,
			Scope:          def.Scope,
		})
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	return resources
}

func GetResourceByType(resourceType string) (*models.ResourceStructure, bool) {
	def, ok := models.ResourceDefinitions[resourceType]
	return def, ok
}
