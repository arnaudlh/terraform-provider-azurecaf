package schemas

import (
    "fmt"
)

type ResourceStructure struct {
    ResourceTypeName  string
    MaxLength        int
    MinLength        int
    ValidationRegExp string
    RegEx           string
    CafPrefix       string
    LowerCase       bool
}

func ValidateResourceType(resourceType string) (*ResourceStructure, error) {
    if resourceType == "" {
        return nil, fmt.Errorf("resource_type parameter is empty")
    }
    if val, ok := ResourceDefinitions[resourceType]; ok {
        return &val, nil
    }
    return nil, fmt.Errorf("resource_type %s is not supported", resourceType)
}

func GetResourceStructure(resourceType string) (*ResourceStructure, error) {
    return ValidateResourceType(resourceType)
}

var ResourceDefinitions = map[string]ResourceStructure{
    "azurerm_resource_group": {
        ResourceTypeName:  "azurerm_resource_group",
        MaxLength:        90,
        MinLength:        1,
        ValidationRegExp: "^[a-zA-Z0-9-_]{1,90}$",
        RegEx:           "[^a-zA-Z0-9-_]",
        CafPrefix:       "rg",
        LowerCase:       true,
    },
    "azurerm_recovery_services_vault": {
        ResourceTypeName:  "azurerm_recovery_services_vault",
        MaxLength:        16,
        MinLength:        2,
        ValidationRegExp: "^[a-zA-Z][a-zA-Z0-9]{1,15}$",
        RegEx:           "[^a-zA-Z0-9]",
        CafPrefix:       "rsv",
        LowerCase:       true,
    },
    "azurerm_container_registry": {
        ResourceTypeName:  "azurerm_container_registry",
        MaxLength:        63,
        MinLength:        5,
        ValidationRegExp: "^[a-zA-Z0-9]{5,63}$",
        RegEx:           "[^a-zA-Z0-9]",
        CafPrefix:       "cr",
        LowerCase:       true,
    },
    "azurerm_container_app": {
        ResourceTypeName:  "azurerm_container_app",
        MaxLength:        27,
        MinLength:        1,
        ValidationRegExp: "^[a-zA-Z0-9][a-zA-Z0-9-]{0,26}$",
        RegEx:           "[^a-zA-Z0-9-]",
        CafPrefix:       "ca",
        LowerCase:       true,
    },
    "azurerm_container_app_environment": {
        ResourceTypeName:  "azurerm_container_app_environment",
        MaxLength:        25,
        MinLength:        1,
        ValidationRegExp: "^[a-zA-Z0-9][a-zA-Z0-9-]{0,24}$",
        RegEx:           "[^a-zA-Z0-9-]",
        CafPrefix:       "cae",
        LowerCase:       true,
    },
}
