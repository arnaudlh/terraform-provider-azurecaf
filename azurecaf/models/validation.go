package models

import (
    "fmt"
    "regexp"
)

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
