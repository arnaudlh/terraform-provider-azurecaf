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
	defs := GetResourceDefinitions()
	resources := make([]*ResourceTestData, 0, len(defs))
	
	for resourceType, defRaw := range defs {
		def, ok := defRaw.(map[string]interface{})
		if !ok {
			continue
		}
		resources = append(resources, &ResourceTestData{
			ResourceType:    resourceType,
			Slug:           def["slug"].(string),
			ValidationRegex: def["validation_regexp"].(string),
			MinLength:      int(def["min_length"].(float64)),
			MaxLength:      int(def["max_length"].(float64)),
			LowerCase:      def["lowercase"].(bool),
			Scope:          def["scope"].(string),
		})
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].ResourceType < resources[j].ResourceType
	})

	return resources
}

func GetResourceByType(resourceType string) (*ResourceTestData, bool) {
	defs := GetResourceDefinitions()
	if defRaw, ok := defs[resourceType]; ok {
		def, ok := defRaw.(map[string]interface{})
		if !ok {
			return nil, false
		}
		return &ResourceTestData{
			ResourceType:    resourceType,
			Slug:           def["slug"].(string),
			ValidationRegex: def["validation_regexp"].(string),
			MinLength:      int(def["min_length"].(float64)),
			MaxLength:      int(def["max_length"].(float64)),
			LowerCase:      def["lowercase"].(bool),
			Scope:          def["scope"].(string),
		}, true
	}
	return nil, false
}
