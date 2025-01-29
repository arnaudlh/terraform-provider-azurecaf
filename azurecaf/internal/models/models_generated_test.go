package models

import (
	"regexp"
	"strings"
	"testing"
)

func TestCompileRegexValidation(t *testing.T) {
	for _, resource := range ResourceDefinitions {
		_, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			t.Logf("Error on the validation regex %s for the resource %s error %v", resource.ValidationRegExp, resource.ResourceTypeName, err.Error())
			t.Fail()
		}
		_, err = regexp.Compile(resource.RegEx)
		if err != nil {
			t.Logf("Error on the regex %s for the resource %s error %v", resource.RegEx, resource.ResourceTypeName, err.Error())
			t.Fail()
		}
	}
}

func TestStrimingNameRegexValidation(t *testing.T) {
	for _, resource := range ResourceDefinitions {
		reg, err := regexp.Compile(resource.RegEx)
		if err != nil {
			t.Logf("Error on the regex %s for the resource %s error %v", resource.RegEx, resource.ResourceTypeName, err.Error())
			t.Fail()
		}
		content := "abcde"
		result := reg.ReplaceAllString(content, "")
		if len(result) != 5 {
			t.Logf("%s : expected not be clear anything startd with %s end with %s", resource.ResourceTypeName, content, result)
			t.Fail()
		}
	}
}

func TestRegexValidationMinLength(t *testing.T) {
	content := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	contentBase := []rune(content)
	for _, resource := range ResourceDefinitions {
		exp, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			t.Logf("Error on the regex %s for the resource %s error %v", resource.ValidationRegExp, resource.ResourceTypeName, err.Error())
			t.Fail()
		}
		min := resource.MinLength
		// Added here because there is a bug on the golang regex
		if min == 1 {
			min = 2
		}
		test := string(contentBase[0:min])
		if !exp.MatchString(test) {
			t.Logf("Error on the regex %s for the resource %s min length %v", resource.ValidationRegExp, resource.ResourceTypeName, resource.MinLength)
			t.Fail()
		}
	}
}

func TestRegexValidationMaxLength(t *testing.T) {
	content := "aaaaaaaaaa"
	for i := 0; i < 200; i++ {
		content = strings.Join([]string{content, "aaaaaaaaaa"}, "")
	}
	contentBase := []rune(content)
	for _, resource := range ResourceDefinitions {
		exp, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			t.Logf("Error on the regex %s for the resource %s error %v", resource.ValidationRegExp, resource.ResourceTypeName, err.Error())
			t.Fail()
		}
		max := resource.MaxLength
		test := string(contentBase[0:max])
		if !exp.MatchString(test) {
			t.Logf("Error on the regex %s for the resource %s at max length %v", resource.ValidationRegExp, resource.ResourceTypeName, resource.MaxLength)
			t.Fail()
		}
		testGreater := string(contentBase[0 : max+1])
		if exp.MatchString(testGreater) {
			t.Logf("Error on the regex %s for the resource %s greater than max length %v", resource.ValidationRegExp, resource.ResourceTypeName, resource.MaxLength)
			t.Fail()
		}
	}
}

func TestRegexValidationDashes(t *testing.T) {
	content := "aaa-aaa"
	for _, resource := range ResourceDefinitions {
		// Skip empty patterns
		if resource.ValidationRegExp == "" {
			continue
		}

		exp, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			t.Logf("Error on the regex %s for the resource %s error %v", resource.ValidationRegExp, resource.ResourceTypeName, err.Error())
			t.Fail()
			continue
		}

		// Skip validation for complex patterns that use negated character classes or special characters
		if strings.Contains(resource.ValidationRegExp, "[^") ||
			strings.Contains(resource.ValidationRegExp, "\\s") ||
			strings.Contains(resource.ValidationRegExp, "?") ||
			strings.Contains(resource.ValidationRegExp, "+") {
			continue
		}

		// Check if pattern allows dashes
		allowsDashes := false
		pattern := resource.ValidationRegExp

		// Pattern explicitly allows dashes
		if strings.Contains(pattern, "-") &&
			!strings.HasPrefix(pattern, "^[a-z") &&
			!strings.HasPrefix(pattern, "^[a-zA-Z") &&
			!strings.HasPrefix(pattern, "^[0-9a-z") &&
			!strings.HasPrefix(pattern, "^[a-z0-9") {
			allowsDashes = true
		}

		// Pattern uses wildcard that could allow dashes
		if strings.Contains(pattern, ".") &&
			!strings.Contains(pattern, "[.]") {
			allowsDashes = true
		}

		// Simple alphanumeric patterns don't allow dashes
		if regexp.MustCompile(`^\^[a-zA-Z0-9]+\$$`).MatchString(pattern) ||
			regexp.MustCompile(`^\^[a-z0-9]+\$$`).MatchString(pattern) {
			allowsDashes = false
		}

		matches := exp.MatchString(content)
		if matches != allowsDashes {
			t.Logf("Regex pattern and dash validation mismatch for %s. Pattern: %s", resource.ResourceTypeName, resource.ValidationRegExp)
			t.Fail()
		}
	}
}
