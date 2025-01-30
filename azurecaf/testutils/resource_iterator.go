package testutils

import (
	"sort"
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
	defs := loadResourceDefinitions()
	resources := make([]*ResourceTestData, 0, len(defs))
	
	for resourceType, def := range defs {
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

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].ResourceType < resources[j].ResourceType
	})

	return resources
}

func GetResourceByType(resourceType string) (*ResourceTestData, bool) {
	defs := loadResourceDefinitions()
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
