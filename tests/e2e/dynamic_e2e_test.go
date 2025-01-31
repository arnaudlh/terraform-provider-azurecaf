package e2e

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/aztfmod/terraform-provider-azurecaf/tests/e2e/testutils"
)

func TestDynamicResourceDefinitions(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()

    t.Log("Loading resource definitions...")
    resourceDefs := testutils.GetResourceDefinitions()
    t.Logf("Found %d resource definitions", len(resourceDefs))

    // Create base test directory
    baseDir, err := os.MkdirTemp("", "azurecaf-e2e-*")
    if err != nil {
        t.Fatalf("Failed to create test directory: %v", err)
    }
    defer os.RemoveAll(baseDir)

    // Write provider configuration
    providerConfig := `
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}

provider "azurecaf" {}`

    if err := os.WriteFile(filepath.Join(baseDir, "provider.tf"), []byte(providerConfig), 0644); err != nil {
        t.Fatalf("Failed to write provider config: %v", err)
    }

    // Create dev override configuration
    devConfig := fmt.Sprintf(`provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "%s"
  }
  direct {}
}`, filepath.Join(os.Getenv("HOME"), ".terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/2.0.0-preview5/linux_amd64"))

    if err := os.WriteFile(filepath.Join(baseDir, ".terraformrc"), []byte(devConfig), 0644); err != nil {
        t.Fatalf("Failed to write dev config: %v", err)
    }

    // Set development override environment variable
    os.Setenv("TF_CLI_CONFIG_FILE", filepath.Join(baseDir, ".terraformrc"))

    // Initialize Terraform
    cmd := exec.CommandContext(ctx, "terraform", "init")
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
            
            // Generate test configuration
            // Create a safe resource name by replacing non-alphanumeric characters
            safeResourceType := strings.ReplaceAll(strings.ReplaceAll(resourceType, "-", "_"), ".", "_")
            
            config := fmt.Sprintf(`
resource "azurecaf_name" "test_%[1]s" {
  name          = "test-%[2]s"
  resource_type = "%[2]s"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_%[1]s" {
  name          = azurecaf_name.test_%[1]s.result
  resource_type = "%[2]s"
  random_length = 5
  clean_input   = true
}

output "resource_output_%[1]s" {
  value = azurecaf_name.test_%[1]s.result
}

output "data_output_%[1]s" {
  value = data.azurecaf_name.test_%[1]s.result
}`, safeResourceType, resourceType)

            configPath := filepath.Join(baseDir, fmt.Sprintf("%s.tf", resourceType))
            if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
                t.Fatalf("Failed to write test config for %s: %v", resourceType, err)
            }

            // Apply configuration
            cmd := exec.CommandContext(ctx, "terraform", "apply", "-auto-approve")
            cmd.Dir = baseDir
            if out, err := cmd.CombinedOutput(); err != nil {
                t.Fatalf("Failed to apply Terraform config for %s: %v\n%s", resourceType, err, out)
            }

            // Get outputs
            cmd = exec.CommandContext(ctx, "terraform", "output", "-json")
            cmd.Dir = baseDir
            output, err := cmd.Output()
            if err != nil {
                t.Fatalf("Failed to get outputs for %s: %v", resourceType, err)
            }

            resourceOutputKey := fmt.Sprintf("resource_output_%s", safeResourceType)
            dataOutputKey := fmt.Sprintf("data_output_%s", safeResourceType)
            
            var outputs map[string]struct{ Value string }
            if err := json.Unmarshal(output, &outputs); err != nil {
                t.Fatalf("Failed to parse outputs for %s: %v", resourceType, err)
            }

            testutils.ValidateResourceOutput(t, resourceType, outputs[resourceOutputKey].Value, outputs[dataOutputKey].Value)
        })
    }
}
