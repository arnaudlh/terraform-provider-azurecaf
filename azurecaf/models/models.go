package models

import (
    "fmt"
    "regexp"
)

type ResourceStructure struct {
    ResourceTypeName  string `json:"name"`
    CafPrefix        string `json:"slug,omitempty"`
    MinLength        int    `json:"min_length"`
    MaxLength        int    `json:"max_length"`
    LowerCase        bool   `json:"lowercase,omitempty"`
    RegEx           string `json:"regex,omitempty"`
    ValidationRegExp string `json:"validation_regexp,omitempty"`
    Dashes          bool   `json:"dashes"`
    Scope           string `json:"scope,omitempty"`
}

var ResourceDefinitions = map[string]*ResourceStructure{}
var ResourceMaps = map[string]string{}

func GetResourceStructure(resourceType string) (*ResourceStructure, error) {
    if resourceKey, existing := ResourceMaps[resourceType]; existing {
        resourceType = resourceKey
    }
    if resource, resourceFound := ResourceDefinitions[resourceType]; resourceFound {
        return resource, nil
    }
    return nil, fmt.Errorf("invalid resource type %s", resourceType)
}

func ValidateResourceType(resourceType string) (*ResourceStructure, error) {
    resource, err := GetResourceStructure(resourceType)
    if err != nil {
        return nil, err
    }

    if resource.ValidationRegExp != "" {
        if _, err := regexp.Compile(resource.ValidationRegExp); err != nil {
            return nil, fmt.Errorf("invalid validation regex pattern for resource type %s: %v", resourceType, err)
        }
    }

    return resource, nil
}
