package testutils

import (
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
)

type ResourceTestData struct {
	ResourceType    string
	Slug           string
	ValidationRegex string
	MinLength      int
	MaxLength      int
	LowerCase      bool
	Scope          string
}

func GetAllResourceTestData() []*ResourceTestData {
	resources := make([]*ResourceTestData, 0, len(models.ResourceDefinitions))
	
	for resourceType, def := range models.ResourceDefinitions {
		resources = append(resources, &ResourceTestData{
			ResourceType:    resourceType,
			Slug:           def.CafPrefix,
			ValidationRegex: def.ValidationRegExp,
			MinLength:      def.MinLength,
			MaxLength:      def.MaxLength,
			LowerCase:      def.LowerCase,
			Scope:          def.Scope,
		})
	}
	return resources
}

func GetResourceByType(resourceType string) (*ResourceTestData, bool) {
	if def, ok := models.ResourceDefinitions[resourceType]; ok {
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
