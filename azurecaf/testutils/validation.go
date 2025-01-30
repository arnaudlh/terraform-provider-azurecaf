package testutils

import (
	"regexp"
	"strings"
	"testing"
)

func ValidateResourceOutput(t *testing.T, resourceType string, resourceOutput, dataOutput string) {
	t.Helper()
	t.Logf("Testing resource type: %s", resourceType)
	t.Logf("Resource output: %s", resourceOutput)
	t.Logf("Data source output: %s", dataOutput)

	defs := loadResourceDefinitions()
	def, ok := defs[resourceType]
	if !ok {
		t.Fatalf("Resource type %s not found in definitions", resourceType)
	}

	if resourceOutput != dataOutput {
		t.Errorf("Resource output %s does not match data source output %s", resourceOutput, dataOutput)
		return
	}

	nameToValidate := resourceOutput
	nameLength := len(nameToValidate)

	if nameLength < def.MinLength || nameLength > def.MaxLength {
		t.Errorf("Resource name %s length %d is outside allowed range [%d, %d]", nameToValidate, nameLength, def.MinLength, def.MaxLength)
	}

	if def.LowerCase && nameToValidate != strings.ToLower(nameToValidate) {
		t.Errorf("Resource name %s should be lowercase", nameToValidate)
	}

	if def.CafPrefix != "" {
		slugIndex := strings.Index(strings.ToLower(nameToValidate), strings.ToLower(def.CafPrefix))
		if slugIndex == -1 {
			t.Errorf("Resource name %s does not contain slug %s", nameToValidate, def.CafPrefix)
		} else if slugIndex > 0 {
			prevChar := rune(nameToValidate[slugIndex-1])
			if !strings.ContainsRune("-_.", prevChar) {
				t.Errorf("Resource name %s has incorrectly placed slug %s - should be at start or after separator (-, _, or .)", nameToValidate, def.CafPrefix)
			}
		}
	}

	if def.ValidationRegExp != "" {
		pattern, err := regexp.Compile(def.ValidationRegExp)
		if err != nil {
			t.Errorf("Invalid validation regex pattern %s: %v", def.ValidationRegExp, err)
		} else if !pattern.MatchString(nameToValidate) {
			t.Errorf("Resource name %s does not match validation pattern %s", nameToValidate, def.ValidationRegExp)
		}
	}
}
