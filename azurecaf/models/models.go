package models

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
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
    ValidationRegExp string `json:"validation_regexp,omitempty"`
    Dashes          bool   `json:"dashes"`
    Scope           string `json:"scope,omitempty"`
}

var ResourceDefinitions = map[string]*ResourceStructure{}
var ResourceMaps = map[string]string{}

func init() {
    // Load resource definitions from JSON file
    wd, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("failed to get working directory: %v", err))
    }

    // Try to find resourceDefinition.json by walking up directories
    var jsonPath string
    for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
        path := filepath.Join(dir, "resourceDefinition.json")
        if _, err := os.Stat(path); err == nil {
            jsonPath = path
            break
        }
    }

    if jsonPath == "" {
        panic("resourceDefinition.json not found in any parent directory")
    }

    data, err := os.ReadFile(jsonPath)
    if err != nil {
        panic(fmt.Sprintf("failed to read resource definitions: %v", err))
    }

    var definitionsArray []ResourceStructure
    if err := json.Unmarshal(data, &definitionsArray); err != nil {
        panic(fmt.Sprintf("failed to parse resource definitions: %v", err))
    }

    // Clean up validation regex patterns by removing quotes
    for i := range definitionsArray {
        if definitionsArray[i].ValidationRegExp != "" {
            // Remove surrounding quotes and unescape inner quotes
            pattern := definitionsArray[i].ValidationRegExp
            pattern = strings.Trim(pattern, "\"")
            pattern = strings.ReplaceAll(pattern, "\\\"", "\"")
            definitionsArray[i].ValidationRegExp = pattern
        }
    }

    // Initialize ResourceDefinitions map
    for i := range definitionsArray {
        def := &definitionsArray[i]
        ResourceDefinitions[def.ResourceTypeName] = def
        ResourceMaps[def.ResourceTypeName] = def.ResourceTypeName
    }
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
