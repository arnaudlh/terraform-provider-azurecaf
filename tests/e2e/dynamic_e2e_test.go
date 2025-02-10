package e2e

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "testing"
    "time"

    "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    return err
}

func TestDynamicResourceTypes(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }
    t.Parallel()

    // Use a background context for the overall test suite
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
    defer cancel()

    // Initialize test statistics with mutex protection
    var stats struct {
        sync.Mutex
        totalResources    int
        passedValidation  int
        failedValidation  int
        skippedResources  int
        terraformErrors   int
        matchErrors       int
        startTime        time.Time
        endTime          time.Time
        resourceDirs     []string
    }
    stats.startTime = time.Now()
    defer func() {
        stats.endTime = time.Now()
        t.Logf("\n=== E2E Test Statistics ===")
        t.Logf("Total Resources Tested: %d", stats.totalResources)
        t.Logf("Passed Validation: %d (%.1f%%)", stats.passedValidation, float64(stats.passedValidation)/float64(stats.totalResources)*100)
        t.Logf("Failed Validation: %d", stats.failedValidation)
        t.Logf("Resource/Data Mismatches: %d", stats.matchErrors)
        t.Logf("Skipped (Timeout): %d", stats.skippedResources)
        t.Logf("Terraform Errors: %d", stats.terraformErrors)
        t.Logf("Total Duration: %v", stats.endTime.Sub(stats.startTime))
        t.Logf("========================")

        // Clean up test directories
        for _, dir := range stats.resourceDirs {
            if err := os.RemoveAll(dir); err != nil {
                t.Logf("Warning: Failed to clean up test directory %s: %v", dir, err)
            }
        }
    }()

    t.Log("Loading resource definitions...")
    // Load resource definitions from JSON file
    searchPaths := []string{
        "resourceDefinition.json",
        "../../resourceDefinition.json",
    }

    var data []byte
    var err error
    for _, path := range searchPaths {
        data, err = os.ReadFile(path)
        if err == nil {
            t.Logf("Found resource definitions at: %s", path)
            break
        }
    }
    if err != nil {
        t.Fatalf("Failed to read resource definitions from any search path: %v", err)
    }

    var definitionsArray []models.ResourceStructure
    if err := json.Unmarshal(data, &definitionsArray); err != nil {
        t.Fatalf("Failed to parse resource definitions: %v", err)
    }

    resourceDefs := make(map[string]*models.ResourceStructure)
    for i := range definitionsArray {
        def := &definitionsArray[i]
        if def.ValidationRegExp != "" {
            def.ValidationRegExp = strings.Trim(def.ValidationRegExp, "\"")
            def.ValidationRegExp = strings.ReplaceAll(def.ValidationRegExp, "\\\"", "\"")
        }
        resourceDefs[def.ResourceTypeName] = def
    }

    stats.totalResources = len(resourceDefs)
    t.Logf("Found %d resource definitions", stats.totalResources)

    // Initialize test environment with proper cleanup
    testDir := filepath.Join(".", "terraform.d", fmt.Sprintf("azurecaf-e2e-%d", time.Now().Unix()))
    if err := os.MkdirAll(testDir, 0755); err != nil {
        t.Fatalf("Failed to create test directory: %v", err)
    }
    defer func() {
        if cleanupErr := os.RemoveAll(testDir); cleanupErr != nil {
            t.Logf("Warning: Failed to clean up test directory: %v", cleanupErr)
        }
    }()

    // Create plugin and cache directories with proper permissions
    pluginDir := filepath.Join(testDir, "terraform.d", "plugins", "registry.terraform.io", "aztfmod", "azurecaf", "2.0.0-preview5", "linux_amd64")
    if err := os.MkdirAll(pluginDir, 0755); err != nil {
        t.Fatalf("Failed to create plugin directory: %v", err)
    }

    // Build provider binary once
    providerBinary := filepath.Join(pluginDir, "terraform-provider-azurecaf_v2.0.0-preview5")
    buildCmd := exec.Command("go", "build", "-o", providerBinary, "../..")
    buildCmd.Env = append(os.Environ(),
        "CGO_ENABLED=0",
        "GOOS=linux",
        "GOARCH=amd64",
    )
    if out, err := buildCmd.CombinedOutput(); err != nil {
        t.Fatalf("Failed to build provider: %v\nOutput: %s", err, out)
    }
    if err := os.Chmod(providerBinary, 0755); err != nil {
        t.Fatalf("Failed to make provider binary executable: %v", err)
    }

    // Create cache directory
    cacheDir := filepath.Join(testDir, "terraform.d", "plugin-cache")
    if err := os.MkdirAll(cacheDir, 0755); err != nil {
        t.Fatalf("Failed to create cache directory: %v", err)
    }

    // Ensure plugin directory is empty
    if err := os.RemoveAll(pluginDir); err != nil {
        t.Fatalf("Failed to clean plugin directory: %v", err)
    }
    if err := os.MkdirAll(pluginDir, 0755); err != nil {
        t.Fatalf("Failed to recreate plugin directory: %v", err)
    }
    
    // Ensure cleanup on test completion or timeout
    defer func() {
        if err := os.RemoveAll(testDir); err != nil {
            t.Logf("Warning: Failed to clean up test directory %s: %v", testDir, err)
        }
    }()

    // Set up test environment variables with proper paths
    envVars := map[string]string{
        "TF_ACC": "1",
        "TF_LOG": "DEBUG",
        "TF_LOG_PATH": filepath.Join(testDir, "terraform.log"),
        "TF_CLI_CONFIG_FILE": filepath.Join(testDir, ".terraformrc"),
        "TF_PROVIDER_VERSION": "2.0.0-preview5",
        "TF_IN_AUTOMATION": "true",
        "TF_INPUT": "false",
    }
    
    for key, value := range envVars {
        if err := os.Setenv(key, value); err != nil {
            t.Fatalf("Failed to set environment variable %s: %v", key, err)
        }
    }

    // Set up concurrency control with disk space management
    maxConcurrent := 3 // Reduced to manage disk space
    if val := os.Getenv("TEST_PARALLELISM"); val != "" {
        if n, err := strconv.Atoi(val); err == nil && n > 0 {
            maxConcurrent = n
        }
    }
    sem := make(chan struct{}, maxConcurrent)
    var wg sync.WaitGroup

    // Create shared provider directory in current directory to avoid /tmp
    var sharedProviderBinary string
    {
        sharedPluginDir := filepath.Join(".", "terraform.d", "plugins", "registry.terraform.io", "aztfmod", "azurecaf", "2.0.0-preview5", "linux_amd64")
        if err := os.MkdirAll(sharedPluginDir, 0755); err != nil {
            t.Fatalf("Failed to create shared plugin directory: %v", err)
        }

        // Build provider once and share it
        sharedProviderBinary = filepath.Join(sharedPluginDir, "terraform-provider-azurecaf_v2.0.0-preview5")
        buildCmd := exec.Command("go", "build", "-o", sharedProviderBinary)
        buildCmd.Dir = filepath.Join("..", "..")
        buildCmd.Env = append(os.Environ(),
            "CGO_ENABLED=0",
            "GOOS=linux",
            "GOARCH=amd64",
            "TMPDIR=.",  // Use current directory for temporary files
        )
        if out, err := buildCmd.CombinedOutput(); err != nil {
            t.Fatalf("Failed to build provider: %v\nOutput: %s", err, out)
        }
        if err := os.Chmod(sharedProviderBinary, 0755); err != nil {
            t.Fatalf("Failed to make provider binary executable: %v", err)
        }
        t.Log("Successfully built shared provider binary")
    }

    // Set TMPDIR for all terraform operations
    os.Setenv("TMPDIR", ".")
    
    // Create error channel for collecting test results
    errChan := make(chan error, len(resourceDefs))
    
    // Clean test environment once at start
    testDirs := []string{
        filepath.Join(os.Getenv("HOME"), ".terraform.d/plugins"),
        filepath.Join(os.Getenv("HOME"), ".terraform.d/plugin-cache"),
    }
    for _, dir := range testDirs {
        if err := os.RemoveAll(dir); err != nil {
            t.Logf("Warning: Failed to clean up directory %s: %v", dir, err)
        }
    }
    
    // Set up provider configuration with proper directory structure
    {
        providerDir := filepath.Join(testDir, "terraform.d", "plugins", "registry.terraform.io", "aztfmod", "azurecaf", "2.0.0-preview5", "linux_amd64")
        if cleanErr := os.RemoveAll(providerDir); cleanErr != nil && !os.IsNotExist(cleanErr) {
            t.Logf("Warning: Failed to clean up provider directory: %v", cleanErr)
        }
        if mkdirErr := os.MkdirAll(providerDir, 0755); mkdirErr != nil {
            t.Fatalf("Failed to create provider directory: %v", mkdirErr)
        }

        // Create plugin cache directory with full provider path structure
        pluginCacheDir := filepath.Join(testDir, "terraform.d", "plugin-cache", "registry.terraform.io", "aztfmod", "azurecaf", "2.0.0-preview5", "linux_amd64")
        if cleanErr := os.RemoveAll(pluginCacheDir); cleanErr != nil && !os.IsNotExist(cleanErr) {
            t.Logf("Warning: Failed to clean up plugin cache directory: %v", cleanErr)
        }
        if mkdirErr := os.MkdirAll(pluginCacheDir, 0755); mkdirErr != nil {
            t.Fatalf("Failed to create plugin cache directory: %v", mkdirErr)
        }

        // Copy resource definition to provider directory
        resourceDefPath := filepath.Join("..", "..", "resourceDefinition.json")
        if copyErr := copyFile(resourceDefPath, filepath.Join(providerDir, "resourceDefinition.json")); copyErr != nil {
            t.Fatalf("Failed to copy resource definitions: %v", copyErr)
        }
    }

    // Verify resource definitions are loaded
    if len(resourceDefs) == 0 {
        t.Fatal("No resource definitions loaded. Check resourceDefinition.json path and content")
    }

    // Add cleanup handler for timeout
    go func() {
        <-ctx.Done()
        if ctx.Err() == context.DeadlineExceeded {
            os.RemoveAll(testDir)
        }
    }()

    t.Logf("Test environment initialized with %d resource definitions", len(resourceDefs))

    // Prepare common configurations with version constraint
    providerConfig := `terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "2.0.0-preview5"
    }
  }
}

provider "azurecaf" {
  random_seed = 12345
  random_length = 5
}`

    devConfig := fmt.Sprintf(`provider_installation {
  filesystem_mirror {
    path = "%s"
    include = ["registry.terraform.io/aztfmod/azurecaf"]
  }
  direct {
    exclude = ["registry.terraform.io/aztfmod/azurecaf"]
  }
}
plugin_cache_dir = "%s"
disable_plugin_tls_verify = true`, 
        filepath.Join(testDir, "terraform.d", "plugins"),
        filepath.Join(testDir, "terraform.d", "plugin-cache"))

    for resourceType := range resourceDefs {
        if ctx.Err() != nil {
            t.Fatal("Test timeout exceeded")
        }

        wg.Add(1)
        go func(resourceType string) {
            defer wg.Done()
            
            // Acquire semaphore for resource limiting
            sem <- struct{}{}
            defer func() { <-sem }()
            
            // Run subtest
            if err := func() error {
            
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
            resourceTestDir := filepath.Join(testDir, strings.ReplaceAll(safeResourceType, "/", "_"))
            if err := os.RemoveAll(resourceTestDir); err != nil {
                t.Fatalf("Failed to clean resource test directory: %v", err)
            }
            if err := os.MkdirAll(resourceTestDir, 0755); err != nil {
                t.Fatalf("Failed to create resource test directory: %v", err)
            }
            
            // Create resource plugin directory
            resourcePluginDir := filepath.Join(resourceTestDir, "terraform.d", "plugins", "registry.terraform.io", "aztfmod", "azurecaf", "2.0.0-preview5", "linux_amd64")
            if mkdirErr := os.MkdirAll(resourcePluginDir, 0755); mkdirErr != nil {
                t.Fatalf("Failed to create resource plugin directory: %v", mkdirErr)
            }
            
            // Copy provider binary to resource directory
            resourceProviderBinary := filepath.Join(resourcePluginDir, "terraform-provider-azurecaf_v2.0.0-preview5")
            if copyErr := copyFile(sharedProviderBinary, resourceProviderBinary); copyErr != nil {
                t.Fatalf("Failed to copy provider binary: %v", copyErr)
            }
            if chmodErr := os.Chmod(resourceProviderBinary, 0755); chmodErr != nil {
                t.Fatalf("Failed to make provider binary executable: %v", chmodErr)
            }
            
            stats.Lock()
            stats.resourceDirs = append(stats.resourceDirs, resourceTestDir)
            stats.Unlock()

            // Initialize Terraform with detailed logging and separate timeout
            initCtx, initCancel := context.WithTimeout(ctx, 60*time.Second)
            defer initCancel()

            // Write provider configuration
            if err := os.WriteFile(filepath.Join(resourceTestDir, "provider.tf"), []byte(providerConfig), 0644); err != nil {
                return fmt.Errorf("failed to write provider config: %v", err)
            }

            // Write dev override configuration
            if err := os.WriteFile(filepath.Join(resourceTestDir, ".terraformrc"), []byte(devConfig), 0644); err != nil {
                return fmt.Errorf("failed to write dev config: %v", err)
            }

            // Set environment variable for this test
            os.Setenv("TF_CLI_CONFIG_FILE", filepath.Join(resourceTestDir, ".terraformrc"))

            initCmd := exec.CommandContext(initCtx, "terraform", "init", "-no-color")
            initCmd.Dir = resourceTestDir
            initCmd.Env = append(os.Environ(),
                "TF_LOG=DEBUG",
                "TF_INPUT=false",
                fmt.Sprintf("TF_CLI_CONFIG_FILE=%s", filepath.Join(resourceTestDir, ".terraformrc")),
            )

            if out, err := initCmd.CombinedOutput(); err != nil {
                return fmt.Errorf("failed to initialize terraform: %v\nOutput: %s\nProvider Config:\n%s\nDev Override:\n%s",
                    err, out, providerConfig, devConfig)
            }

            // Write new test configuration
            configPath := filepath.Join(resourceTestDir, "test.tf")
            if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
                return fmt.Errorf("failed to write test config for %s: %v", resourceType, err)
            }

            // Apply configuration with timeout for individual resource
            resourceCtx, resourceCancel := context.WithTimeout(ctx, 2*time.Minute)
            defer resourceCancel()
            
            applyCmd := exec.CommandContext(resourceCtx, "terraform", "apply", "-auto-approve", "-no-color")
            applyCmd.Dir = resourceTestDir
            applyCmd.Env = append(os.Environ(),
                "TF_LOG=DEBUG",
                "TF_INPUT=false",
                fmt.Sprintf("TF_CLI_CONFIG_FILE=%s", filepath.Join(resourceTestDir, ".terraformrc")),
            )
            if out, err := applyCmd.CombinedOutput(); err != nil {
                if resourceCtx.Err() == context.DeadlineExceeded {
                    stats.Lock()
                    stats.skippedResources++
                    stats.Unlock()
                    return fmt.Errorf("timeout after 2 minutes for resource type %q", resourceType)
                }
                stats.Lock()
                stats.terraformErrors++
                stats.Unlock()
                def := resourceDefs[resourceType]
                return fmt.Errorf("terraform apply failed for resource type %q:\nError: %v\nResource Definition:\n  - Type: %s\n  - Prefix: %q\n  - Validation Pattern: %q\nOutput:\n%s",
                    resourceType, err, resourceType, def.CafPrefix, def.ValidationRegExp, out)
            }

            // Get outputs
            outputCmd := exec.CommandContext(ctx, "terraform", "output", "-json")
            outputCmd.Dir = resourceTestDir
            output, err := outputCmd.Output()
            if err != nil {
                t.Fatalf("Failed to get outputs for %s: %v", resourceType, err)
            }

            resourceOutputKey := fmt.Sprintf("resource_output_%s", safeResourceType)
            dataOutputKey := fmt.Sprintf("data_output_%s", safeResourceType)
            
            var outputs map[string]struct{ Value string }
            if err := json.Unmarshal(output, &outputs); err != nil {
                return fmt.Errorf("failed to parse outputs for %s: %v", resourceType, err)
            }

            // Validate that resource and data source outputs match
            if outputs[resourceOutputKey].Value != outputs[dataOutputKey].Value {
                stats.Lock()
                stats.matchErrors++
                stats.Unlock()
                return fmt.Errorf("resource/data source mismatch for type %q:\n  - Resource Output: %q\n  - Data Source Output: %q",
                    resourceType, outputs[resourceOutputKey].Value, outputs[dataOutputKey].Value)
            }

            // Validate against resource definition
            def := resourceDefs[resourceType]
            if def != nil && def.ValidationRegExp != "" {
                re, err := regexp.Compile(def.ValidationRegExp)
                if err != nil {
                    return fmt.Errorf("invalid validation regex for %s: %v", resourceType, err)
                }

                if !re.MatchString(outputs[resourceOutputKey].Value) {
                    stats.Lock()
                    stats.failedValidation++
                    stats.Unlock()
                    return fmt.Errorf("resource validation failed for type %q:\nResource Definition:\n  - Prefix: %q\n  - Min Length: %d\n  - Max Length: %d\n  - Validation Pattern: %q\nGenerated Output:\n  - Resource Output: %q\n  - Data Source Output: %q",
                        resourceType, def.CafPrefix, def.MinLength, def.MaxLength, def.ValidationRegExp, outputs[resourceOutputKey].Value, outputs[dataOutputKey].Value)
                }
                
                stats.Lock()
                stats.passedValidation++
                stats.Unlock()
            }
            return nil
        }(); err != nil {
            errChan <- fmt.Errorf("test failed for %s: %v", resourceType, err)
        }
    }(resourceType)
    }

    // Wait for all tests to complete
    wg.Wait()
    close(errChan)

    // Check for any test failures
    var failures []string
    for err := range errChan {
        failures = append(failures, err.Error())
    }
    if len(failures) > 0 {
        t.Errorf("Test failures:\n%s", strings.Join(failures, "\n"))
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
