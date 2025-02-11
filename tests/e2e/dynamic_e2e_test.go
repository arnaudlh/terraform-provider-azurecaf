package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDynamicResourceTypes(t *testing.T) {
	t.Log("Loading resource definitions...")
	data, err := os.ReadFile("resourceDefinition.json")
	if err != nil {
		t.Fatalf("Failed to read resource definitions: %v", err)
	}

	var definitions []models.ResourceStructure
	if err := json.Unmarshal(data, &definitions); err != nil {
		t.Fatalf("Failed to parse resource definitions: %v", err)
	}

	t.Logf("Found %d resource definitions", len(definitions))

	// Set up provider binary path
	providerPath := filepath.Join("terraform.d", "plugins", "registry.terraform.io", "aztfmod", "azurecaf", "2.0.0-preview5", "linux_amd64")
	if err := os.MkdirAll(providerPath, 0755); err != nil {
		t.Fatalf("Failed to create provider directory: %v", err)
	}

	// Copy provider binary from parent directory
	srcPath := "../../terraform-provider-azurecaf"
	dstPath := filepath.Join(providerPath, "terraform-provider-azurecaf_v2.0.0-preview5")
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("Failed to copy provider binary: %v", err)
	}

	// Make provider binary executable
	if err := os.Chmod(dstPath, 0755); err != nil {
		t.Fatalf("Failed to make provider binary executable: %v", err)
	}

	var totalPassed, totalFailed, totalMismatches, totalSkipped, totalErrors int
	startTime := time.Now()

	for _, def := range definitions {
		t.Run(def.ResourceTypeName, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: generateTestConfig(def),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("azurecaf_name.test", "result", generateExpectedName(def)),
						),
					},
				},
			})
		})
	}

	duration := time.Since(startTime)
	t.Logf(`
		=== E2E Test Statistics ===
		Total Resources Tested: %d
		Passed Validation: %d (%.1f%%)
		Failed Validation: %d
		Resource/Data Mismatches: %d
		Skipped (Timeout): %d
		Terraform Errors: %d
		Total Duration: %v
		========================
	`, len(definitions), totalPassed, float64(totalPassed)/float64(len(definitions))*100,
		totalFailed, totalMismatches, totalSkipped, totalErrors, duration)
}

func generateTestConfig(def models.ResourceStructure) string {
	return fmt.Sprintf(`
resource "azurecaf_name" "test" {
	name           = "test"
	resource_type  = "%s"
	prefixes       = ["dev"]
	random_length  = 5
	clean_input    = true
}`, def.ResourceTypeName)
}

func generateExpectedName(def models.ResourceStructure) string {
	return fmt.Sprintf("dev-test-xvlbz-%s", def.ResourceTypeName)
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}
