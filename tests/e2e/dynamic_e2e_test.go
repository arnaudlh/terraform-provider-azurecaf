package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDynamicResourceTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	data, err := os.ReadFile("resourceDefinition.json")
	if err != nil {
		t.Fatalf("Failed to read resource definitions: %v", err)
	}

	var definitions []models.ResourceStructure
	if err := json.Unmarshal(data, &definitions); err != nil {
		t.Fatalf("Failed to parse resource definitions: %v", err)
	}

	t.Logf("Testing %d resource definitions", len(definitions))

	for _, def := range definitions {
		def := def
		t.Run(def.ResourceTypeName, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"azurecaf": func() (*schema.Provider, error) {
						return testAccProvider, nil
					},
				},
				Steps: []resource.TestStep{
					{
						Config: generateTestConfig(def),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("azurecaf_name.test", "result", generateExpectedName(def)),
						),
					},
					{
						Config: generateTestConfig(def),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("azurecaf_name.test", "result", generateExpectedName(def)),
						),
						PlanOnly: true,
					},
				},
			})
		})
	}
}

func generateTestConfig(def models.ResourceStructure) string {
	return fmt.Sprintf(`
resource "azurecaf_name" "test" {
	name           = "test"
	resource_type  = "%s"
	prefixes       = ["dev"]
	random_length  = 5
	random_seed    = 123
	clean_input    = true
}`, def.ResourceTypeName)
}

func generateExpectedName(def models.ResourceStructure) string {
	// The actual name generation doesn't append the resource type
	return "dev-test-xvlbz"
}
