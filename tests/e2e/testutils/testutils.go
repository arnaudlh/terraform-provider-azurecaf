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

func GetResourceDefinitions() map[string]interface{} {
    wd, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("failed to get working directory: %v", err))
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
        panic("resourceDefinition.json not found in any parent directory")
    }

    data, err := os.ReadFile(jsonPath)
    if err != nil {
        panic(fmt.Sprintf("failed to read resource definitions: %v", err))
    }

    var definitionsArray []map[string]interface{}
    if err := json.Unmarshal(data, &definitionsArray); err != nil {
        panic(fmt.Sprintf("failed to parse resource definitions: %v", err))
    }

    definitions := make(map[string]interface{})
    for _, def := range definitionsArray {
        if name, ok := def["name"].(string); ok {
            definitions[name] = def
        }
    }

    return definitions
}

func ValidateResourceOutput(t *testing.T, resourceType, resourceOutput, dataOutput string) {
    if resourceOutput != dataOutput {
        t.Errorf("Resource output (%s) does not match data source output (%s)", resourceOutput, dataOutput)
        return
    }

    defs := GetResourceDefinitions()
    def, ok := defs[resourceType].(map[string]interface{})
    if !ok {
        t.Fatalf("Resource type %s not found in definitions", resourceType)
        return
    }

    if pattern, ok := def["validation_regexp"].(string); ok && pattern != "" {
        re, err := regexp.Compile(pattern)
        if err != nil {
            t.Fatalf("Invalid validation regex for %s: %v", resourceType, err)
            return
        }

        if !re.MatchString(resourceOutput) {
            t.Errorf("Resource output %q does not match validation pattern %q", resourceOutput, pattern)
        }
    }

    if prefix, ok := def["slug"].(string); ok && prefix != "" {
        if lowercase, ok := def["lowercase"].(bool); ok && lowercase {
            if resourceOutput != "" && !regexp.MustCompile("^"+prefix+"(-|$)").MatchString(resourceOutput) {
                t.Errorf("Resource output %q does not start with expected prefix %q", resourceOutput, prefix)
            }
        }
    }
}
