package models

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
)

type ResourceStructure struct {
    ResourceTypeName  string `json:"name"`
    CafPrefix        string `json:"slug,omitempty"`
    MinLength        int    `json:"min_length"`
    MaxLength        int    `json:"max_length"`
    LowerCase        bool   `json:"lowercase,omitempty"`
    RegEx           string `json:"regex,omitempty"`
    ValidationRegExp string `json:"validation_regex,omitempty"`
    Dashes          bool   `json:"dashes"`
    Scope           string `json:"scope,omitempty"`
}

var ResourceDefinitions = map[string]*ResourceStructure{}
var ResourceMaps = map[string]string{}

func init() {
    // Load resource definitions from JSON file
    searchPaths := []string{
        "/home/ubuntu/repos/terraform-provider-azurecaf/resourceDefinition.json",
        "/home/ubuntu/repos/terraform-provider-azurecaf/tests/e2e/resourceDefinition.json",
    }

    var jsonPath string
    for _, path := range searchPaths {
        if _, err := os.Stat(path); err == nil {
            jsonPath = path
            log.Printf("[INFO] Found resource definition at: %s", path)
            break
        }
    }

    if jsonPath == "" {
        log.Printf("[WARN] resourceDefinition.json not found in search paths")
        return
    }

    data, err := os.ReadFile(jsonPath)
    if err != nil {
        log.Printf("[ERROR] Failed to read resource definitions: %v", err)
        return
    }

    var definitionsArray []ResourceStructure
    if err := json.Unmarshal(data, &definitionsArray); err != nil {
        log.Printf("[ERROR] Failed to parse resource definitions: %v", err)
        return
    }

    // Initialize ResourceDefinitions map
    for i := range definitionsArray {
        def := &definitionsArray[i]
        if def.ValidationRegExp != "" {
            // Remove surrounding quotes and unescape internal quotes
            def.ValidationRegExp = strings.Trim(def.ValidationRegExp, "\"")
            def.ValidationRegExp = strings.ReplaceAll(def.ValidationRegExp, "\\\"", "\"")
            // Remove any remaining escaped backslashes
            def.ValidationRegExp = strings.ReplaceAll(def.ValidationRegExp, "\\\\", "\\")
        }
        ResourceDefinitions[def.ResourceTypeName] = def
        ResourceMaps[def.ResourceTypeName] = def.ResourceTypeName
    }
    log.Printf("[DEBUG] Loaded %d resource definitions", len(ResourceDefinitions))
}

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
