package e2e

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/aztfmod/terraform-provider-azurecaf/tests/e2e/testutils"
)

func TestDynamicResourceDefinitions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	t.Log("Loading resource definitions...")
	resourceDefs := testutils.GetResourceDefinitions()
	t.Logf("Found %d resource definitions", len(resourceDefs))

	// Create base test directory
	baseDir, err := os.MkdirTemp("", "azurecaf-e2e-*")
	if err != nil {
		t.Fatalf("Failed to create base test dir: %v", err)
	}
	defer os.RemoveAll(baseDir)

	// Write provider configuration
	providerConfig := `
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "2.0.0-preview5"
    }
  }
}

provider "azurecaf" {}`

	if err := os.WriteFile(filepath.Join(baseDir, "provider.tf"), []byte(providerConfig), 0644); err != nil {
		t.Fatalf("Failed to write provider config: %v", err)
	}

	// Build and install provider locally
	buildCmd := exec.Command("go", "build", "-o", "terraform-provider-azurecaf_v2.0.0-preview5")
	buildCmd.Dir = filepath.Join(baseDir, "..", "..")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build provider: %v\n%s", err, out)
	}

	// Create plugin directory
	pluginDir := filepath.Join(baseDir, ".terraform.d", "plugins", "registry.terraform.io", "aztfmod", "azurecaf", "2.0.0-preview5", "linux_amd64")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("Failed to create plugin directory: %v", err)
	}

	// Copy provider to plugin directory with correct name
	if err := exec.Command("cp", 
		filepath.Join(baseDir, "..", "..", "terraform-provider-azurecaf_v2.0.0-preview5"),
		filepath.Join(pluginDir, "terraform-provider-azurecaf_v2.0.0-preview5"),
	).Run(); err != nil {
		t.Fatalf("Failed to copy provider: %v", err)
	}

	// Create provider configuration
	terraformConfig := `terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "2.0.0-preview5"
    }
  }
}

provider "azurecaf" {}`

	if err := os.WriteFile(filepath.Join(baseDir, "provider.tf"), []byte(terraformConfig), 0644); err != nil {
		t.Fatalf("Failed to write provider config: %v", err)
	}

	// Create dev override configuration
	devConfig := fmt.Sprintf(`provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "%s"
  }
  direct {}
}`, pluginDir)

	if err := os.WriteFile(filepath.Join(os.Getenv("HOME"), ".terraformrc"), []byte(devConfig), 0644); err != nil {
		t.Fatalf("Failed to write dev config: %v", err)
	}

	// Set development override environment variable
	os.Setenv("TF_PLUGIN_MIRROR_DIR", pluginDir)

	// Initialize Terraform
	cmd := exec.Command("terraform", "init")
	cmd.Dir = baseDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to initialize Terraform: %v\n%s", err, out)
	}

	for resourceType := range resourceDefs {
		if ctx.Err() != nil {
			t.Fatal("Test timeout exceeded")
		}

		t.Run(resourceType, func(t *testing.T) {
			t.Logf("Testing resource type: %s", resourceType)
			
			// Generate unique resource names based on resource type hash to avoid conflicts
			var resourceHash string
			resourceHash = fmt.Sprintf("%x", sha256.Sum256([]byte(resourceType)))[:8]
			
			config := fmt.Sprintf(`
resource "azurecaf_name" "test_%s" {
  name          = "test"
  resource_type = "%s"
  random_length = 5
  random_seed   = 123
  clean_input   = true
  use_slug      = true
}

data "azurecaf_name" "test_%s_data" {
  name          = azurecaf_name.test_%s.result
  resource_type = "%s"
  random_length = 5
  random_seed   = 123
  clean_input   = true
  use_slug      = true
}

output "resource_%s" {
  value = azurecaf_name.test_%s.result
}

output "data_%s" {
  value = data.azurecaf_name.test_%s_data.result
}
`, resourceHash, resourceType,
   resourceHash, resourceHash, resourceType,
   resourceHash, resourceHash,
   resourceHash, resourceHash)

			// Write resource config
			if err := os.WriteFile(filepath.Join(baseDir, "main.tf"), []byte(config), 0644); err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}

			// Apply configuration
			cmd = exec.Command("terraform", "apply", "-auto-approve")
			cmd.Dir = baseDir
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("Failed to apply Terraform: %v\n%s", err, out)
			}

			// Get outputs
			cmd = exec.Command("terraform", "output", "-json")
			cmd.Dir = baseDir
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

			resourceOutput := outputs[fmt.Sprintf("resource_%s", resourceHash)].Value
			dataOutput := outputs[fmt.Sprintf("data_%s", resourceHash)].Value

			testutils.ValidateResourceOutput(t, resourceType, resourceOutput, dataOutput)
		})
	}
}
