package azurecaf

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCafNamingValidation(id string, name string, expectedLength int, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		attrs := rs.Primary.Attributes

		result := attrs["result"]
		if len(result) != expectedLength {
			return fmt.Errorf("got %s %d result items; want %d", result, len(result), expectedLength)
		}
		if !strings.HasPrefix(result, prefix) {
			return fmt.Errorf("got %s which doesn't start with %s", result, prefix)
		}
		if !strings.Contains(result, name) {
			return fmt.Errorf("got %s which doesn't contain the name %s", result, name)
		}
		return nil
	}
}

func regexMatch(id string, exp *regexp.Regexp, requiredMatches int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		result := rs.Primary.Attributes["result"]

		if matches := exp.FindAllStringSubmatchIndex(result, -1); len(matches) != requiredMatches {
			return fmt.Errorf("result string is %s; did not match %s, got %d", result, exp, len(matches))
		}

		return nil
	}
}

// Unit tests for these functions have been moved to resource_name_engine_test.go

// Integration tests moved to resource_name_integration_test.go

func TestComposeName(t *testing.T) {
	namePrecedence := []string{"name", "random", "slug", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name := composeName("-", prefixes, "name", "slug", suffixes, "rd", 21, namePrecedence)
	expected := "a-b-slug-name-rd-c-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutCorrect(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name := composeName("-", prefixes, "name", "slug", suffixes, "rd", 19, namePrecedence)
	expected := "b-slug-name-rd-c-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutMaxLength(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{}
	suffixes := []string{}
	name := composeName("-", prefixes, "aaaaaaaaaa", "bla", suffixes, "", 10, namePrecedence)
	expected := "aaaaaaaaaa"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeNameCutCorrectSuffixes(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"a", "b"}
	suffixes := []string{"c", "d"}
	name := composeName("-", prefixes, "name", "slug", suffixes, "rd", 15, namePrecedence)
	expected := "slug-name-rd-c"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestComposeEmptyStringArray(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	prefixes := []string{"", "b"}
	suffixes := []string{"", "d"}
	name := composeName("-", prefixes, "", "", suffixes, "", 15, namePrecedence)
	expected := "b-d"
	if name != expected {
		t.Logf("Fail to generate name expected %s received %s", expected, name)
		t.Fail()
	}
}

func TestValidResourceType_validParameters(t *testing.T) {
	resourceType := "azurerm_resource_group"
	err := validateResourceType(resourceType)
	if err != nil {
		t.Logf("resource validation generated an unexpected error %s", err.Error())
		t.Fail()
	}
}

func TestValidResourceType_invalidParameters(t *testing.T) {
	resourceType := "azurerm_not_supported"
	err := validateResourceType(resourceType)
	if err == nil {
		t.Logf("resource validation did not generate an error while the input is invalid")
		t.Fail()
	}
}

func TestGetResourceNameValid(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", true, false, true, namePrecedence)
	expected := "a-b-rg-myrg-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameValidRsv(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_recovery_services_vault", "-", []string{"a", "b"}, "test", nil, "1234", true, false, true, namePrecedence)
	expected := "a-b-rsv-test-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameValidNoSlug(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", true, false, false, namePrecedence)
	expected := "a-b-myrg-1234"

	if err != nil {
		t.Logf("getResource Name generated an error %s", err.Error())
		t.Fail()
	}
	if expected != resourceName {
		t.Logf("invalid name, expected %s got %s", expected, resourceName)
		t.Fail()
	}
}

func TestGetResourceNameInvalidResourceType(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, err := getResourceName("azurerm_invalid", "-", []string{"a", "b"}, "myrg", nil, "1234", true, false, true, namePrecedence)
	expected := "a-b-rg-myrg-1234"

	if err == nil {
		t.Logf("Expected a validation error, got nil")
		t.Fail()
	}
	if expected == resourceName {
		t.Logf("valid name received while an error is expected")
		t.Fail()
	}
}

func TestGetResourceNamePassthrough(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	resourceName, _ := getResourceName("azurerm_resource_group", "-", []string{"a", "b"}, "myrg", nil, "1234", true, true, true, namePrecedence)
	expected := "myrg"

	if expected != resourceName {
		t.Logf("valid name received while an error is expected")
		t.Fail()
	}
}

// Test configurations moved to resource_name_integration_test.go
