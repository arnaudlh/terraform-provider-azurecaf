package testutils

import (
	"regexp"
	"strings"
	"testing"
)

// ValidateResourceOutput validates that the resource and data source outputs match the expected format
func ValidateResourceOutput(t *testing.T, resourceType string, resourceOutput, dataOutput string) {
	t.Helper()
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

	// Validate that resource and data source outputs match
	// Validate that resource and data source outputs match
	if resourceOutput != dataOutput {
		t.Errorf("Resource output %s does not match data source output %s", resourceOutput, dataOutput)
		return
	}

	nameToValidate := resourceOutput

	// Validate length constraints
	nameLength := len(nameToValidate)
	if nameLength < def.MinLength {
		t.Errorf("Resource name %s length %d is less than minimum length %d", nameToValidate, nameLength, def.MinLength)
	}
	if nameLength > def.MaxLength {
		t.Errorf("Resource name %s length %d exceeds maximum length %d", nameToValidate, nameLength, def.MaxLength)
	}

	// Validate case sensitivity
	if def.LowerCase && nameToValidate != strings.ToLower(nameToValidate) {
		t.Errorf("Resource name %s should be lowercase", nameToValidate)
	}

	// Validate slug presence and placement
	if def.Slug != "" {
		slugIndex := strings.Index(strings.ToLower(nameToValidate), strings.ToLower(def.Slug))
		if slugIndex == -1 {
			t.Errorf("Resource name %s does not contain slug %s", nameToValidate, def.Slug)
		} else if slugIndex > 0 {
			prevChar := rune(nameToValidate[slugIndex-1])
			if !strings.ContainsRune("-_.", prevChar) {
				t.Errorf("Resource name %s has incorrectly placed slug %s - should be at start or after separator (-, _, or .)", nameToValidate, def.Slug)
			}
		}
	}

	// Validate regex pattern
	if def.ValidationRegex != "" {
		pattern, err := regexp.Compile(def.ValidationRegex)
		if err != nil {
			t.Errorf("Invalid validation regex pattern %s: %v", def.ValidationRegex, err)
		} else if !pattern.MatchString(nameToValidate) {
			t.Errorf("Resource name %s does not match validation pattern %s", nameToValidate, def.ValidationRegex)
		}
	}
}
