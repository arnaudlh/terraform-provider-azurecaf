// Package testutils provides utilities for E2E testing of the azurecaf provider
package testutils

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "testing"
)

type ResourceDefinition struct {
    ResourceTypeName  string `json:"name"`
    CafPrefix        string `json:"slug,omitempty"`
    MinLength        int    `json:"min_length"`
    MaxLength        int    `json:"max_length"`
    LowerCase        bool   `json:"lowercase,omitempty"`
    RegEx            string `json:"regex,omitempty"`
    ValidationRegExp string `json:"validation_regexp,omitempty"`
    Scope           string `json:"scope,omitempty"`
}

func GetResourceDefinitions() (map[string]ResourceDefinition, error) {
    wd, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("failed to get working directory: %v", err)
    }

    var jsonPath string
    for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
        path := filepath.Join(dir, "resourceDefinition.json")
        if _, err := os.Stat(path); err == nil {
            jsonPath = path
            break
        }
    }

    if jsonPath == "" {
        return nil, fmt.Errorf("resourceDefinition.json not found in any parent directory")
    }

    data, err := os.ReadFile(jsonPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read resource definitions: %v", err)
    }

    var definitionsArray []ResourceDefinition
    if err := json.Unmarshal(data, &definitionsArray); err != nil {
        return nil, fmt.Errorf("failed to parse resource definitions: %v", err)
    }

    definitions := make(map[string]ResourceDefinition)
    for _, def := range definitionsArray {
        definitions[def.ResourceTypeName] = def
    }

    return definitions, nil
}

func ValidateResourceOutput(t *testing.T, resourceType, resourceOutput, dataOutput string) error {
    if resourceOutput != dataOutput {
        return fmt.Errorf("resource output (%s) does not match data source output (%s)", resourceOutput, dataOutput)
    }

    defs, err := GetResourceDefinitions()
    if err != nil {
        return fmt.Errorf("failed to get resource definitions: %v", err)
    }

    def, ok := defs[resourceType]
    if !ok {
        return fmt.Errorf("resource type %s not found in definitions", resourceType)
    }

    if def.ValidationRegExp != "" {
        pattern, err := regexp.Compile(def.ValidationRegExp)
        if err != nil {
            return fmt.Errorf("invalid validation regex for %s: %v", resourceType, err)
        }

        if !pattern.MatchString(resourceOutput) {
            return fmt.Errorf("resource output %q does not match validation pattern %q", resourceOutput, def.ValidationRegExp)
        }
    }

    if def.CafPrefix != "" && def.LowerCase {
        if resourceOutput != "" && !regexp.MustCompile("^"+def.CafPrefix+"(-|$)").MatchString(resourceOutput) {
            return fmt.Errorf("resource output %q does not start with expected prefix %q", resourceOutput, def.CafPrefix)
        }
    }

    return nil
}
