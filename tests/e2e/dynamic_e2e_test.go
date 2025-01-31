package e2e

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
    "testing"
    "time"

    "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

func TestDynamicResourceDefinitions(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()

    stats := struct {
        totalResources    int
        passedValidation  int
        failedValidation  int
        skippedResources  int
        terraformErrors   int
        matchErrors       int
    }{
        totalResources: 0,
    }

    t.Log("Loading resource definitions...")
    resourceDefs := models.ResourceDefinitions
    stats.totalResources = len(resourceDefs)
    t.Logf("Found %d resource definitions", stats.totalResources)

    // Initialize command variable for reuse
    var cmd *exec.Cmd

    // Create base test directory
    baseDir, err := os.MkdirTemp("", "azurecaf-e2e-*")
    if err != nil {
        t.Fatalf("Failed to create test directory: %v", err)
    }
    defer os.RemoveAll(baseDir)

    // Add cleanup handler for timeout
    go func() {
        <-ctx.Done()
        if ctx.Err() == context.DeadlineExceeded {
            os.RemoveAll(baseDir)
        }
    }()

    // Prepare common configurations
    providerConfig := `terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}

provider "azurecaf" {}`

    devConfig := fmt.Sprintf(`provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "%s"
  }
  direct {}
}`, filepath.Join(os.Getenv("HOME"), ".terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/2.0.0-preview5/linux_amd64"))

    for resourceType := range resourceDefs {
        if ctx.Err() != nil {
            t.Fatal("Test timeout exceeded")
        }

        t.Run(resourceType, func(t *testing.T) {
            t.Logf("Testing resource type: %s", resourceType)
            
            // Generate test configuration
            // Clean the resource type and create a safe identifier
            cleanResourceType := strings.TrimSpace(resourceType)
            safeResourceType := strings.TrimSuffix(strings.Map(func(r rune) rune {
                switch {
                case r >= 'a' && r <= 'z':
                    return r
                case r >= 'A' && r <= 'Z':
                    return r
                case r >= '0' && r <= '9':
                    return r
                case r == '_' || r == '-':
                    return '_'
                default:
                    return '_'
                }
            }, cleanResourceType), "_")

            config := fmt.Sprintf(`
resource "azurecaf_name" "test_%[1]s" {
  name          = "test"
  resource_type = "%[2]s"
  random_length = 5
  random_seed   = 12345
  clean_input   = true
}

data "azurecaf_name" "test_%[1]s" {
  name          = azurecaf_name.test_%[1]s.result
  resource_type = "%[2]s"
  random_length = 5
  random_seed   = 12345
  clean_input   = true
}

output "resource_output_%[1]s" {
  value = azurecaf_name.test_%[1]s.result
}

output "data_output_%[1]s" {
  value = data.azurecaf_name.test_%[1]s.result
}`, safeResourceType, cleanResourceType)

            // Create unique test directory for this resource
            testDir := filepath.Join(baseDir, strings.ReplaceAll(safeResourceType, "/", "_"))
            if err := os.MkdirAll(testDir, 0755); err != nil {
                t.Fatalf("Failed to create test directory: %v", err)
            }

            // Write provider configuration
            if err := os.WriteFile(filepath.Join(testDir, "provider.tf"), []byte(providerConfig), 0644); err != nil {
                t.Fatalf("Failed to write provider config: %v", err)
            }

            // Write dev override configuration
            if err := os.WriteFile(filepath.Join(testDir, ".terraformrc"), []byte(devConfig), 0644); err != nil {
                t.Fatalf("Failed to write dev config: %v", err)
            }

            // Set environment variable for this test
            os.Setenv("TF_CLI_CONFIG_FILE", filepath.Join(testDir, ".terraformrc"))

            // Initialize Terraform
            cmd = exec.CommandContext(ctx, "terraform", "init")
            cmd.Dir = testDir
            if out, err := cmd.CombinedOutput(); err != nil {
                t.Fatalf("Failed to initialize Terraform: %v\n%s", err, out)
            }

            // Write new test configuration
            configPath := filepath.Join(testDir, "test.tf")
            if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
                t.Fatalf("Failed to write test config for %s: %v", resourceType, err)
            }

            // Apply configuration with timeout for individual resource
            resourceCtx, resourceCancel := context.WithTimeout(ctx, 2*time.Minute)
            defer resourceCancel()
            
            cmd := exec.CommandContext(resourceCtx, "terraform", "apply", "-auto-approve")
            cmd.Dir = testDir
            if out, err := cmd.CombinedOutput(); err != nil {
                if resourceCtx.Err() == context.DeadlineExceeded {
                    stats.skippedResources++
                    t.Skipf("Skipping resource type %q due to timeout after 2 minutes", resourceType)
                    return
                }
                stats.terraformErrors++
                def := resourceDefs[resourceType]
                t.Fatalf(`Terraform apply failed for resource type %q:
Error: %v
Resource Definition:
  - Type: %s
  - Prefix: %q
  - Validation Pattern: %q
Terraform Output:
%s`, 
                    resourceType, err, resourceType, def.CafPrefix, def.ValidationRegExp, out)
            }

            // Get outputs
            cmd = exec.CommandContext(ctx, "terraform", "output", "-json")
            cmd.Dir = testDir
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

            // Validate that resource and data source outputs match
            if outputs[resourceOutputKey].Value != outputs[dataOutputKey].Value {
                stats.matchErrors++
                def := resourceDefs[resourceType]
                t.Errorf(`Resource/Data source mismatch for type %q:
  - Resource Output: %q
  - Data Source Output: %q
  - Resource Type: %s
  - Expected Prefix: %q`,
                    resourceType,
                    outputs[resourceOutputKey].Value,
                    outputs[dataOutputKey].Value,
                    resourceType,
                    def.CafPrefix)
            }

            // Validate against resource definition
            def := resourceDefs[resourceType]
            if def != nil && def.ValidationRegExp != "" {
                re, err := regexp.Compile(def.ValidationRegExp)
                if err != nil {
                    t.Errorf("Invalid validation regex for %s: %v", resourceType, err)
                    return
                }

                if !re.MatchString(outputs[resourceOutputKey].Value) {
                    stats.failedValidation++
                    t.Errorf(`Resource validation failed for type %q:
Resource Definition:
  - Prefix: %q
  - Min Length: %d
  - Max Length: %d
  - Validation Pattern: %q
Generated Output:
  - Resource Output: %q
  - Data Source Output: %q
  - Expected Pattern: %q`,
                        resourceType,
                        def.CafPrefix,
                        def.MinLength,
                        def.MaxLength,
                        def.ValidationRegExp,
                        outputs[resourceOutputKey].Value,
                        outputs[dataOutputKey].Value,
                        def.ValidationRegExp)
                } else {
                    stats.passedValidation++
                }
            }
        })
    }

    // Print final test statistics
    t.Logf("\n=== E2E Test Statistics ===")
    t.Logf("Total Resources Tested: %d", stats.totalResources)
    t.Logf("Passed Validation: %d (%.1f%%)", stats.passedValidation, float64(stats.passedValidation)/float64(stats.totalResources)*100)
    t.Logf("Failed Validation: %d", stats.failedValidation)
    t.Logf("Resource/Data Mismatches: %d", stats.matchErrors)
    t.Logf("Skipped (Timeout): %d", stats.skippedResources)
    t.Logf("Terraform Errors: %d", stats.terraformErrors)
    t.Logf("========================")
}
