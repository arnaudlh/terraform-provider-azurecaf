package e2e

import (
    "context"
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
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()

    t.Log("Loading resource definitions...")
    defs := testutils.GetResourceDefinitions()
    t.Logf("Found %d resource definitions", len(defs))

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

    // Initialize Terraform
    cmd := exec.CommandContext(ctx, "terraform", "init")
    cmd.Dir = baseDir
    if out, err := cmd.CombinedOutput(); err != nil {
        t.Fatalf("Failed to initialize Terraform: %v\n%s", err, out)
    }

    for resourceType := range defs {
        if ctx.Err() != nil {
            t.Fatal("Test timeout exceeded")
        }

        t.Run(resourceType, func(t *testing.T) {
            t.Logf("Testing resource type: %s", resourceType)
            
            // Generate test configuration
            config := fmt.Sprintf(`
resource "azurecaf_name" "test" {
  name          = "test-%s"
  resource_type = "%s"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test" {
  name          = azurecaf_name.test.result
  resource_type = "%s"
  random_length = 5
  clean_input   = true
}

output "resource_output" {
  value = azurecaf_name.test.result
}

output "data_output" {
  value = data.azurecaf_name.test.result
}`, resourceType, resourceType, resourceType)

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

            var outputs struct {
                ResourceOutput struct{ Value string }
                DataOutput    struct{ Value string }
            }
            if err := json.Unmarshal(output, &outputs); err != nil {
                t.Fatalf("Failed to parse outputs for %s: %v", resourceType, err)
            }

            testutils.ValidateResourceOutput(t, resourceType, outputs.ResourceOutput.Value, outputs.DataOutput.Value)
        })
    }
}
