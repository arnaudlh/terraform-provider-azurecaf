package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aztfmod/terraform-provider-azurecaf/tests/e2e/testutils"
)

func TestDynamicResourceDefinitions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	t.Log("Loading resource definitions...")
	resourceDefs := testutils.GetResourceDefinitions()
	resourceCount := len(resourceDefs)
	t.Logf("Found %d resource definitions", resourceCount)

	for resourceType := range resourceDefs {
		if ctx.Err() != nil {
			t.Fatal("Test timeout exceeded")
		}

		t.Run(resourceType, func(t *testing.T) {
			t.Logf("Testing resource type: %s", resourceType)
			
			config := fmt.Sprintf(`
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "2.0.0-preview5"
    }
  }
}

provider "azurecaf" {}

resource "azurecaf_name" "test" {
  name          = "test-%s"
  resource_type = "%s"
  random_length = 5
  random_seed   = 123
  clean_input   = true
  use_slug      = true
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
}

data "azurecaf_name" "test_data" {
  name          = azurecaf_name.test.result
  resource_type = "%s"
  random_length = 5
  random_seed   = 123
  clean_input   = true
  use_slug      = true
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
}

output "resource_result" {
  value = azurecaf_name.test.result
}

output "data_result" {
  value = data.azurecaf_name.test_data.result
}
`, strings.ReplaceAll(resourceType, "_", "-"), resourceType, resourceType)

			// Create temp dir for test
			tmpDir, err := os.MkdirTemp("", "azurecaf-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Write config
			if err := os.WriteFile(filepath.Join(tmpDir, "main.tf"), []byte(config), 0644); err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}

			// Initialize Terraform
			cmd := exec.Command("terraform", "init")
			cmd.Dir = tmpDir
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("Failed to initialize Terraform: %v\n%s", err, out)
			}

			// Apply configuration
			cmd = exec.Command("terraform", "apply", "-auto-approve")
			cmd.Dir = tmpDir
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("Failed to apply Terraform: %v\n%s", err, out)
			}

			// Get outputs
			cmd = exec.Command("terraform", "output", "-json")
			cmd.Dir = tmpDir
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to get outputs: %v", err)
			}

			var outputs map[string]struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(output, &outputs); err != nil {
				t.Fatalf("Failed to parse outputs: %v", err)
			}

			resourceOutput := outputs["resource_result"].Value
			dataOutput := outputs["data_result"].Value

			// Create test data from resource definition
			testutils.ValidateResourceOutput(t, resourceType, resourceOutput, dataOutput)
		})
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Chmod(0755)
}
