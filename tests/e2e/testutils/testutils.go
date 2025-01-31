// Package testutils provides utilities for E2E testing of the azurecaf provider
package testutils

import (
    "regexp"
    "testing"

    "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

func GetResourceDefinitions() map[string]*models.ResourceStructure {
    return models.ResourceDefinitions
}

func ValidateResourceOutput(t *testing.T, resourceType, resourceOutput, dataOutput string) {
    if resourceOutput != dataOutput {
        t.Errorf("Resource output (%s) does not match data source output (%s)", resourceOutput, dataOutput)
        return
    }

    defs := GetResourceDefinitions()
    def, ok := defs[resourceType]
    if !ok {
        t.Fatalf("Resource type %s not found in definitions", resourceType)
        return
    }

    if def.ValidationRegExp != "" {
        re, err := regexp.Compile(def.ValidationRegExp)
        if err != nil {
            t.Fatalf("Invalid validation regex for %s: %v", resourceType, err)
            return
        }

        if !re.MatchString(resourceOutput) {
            t.Errorf("Resource output %q does not match validation pattern %q", resourceOutput, def.ValidationRegExp)
        }
    }

    if def.CafPrefix != "" && def.LowerCase {
        if resourceOutput != "" && !regexp.MustCompile("^"+def.CafPrefix+"(-|$)").MatchString(resourceOutput) {
            t.Errorf("Resource output %q does not start with expected prefix %q", resourceOutput, def.CafPrefix)
        }
    }
}
