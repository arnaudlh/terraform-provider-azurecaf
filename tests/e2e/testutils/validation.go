package testutils

import (
	"regexp"
	"strings"
	"testing"
)

func ValidateResourceOutput(t *testing.T, resourceType string, resourceOutput, dataOutput string) {
	t.Logf("Testing resource type: %s", resourceType)
	t.Logf("Resource output: %s", resourceOutput)
	t.Logf("Data source output: %s", dataOutput)

	def, ok := GetResourceByType(resourceType)
	if !ok {
		t.Fatalf("Resource type %s not found in definitions", resourceType)
	}

	if resourceOutput != dataOutput {
		t.Errorf("Resource output %s does not match data source output %s", resourceOutput, dataOutput)
		return
	}

	nameToValidate := resourceOutput

	// Validate length constraints
	if len(nameToValidate) < def.MinLength {
		t.Errorf("Resource name %s length %d is less than minimum length %d", nameToValidate, len(nameToValidate), def.MinLength)
	}
	if len(nameToValidate) > def.MaxLength {
		t.Errorf("Resource name %s length %d exceeds maximum length %d", nameToValidate, len(nameToValidate), def.MaxLength)
	}

	// Validate case sensitivity
	if def.LowerCase && nameToValidate != strings.ToLower(nameToValidate) {
		t.Errorf("Resource name %s should be lowercase", nameToValidate)
	}

	// Validate slug presence and placement
	if def.CafPrefix != "" {
		slugIndex := strings.Index(strings.ToLower(nameToValidate), strings.ToLower(def.CafPrefix))
		if slugIndex == -1 {
			t.Errorf("Resource name %s does not contain slug %s", nameToValidate, def.CafPrefix)
		} else if slugIndex > 0 && !strings.ContainsRune("-_.", rune(nameToValidate[slugIndex-1])) {
			t.Errorf("Resource name %s has incorrectly placed slug %s - should be at start or after separator", nameToValidate, def.CafPrefix)
		}
	}

	// Validate regex pattern
	if def.ValidationRegExp != "" {
		cleanRegex := strings.Trim(def.ValidationRegExp, "\"")
		pattern, err := regexp.Compile(cleanRegex)
		if err != nil {
			t.Logf("Warning: Invalid validation regex pattern %s: %v", cleanRegex, err)
		} else if !pattern.MatchString(nameToValidate) {
			t.Errorf("Resource name %s does not match validation pattern %s", nameToValidate, cleanRegex)
		}
	}
}
