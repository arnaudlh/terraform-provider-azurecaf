# End-to-End Terraform Validation

**Status**: In Development  
**Version**: 1.0.0  
**Last Updated**: 2025-10-30

## Overview

The E2E Terraform validation pipeline provides comprehensive automated testing of the azurecaf provider by building the provider binary, executing `terraform plan` and `terraform apply`, and performing drift detection through a second apply. This ensures the provider generates Azure-compliant resource names across all supported naming conventions and resource types.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Validation Modes](#validation-modes)
- [CI/CD Integration](#cicd-integration)
- [Local Development](#local-development)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Performance](#performance)
- [Contributing](#contributing)

## Architecture

### Pipeline Stages

The E2E validation pipeline consists of five sequential stages:

1. **Build**: Compile the provider binary from source
2. **Plan**: Execute `terraform init` and `terraform plan`
3. **Apply**: Execute `terraform apply` to validate names
4. **Drift Check**: Run a second apply to detect drift
5. **Cleanup**: Remove test resources and temporary files

### Validation Modes

| Mode | Duration | Resources | Conventions | Use Case |
|------|----------|-----------|-------------|----------|
| **Quick** | 5-8 min | 20 representative | cafclassic, cafrandom | PR validation |
| **Standard** | 10-15 min | 100 major services | All 4 | Feature testing |
| **Comprehensive** | 25-30 min | All 395 types | All 4 | Main branch, releases |

### Technology Stack

- **Language**: Bash scripts for orchestration
- **Terraform**: Multiple versions (1.5.x through 1.8.x)
- **Provider**: Built from source for each run
- **CI/CD**: GitHub Actions workflows
- **State**: Local backend (ephemeral)

## Quick Start

### Prerequisites

- Go 1.24+
- Terraform CLI (1.5.x - 1.8.x)
- Git
- Bash shell

### Run Quick Validation Locally

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

### Run in CI/CD

E2E validation runs automatically:

- **On Pull Requests**: Quick validation (~7-8 minutes)
- **On Main Branch**: Comprehensive validation (~25-30 minutes)
- **Manual Trigger**: Configurable mode and resource selection

See [CI/CD Integration](#cicd-integration) for details.

## Validation Modes

### Quick Mode

**Purpose**: Fast feedback for pull requests

**Configuration**:
- **Resources**: 20 representative types across categories
  - 5 compute (VM, VMSS, AKS, etc.)
  - 5 storage (storage account, key vault, etc.)
  - 5 networking (vnet, subnet, NSG, etc.)
  - 5 database (SQL, PostgreSQL, Cosmos, etc.)
- **Conventions**: cafclassic, cafrandom
- **Terraform Versions**: Min (1.5.7) and Max (1.8.0)
- **Parallelization**: 10 concurrent jobs

**Usage**:
```bash
./.github/e2e-tests/scripts/run-validation.sh --mode quick
```

### Standard Mode

**Purpose**: Thorough testing of specific changes

**Configuration**:
- **Resources**: 100 types covering major Azure services
- **Conventions**: All 4 (cafclassic, cafrandom, random, passthrough)
- **Terraform Versions**: Min and Max
- **Parallelization**: 20 concurrent jobs

**Usage**:
```bash
./.github/e2e-tests/scripts/run-validation.sh --mode standard
```

### Comprehensive Mode

**Purpose**: Complete validation before releases

**Configuration**:
- **Resources**: All 395 supported resource types
- **Conventions**: All 4
- **Terraform Versions**: All supported (1.5.x through 1.8.x)
- **Parallelization**: 50 concurrent jobs

**Usage**:
```bash
./.github/e2e-tests/scripts/run-validation.sh --mode comprehensive
```

## CI/CD Integration

### GitHub Actions Workflows

Three workflows provide automated validation:

#### 1. PR Validation (`e2e-validation-pr.yml`)

**Trigger**: Pull request opened, updated, synchronized  
**Mode**: Quick  
**Duration**: ~7-8 minutes  
**Purpose**: Fast feedback for developers

**Workflow Steps**:
1. Checkout code
2. Setup Go with module caching
3. Build provider binary
4. Setup Terraform
5. Run quick validation
6. Upload artifacts (logs, results)

#### 2. Main Branch Validation (`e2e-validation-main.yml`)

**Trigger**: Push to main branch  
**Mode**: Comprehensive  
**Duration**: ~25-30 minutes  
**Purpose**: Complete validation before releases

**Workflow Steps**:
1. Checkout code
2. Setup Go with module caching
3. Build provider binary
4. Setup Terraform (matrix: all versions)
5. Run comprehensive validation
6. Upload artifacts
7. Publish test report

#### 3. Manual Validation (`e2e-validation-manual.yml`)

**Trigger**: Manual workflow dispatch  
**Mode**: Configurable  
**Duration**: Varies  
**Purpose**: Ad-hoc testing and debugging

**Inputs**:
- `mode`: quick | standard | comprehensive
- `resource_types`: Optional filter (comma-separated)
- `conventions`: Optional filter (comma-separated)
- `terraform_versions`: Optional version list

**Usage**:
```bash
gh workflow run e2e-validation-manual.yml \
  --field mode=standard \
  --field resource_types=azurerm_storage_account,azurerm_key_vault
```

## Local Development

### Setup

1. **Clone repository**:
   ```bash
   git clone https://github.com/aztfmod/terraform-provider-azurecaf.git
   cd terraform-provider-azurecaf
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Verify Terraform installation**:
   ```bash
   terraform version
   ```

### Running Tests

#### Test Specific Resource Types

```bash
./.github/e2e-tests/scripts/run-validation.sh --mode quick \
  --resources "azurerm_storage_account,azurerm_key_vault"
```

#### Test Specific Naming Convention

```bash
./.github/e2e-tests/scripts/run-validation.sh --mode standard \
  --conventions "cafclassic"
```

#### Test Specific Terraform Version

```bash
TF_VERSION=1.8.0 ./.github/e2e-tests/scripts/run-validation.sh --mode quick
```

#### Debug Mode

```bash
DEBUG=1 ./.github/e2e-tests/scripts/run-validation.sh --mode quick \
  --skip-cleanup
```

### Development Workflow

1. **Make changes** to provider code
2. **Run quick validation** locally:
   ```bash
   ./.github/e2e-tests/scripts/run-validation.sh --mode quick
   ```
3. **Commit changes** if validation passes
4. **Create pull request** - CI will run quick validation
5. **Merge to main** - CI will run comprehensive validation

## Configuration

### Test Configuration Files

Test configurations are stored in `.github/e2e-tests/configs/`:

```
configs/
├── quick/                    # Quick validation configs
│   ├── cafclassic.tf        # CAF classic naming tests
│   ├── cafrandom.tf         # CAF random naming tests
│   ├── random.tf            # Random naming tests
│   ├── passthrough.tf       # Passthrough naming tests
│   └── provider.tf          # Provider configuration
└── comprehensive/            # Comprehensive validation configs
    ├── compute.tf           # Compute resources
    ├── storage.tf           # Storage resources
    ├── networking.tf        # Network resources
    ├── databases.tf         # Database resources
    └── ...                  # Other categories
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TF_VERSION` | 1.8.0 | Terraform version to use |
| `DEBUG` | false | Enable debug output |
| `MAX_PARALLEL` | 10 | Maximum parallel jobs |
| `TIMEOUT` | 1800 | Total timeout (seconds) |
| `SKIP_CLEANUP` | false | Skip cleanup (debugging) |

### Script Options

```bash
./.github/e2e-tests/scripts/run-validation.sh [OPTIONS]

Options:
  --mode MODE               Validation mode: quick | standard | comprehensive
  --resources TYPES         Comma-separated resource types to test
  --conventions CONVS       Comma-separated conventions to test
  --terraform-version VER   Terraform version (e.g., 1.8.0)
  --parallel N              Number of parallel jobs (1-50)
  --timeout SECONDS         Total timeout in seconds
  --skip-cleanup            Skip cleanup (for debugging)
  --debug                   Enable debug output
  --output-dir DIR          Custom output directory
  --help                    Show help message
```

## Troubleshooting

### Common Issues

#### Provider Build Fails

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
```

#### Terraform Init Fails

**Symptoms**:
```
Error: Failed to install provider
Could not find provider azurecaf
```

**Solution**:
```bash
# Verify provider binary exists
ls -lh /tmp/terraform-provider-azurecaf

# Enable debug mode
export TF_LOG=DEBUG
terraform init
```

#### Drift Detected

**Symptoms**:
```
Error: Drift detected
Resource name changed
```

**Solution**:
```bash
# Check random seed stability
cat terraform.tfstate | jq '.resources[] | select(.type == "azurecaf_name") | .instances[].attributes.random_seed'

# Verify no timestamps in names
grep -r "timestamp" azurecaf/*.go
```

#### Cleanup Fails

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
```

### Getting Help

- **Documentation**: See [quickstart guide](../specs/001-e2e-terraform-validation/quickstart.md)
- **GitHub Issues**: Report bugs or request features
- **CI Logs**: Check GitHub Actions logs for detailed errors

## Performance

### Performance Targets

| Metric | Target | Current |
|--------|--------|---------|
| PR validation | < 10 min | TBD |
| Main validation | < 30 min | TBD |
| Provider build (cached) | < 1 min | TBD |
| Test per resource | < 5 sec | TBD |

### Optimization Strategies

1. **Caching**: Go modules, Terraform binaries
2. **Parallelization**: Up to 50 concurrent jobs
3. **Smart selection**: Subset testing for PRs
4. **Early termination**: Fail fast on errors

### Monitoring

Track performance metrics via GitHub Actions:
- Build duration
- Test execution time per resource
- Cache hit rates
- Time to first failure

## Contributing

### Adding New Test Resources

1. Update test configuration in `.github/e2e-tests/configs/`
2. Follow naming convention patterns
3. Include all required attributes
4. Test locally before committing

### Modifying Validation Logic

1. Update scripts in `.github/e2e-tests/scripts/`
2. Maintain backward compatibility
3. Update documentation
4. Run full validation suite

### Reporting Issues

When reporting E2E validation issues, include:
- Validation mode (quick/standard/comprehensive)
- Terraform version
- Provider version
- Full error logs
- Steps to reproduce

## References

- [Feature Specification](../specs/001-e2e-terraform-validation/spec.md)
- [Implementation Plan](../specs/001-e2e-terraform-validation/plan.md)
- [Quick Start Guide](../specs/001-e2e-terraform-validation/quickstart.md)
- [Pipeline Contracts](../specs/001-e2e-terraform-validation/contracts/validation-pipeline.md)
- [Data Models](../specs/001-e2e-terraform-validation/data-model.md)

---

**Status**: This documentation is actively being developed as part of feature 001-e2e-terraform-validation.
