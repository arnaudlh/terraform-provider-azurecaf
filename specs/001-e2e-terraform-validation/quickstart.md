# Quick Start: E2E Terraform Validation

**Feature**: 001-e2e-terraform-validation  
**Date**: 2025-10-30  
**Audience**: Provider maintainers and contributors

## Overview

This guide shows how to run the End-to-End Terraform validation pipeline locally or in CI/CD. The E2E validation tests the provider by building it, running Terraform plan/apply, and detecting drift to ensure zero regressions.

## Prerequisites

### Local Development

- **Go 1.24+**: Provider compilation
- **Terraform CLI**: One or more versions (1.5.x through 1.8.x)
- **Git**: Source code management
- **Make**: Build automation (optional)
- **Bash**: Script execution

### CI/CD (GitHub Actions)

- **GitHub Repository**: With Actions enabled
- **Secrets** (optional): Azure credentials if testing with real resources
- **Runner**: GitHub-hosted runners (Linux recommended)

## Quick Start (Local)

### 1. Run Quick Validation (5 minutes)

```bash
# From repository root
cd /path/to/terraform-provider-azurecaf

# Run quick validation (20 resources, 2 conventions)
./.github/e2e-tests/scripts/run-validation.sh --mode quick

# Expected output:
# ✓ Build provider (30s)
# ✓ Test 20 resources across 2 conventions (3m)
# ✓ Drift check (1m)
# ✓ Cleanup (30s)
# Result: 20/20 tests passed (100%)
```

### 2. Run Standard Validation (10 minutes)

```bash
# Standard validation (100 resources, all conventions)
./.github/e2e-tests/scripts/run-validation.sh --mode standard

# Or with specific resource filter
./.github/e2e-tests/scripts/run-validation.sh --mode standard \
  --resources "azurerm_storage_account,azurerm_key_vault"
```

### 3. Run Comprehensive Validation (30 minutes)

```bash
# Comprehensive validation (all 395 resources, all conventions)
./.github/e2e-tests/scripts/run-validation.sh --mode comprehensive

# Run in background and save log
nohup ./.github/e2e-tests/scripts/run-validation.sh --mode comprehensive > e2e.log 2>&1 &
```

## Quick Start (CI/CD)

### Automatic Triggers

**Pull Requests**: Quick validation runs automatically
```yaml
# Triggered on: PR opened, updated, synchronized
# Duration: ~7-8 minutes
# Resources tested: 20 representative types
# Conventions: cafclassic, cafrandom
```

**Main Branch**: Comprehensive validation runs on merge
```yaml
# Triggered on: Push to main
# Duration: ~25-30 minutes
# Resources tested: All 395 types
# Conventions: All 4 (cafclassic, cafrandom, random, passthrough)
```

### Manual Trigger

```bash
# Via GitHub CLI
gh workflow run e2e-validation-manual.yml \
  --field mode=standard \
  --field resource_types=azurerm_storage_account,azurerm_key_vault

# Via GitHub UI
# 1. Go to Actions tab
# 2. Select "E2E Validation (Manual)"
# 3. Click "Run workflow"
# 4. Choose options:
#    - Mode: quick | standard | comprehensive
#    - Resource types (optional filter)
#    - Naming conventions (optional filter)
#    - Terraform versions (optional)
```

## Understanding Results

### Success Output

```
✅ E2E Validation PASSED

Summary:
  Total Tests: 20
  Passed: 20
  Failed: 0
  Pass Rate: 100%
  Duration: 7m 30s

Stages:
  ✓ Build: 45s
  ✓ Plan: 2m 15s
  ✓ Apply: 3m 0s
  ✓ Drift Check: 1m 0s
  ✓ Cleanup: 30s

Details:
  See logs in /tmp/e2e-validation-{timestamp}/
```

### Failure Output

```
❌ E2E Validation FAILED

Summary:
  Total Tests: 20
  Passed: 18
  Failed: 2
  Pass Rate: 90%
  Duration: 5m 12s (stopped early)

Failed Tests:
  1. azurerm_storage_account (cafclassic)
     Stage: drift_check
     Error: Drift detected in resource name
     Details: Random suffix changed from 'abc' to 'xyz'
     
  2. azurerm_key_vault (cafrandom)
     Stage: plan
     Error: Name validation failed
     Details: Generated name exceeds maximum length (25 > 24)

Logs: /tmp/e2e-validation-{timestamp}/
```

## Common Use Cases

### Use Case 1: Test a Resource Type Change

**Scenario**: You modified resourceDefinition.json for azurerm_storage_account

```bash
# Test just storage account resources
./.github/e2e-tests/scripts/run-validation.sh --mode quick \
  --resources "azurerm_storage_account"

# Expected duration: ~2 minutes
# Tests: storage account with all 4 conventions
```

### Use Case 2: Test a Naming Convention Change

**Scenario**: You modified the cafclassic naming logic

```bash
# Test cafclassic convention across multiple resources
./.github/e2e-tests/scripts/run-validation.sh --mode standard \
  --conventions "cafclassic"

# Expected duration: ~8 minutes
# Tests: 100 resources with cafclassic only
```

### Use Case 3: Test Terraform Version Compatibility

**Scenario**: You want to ensure compatibility with Terraform 1.9.0

```bash
# Test with specific Terraform version
TF_VERSION=1.9.0 ./.github/e2e-tests/scripts/run-validation.sh --mode quick

# Or test multiple versions
for version in 1.5.7 1.6.6 1.7.5 1.8.0; do
  echo "Testing with Terraform $version..."
  TF_VERSION=$version ./.github/e2e-tests/scripts/run-validation.sh --mode quick
done
```

### Use Case 4: Debug a Failing Test

**Scenario**: A test is failing in CI and you want to reproduce locally

```bash
# Enable debug mode
DEBUG=1 ./.github/e2e-tests/scripts/run-validation.sh --mode quick \
  --resources "azurerm_storage_account" \
  --skip-cleanup

# Debug output includes:
# - Detailed Terraform logs
# - Provider debug logs
# - State file contents
# - Generated configurations

# State files preserved in: /tmp/e2e-validation-{timestamp}/
```

### Use Case 5: Performance Testing

**Scenario**: You want to measure validation performance

```bash
# Run with timing information
time ./.github/e2e-tests/scripts/run-validation.sh --mode quick \
  --parallel 20

# Compare with sequential execution
time ./.github/e2e-tests/scripts/run-validation.sh --mode quick \
  --parallel 1

# Expected results:
# - Parallel (20): ~3 minutes
# - Sequential (1): ~15 minutes
# - Speedup: ~5x
```

## Configuration Options

### Environment Variables

```bash
# Terraform version
export TF_VERSION=1.8.0

# Provider version (for testing against installed version)
export AZURECAF_VERSION=1.2.30

# Debug mode (verbose output)
export DEBUG=1

# Parallel jobs
export MAX_PARALLEL=10

# Timeout (seconds)
export TIMEOUT=1800

# Skip cleanup (for debugging)
export SKIP_CLEANUP=1

# Terraform log level
export TF_LOG=DEBUG
export TF_LOG_PATH=/tmp/terraform.log
```

### Command-Line Options

```bash
./.github/e2e-tests/scripts/run-validation.sh [OPTIONS]

Options:
  --mode MODE               Validation mode: quick | standard | comprehensive
  --resources TYPES         Comma-separated resource types to test
  --conventions CONVS       Comma-separated conventions: cafclassic,cafrandom,random,passthrough
  --terraform-version VER   Terraform version to use (e.g., 1.8.0)
  --parallel N              Number of parallel jobs (1-50)
  --timeout SECONDS         Total timeout in seconds
  --skip-cleanup            Skip cleanup (for debugging)
  --debug                   Enable debug output
  --output-dir DIR          Custom output directory
  --help                    Show help message

Examples:
  # Quick validation
  ./run-validation.sh --mode quick
  
  # Test specific resources
  ./run-validation.sh --mode quick --resources "azurerm_storage_account,azurerm_key_vault"
  
  # Debug failing test
  ./run-validation.sh --mode quick --debug --skip-cleanup
  
  # Custom parallelization
  ./run-validation.sh --mode standard --parallel 30
```

## Validation Modes

### Quick Mode (PR Context)

**Purpose**: Fast feedback for pull requests  
**Duration**: 5-8 minutes  
**Resources**: 20 representative types  
**Conventions**: cafclassic, cafrandom  
**TF Versions**: Min (1.5.7) and Max (1.8.0)  

**Resource Selection**:
- 5 compute resources (VM, VMSS, AKS, etc.)
- 5 storage resources (storage account, key vault, etc.)
- 5 networking resources (vnet, subnet, NSG, etc.)
- 5 database resources (SQL, PostgreSQL, Cosmos, etc.)

### Standard Mode (Feature Testing)

**Purpose**: Thorough testing of specific changes  
**Duration**: 10-15 minutes  
**Resources**: 100 types (major services)  
**Conventions**: All 4  
**TF Versions**: Min (1.5.7) and Max (1.8.0)  

**Resource Selection**:
- All compute resources
- All storage resources
- All networking resources
- All database resources
- Major application services

### Comprehensive Mode (Main Branch)

**Purpose**: Complete validation before release  
**Duration**: 25-30 minutes  
**Resources**: All 395 types  
**Conventions**: All 4  
**TF Versions**: All supported (1.5.x through 1.8.x)  

**Resource Coverage**:
- 100% of supported resource types
- All edge cases (min/max length, special chars)
- All Terraform version combinations

## Troubleshooting

### Issue: Provider Build Fails

**Symptoms**:
```
Error: compilation failed
go: error loading module dependencies
```

**Solution**:
```bash
# Update dependencies
go mod tidy
go mod download

# Verify Go version
go version  # Should be 1.24+

# Clean build cache
go clean -modcache
go build -v
```

### Issue: Terraform Init Fails

**Symptoms**:
```
Error: Failed to install provider
Could not find provider azurecaf
```

**Solution**:
```bash
# Verify provider binary exists
ls -lh /tmp/terraform-provider-azurecaf

# Check provider installation
export TF_LOG=DEBUG
terraform init

# Manually install provider
mkdir -p ~/.terraform.d/plugins/local/aztfmod/azurecaf/1.0.0/darwin_arm64/
cp terraform-provider-azurecaf ~/.terraform.d/plugins/local/aztfmod/azurecaf/1.0.0/darwin_arm64/
```

### Issue: Drift Detected

**Symptoms**:
```
Error: Drift detected
Resource name changed from 'st-dev-app-abc' to 'st-dev-app-xyz'
```

**Solution**:
```bash
# Check if random seed is stable
cat terraform.tfstate | jq '.resources[] | select(.type == "azurecaf_name") | .instances[].attributes.random_seed'

# Verify no timestamps in names
grep -r "timestamp" azurecaf/*.go

# Test random stability
./.github/e2e-tests/scripts/test-random-stability.sh
```

### Issue: Cleanup Fails

**Symptoms**:
```
Warning: Cleanup failed
Some resources could not be destroyed
```

**Solution**:
```bash
# Manual cleanup
cd /tmp/e2e-validation-{timestamp}
terraform destroy -auto-approve

# Force cleanup
rm -rf /tmp/e2e-validation-*

# For Azure resources (if applicable)
az group delete --name e2e-test-rg --yes --no-wait
```

### Issue: Timeout

**Symptoms**:
```
Error: Validation timed out after 30 minutes
```

**Solution**:
```bash
# Increase timeout
TIMEOUT=3600 ./.github/e2e-tests/scripts/run-validation.sh --mode comprehensive

# Run in smaller batches
./.github/e2e-tests/scripts/run-validation.sh --mode standard --parallel 10

# Check for hanging resources
ps aux | grep terraform
```

## CI/CD Integration

### GitHub Actions Workflows

Three workflows are available:

1. **e2e-validation-pr.yml**: Automatic PR validation
2. **e2e-validation-main.yml**: Automatic main branch validation
3. **e2e-validation-manual.yml**: Manual trigger with options

### Monitoring Results

```bash
# Via GitHub CLI
gh run list --workflow=e2e-validation-pr.yml --limit 10
gh run view {run-id} --log

# Via GitHub UI
# Navigate to: Actions → E2E Validation → Select run → View logs
```

### Artifacts

After each run, the following artifacts are available (7-day retention):

- **test-results.json**: Detailed test results
- **terraform-plans.tar.gz**: All plan outputs
- **build-logs.tar.gz**: Build and execution logs
- **drift-report.md**: Drift detection report

```bash
# Download artifacts
gh run download {run-id}

# Extract and view
tar -xzf terraform-plans.tar.gz
cat test-results.json | jq '.resource_test_results[] | select(.status == "fail")'
```

## Best Practices

### Before Committing

```bash
# 1. Run quick validation locally
./.github/e2e-tests/scripts/run-validation.sh --mode quick

# 2. If you changed specific resources, test them
./.github/e2e-tests/scripts/run-validation.sh --mode quick \
  --resources "azurerm_your_resource"

# 3. Verify no drift
# (Quick validation includes drift check automatically)
```

### During Code Review

```bash
# Reviewer: Verify E2E validation passed
gh pr checks

# If validation failed, download logs
gh run view {run-id} --log > e2e-failure.log

# Review failed tests
cat e2e-failure.log | grep "FAILED"
```

### Before Release

```bash
# Run comprehensive validation
./.github/e2e-tests/scripts/run-validation.sh --mode comprehensive

# Test all Terraform versions
for version in 1.5.7 1.6.6 1.7.5 1.8.0; do
  TF_VERSION=$version ./.github/e2e-tests/scripts/run-validation.sh --mode comprehensive
done

# Verify 100% pass rate
cat /tmp/e2e-validation-*/test-results.json | jq '.pass_rate'
```

## Performance Tips

### Speed Up Validation

1. **Use caching**:
```bash
# Cache Go modules
export GOCACHE=~/.cache/go-build
export GOMODCACHE=~/go/pkg/mod

# Cache Terraform binaries
export TF_PLUGIN_CACHE_DIR=~/.terraform.d/plugin-cache
```

2. **Increase parallelization**:
```bash
# Use more parallel jobs
./.github/e2e-tests/scripts/run-validation.sh --mode quick --parallel 20
```

3. **Filter resources**:
```bash
# Test only changed resources
CHANGED=$(git diff --name-only main... | grep resourceDefinition.json)
if [ -n "$CHANGED" ]; then
  # Extract changed resource types and test only those
  ./.github/e2e-tests/scripts/run-validation.sh --mode quick --resources "$CHANGED_TYPES"
fi
```

## Getting Help

### Documentation

- **Full Documentation**: `docs/E2E_VALIDATION.md`
- **Pipeline Contracts**: `specs/001-e2e-terraform-validation/contracts/`
- **Data Models**: `specs/001-e2e-terraform-validation/data-model.md`
- **Research**: `specs/001-e2e-terraform-validation/research.md`

### Support

- **GitHub Issues**: Report bugs or request features
- **GitHub Discussions**: Ask questions or share ideas
- **CI/CD Logs**: Check GitHub Actions logs for detailed error messages

### Examples

More examples are available in:
- `.github/e2e-tests/configs/`: Test configurations
- `.github/e2e-tests/scripts/`: Validation scripts
- `e2e/`: E2E testing framework

---

**Need more help?** See the [complete documentation](docs/E2E_VALIDATION.md) or [open an issue](https://github.com/aztfmod/terraform-provider-azurecaf/issues/new).
