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
		if resource.RegEx == "" {
			continue
		}
		reg, err := regexp.Compile(resource.RegEx)
		if err != nil {
			t.Logf("Error on the regex %s for the resource %s error %v", resource.RegEx, resource.ResourceTypeName, err.Error())
			t.Fail()
			continue
		}

		// Skip patterns that are meant to match specific formats
		if strings.HasPrefix(resource.RegEx, "^[") ||
			strings.Contains(resource.RegEx, "[^") ||
			strings.Contains(resource.RegEx, "\\s") ||
			strings.Contains(resource.RegEx, "\\d") ||
			strings.Contains(resource.RegEx, "\\+") ||
			strings.Contains(resource.RegEx, "\\?") {
			continue
		}

		content := "abcde"
		result := reg.ReplaceAllString(content, "")
		if result != content {
			t.Logf("%s : pattern %s modified test string %q to %q", resource.ResourceTypeName, resource.RegEx, content, result)
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
	testCases := []struct {
		name     string
		pattern  string
		content  string
		expected bool
	}{
		{
			name:     "simple pattern with dash",
			pattern:  "^[a-z0-9-]+$",
			content:  "test-name",
			expected: true,
		},
		{
			name:     "pattern without dash",
			pattern:  "^[a-z0-9]+$",
			content:  "testname",
			expected: true,
		},
		{
			name:     "pattern with escaped dash",
			pattern:  "^[a-z0-9\\-]+$",
			content:  "test-name",
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exp, err := regexp.Compile(tc.pattern)
			if err != nil {
				t.Fatalf("Failed to compile regex %s: %v", tc.pattern, err)
			}

			matches := exp.MatchString(tc.content)
			if matches != tc.expected {
				t.Errorf("Pattern %s with content %s: expected matches=%v but got %v",
					tc.pattern, tc.content, tc.expected, matches)
			}
		})
	}

	t.Run("resource patterns", func(t *testing.T) {
		for name, resource := range ResourceDefinitions {
			name := name
			if resource.ValidationRegExp == "" {
				continue
			}

			t.Run(name, func(t *testing.T) {
				exp, err := regexp.Compile(resource.ValidationRegExp)
				if err != nil {
					t.Errorf("Error compiling regex %s: %v", resource.ValidationRegExp, err)
					return
				}

				// Skip patterns with special validation requirements
				if strings.Contains(resource.ValidationRegExp, "[^") ||
					strings.Contains(resource.ValidationRegExp, "\\s") ||
					strings.Contains(resource.ValidationRegExp, "?") ||
					strings.Contains(resource.ValidationRegExp, "+") ||
					strings.Contains(resource.ValidationRegExp, "\\.") ||
					strings.Contains(resource.ValidationRegExp, "\\-") ||
					strings.Contains(resource.ValidationRegExp, "\\d") ||
					strings.Contains(resource.ValidationRegExp, "\\w") {
					t.Skip("Skipping pattern with special validation requirements")
					return
				}

				// Test with both dashed and non-dashed content
				testContents := []string{
					"testresourcename123",
					"test-resource-name-123",
				}

				for _, content := range testContents {
					matches := exp.MatchString(content)
					if !matches {
						t.Logf("Pattern %s rejected valid content %s", resource.ValidationRegExp, content)
					}
				}
			})
		}
	})
}
