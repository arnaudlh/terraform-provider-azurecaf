package testutils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// ValidateResourceOutput validates the output of a resource against a data source
func ValidateResourceOutput(t *testing.T, resourceType string, resourceOutput string, dataOutput string) {
	if resourceOutput != dataOutput {
		t.Errorf("Resource output %q does not match data source output %q for resource type %s", resourceOutput, dataOutput, resourceType)
	}
}

// ValidateResourceName checks if a resource name matches the expected pattern
func ValidateResourceName(t *testing.T, resourceName string, pattern string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["result"]
		matched, err := regexp.MatchString(pattern, name)
		if err != nil {
			return fmt.Errorf("error matching pattern: %v", err)
		}
		if !matched {
			return fmt.Errorf("name %s does not match pattern %s", name, pattern)
		}
		return nil
	}
}
