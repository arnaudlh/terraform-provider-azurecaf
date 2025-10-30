# Data Model: End-to-End Terraform Validation Pipeline

**Date**: 2025-10-30  
**Feature**: 001-e2e-terraform-validation

## Overview

This document defines the data structures and models used in the E2E validation pipeline. These models represent test configurations, validation results, and pipeline state.

## Core Entities

### 1. ValidationConfiguration

Represents a complete E2E validation test configuration.

**Attributes**:
- `name`: string - Unique identifier for this configuration (e.g., "quick-pr", "comprehensive-main")
- `resource_types`: []string - List of Azure resource types to test (e.g., "azurerm_storage_account")
- `naming_conventions`: []string - Conventions to test ("cafclassic", "cafrandom", "random", "passthrough")
- `terraform_versions`: []string - Terraform versions to test against (e.g., "1.5.7", "1.8.0")
- `test_mode`: enum - "quick" | "standard" | "comprehensive"
- `max_parallel_jobs`: integer - Number of parallel test jobs (1-50)
- `timeout_minutes`: integer - Maximum execution time
- `cleanup_enabled`: boolean - Whether to clean up test resources

**Validation Rules**:
- `name` must be unique across configurations
- `resource_types` must contain at least 1 valid resource type
- `naming_conventions` must contain at least 1 valid convention
- `terraform_versions` must contain at least 1 supported version
- `max_parallel_jobs` must be between 1 and 50
- `timeout_minutes` must be between 1 and 180

**State Transitions**:
```
pending → running → completed | failed | timeout
```

**Example**:
```yaml
name: "quick-pr"
resource_types:
  - "azurerm_storage_account"
  - "azurerm_key_vault"
  - "azurerm_virtual_network"
naming_conventions:
  - "cafclassic"
  - "cafrandom"
terraform_versions:
  - "1.5.7"
  - "1.8.0"
test_mode: "quick"
max_parallel_jobs: 10
timeout_minutes: 15
cleanup_enabled: true
```

### 2. TestResource

Represents a single Azure resource type being tested.

**Attributes**:
- `resource_type`: string - Azure resource type (e.g., "azurerm_storage_account")
- `slug`: string - CAF slug for this resource (e.g., "st")
- `min_length`: integer - Minimum allowed name length
- `max_length`: integer - Maximum allowed name length
- `regex_pattern`: string - Validation regex pattern
- `test_inputs`: map[string]interface{} - Test input parameters
  - `name`: string - Base name for testing
  - `prefixes`: []string - Prefix values
  - `suffixes`: []string - Suffix values
  - `random_length`: integer - Random suffix length
  - `random_seed`: integer - Deterministic random seed

**Validation Rules**:
- `resource_type` must match pattern `azurerm_[a-z_]+`
- `slug` must be 2-10 characters
- `min_length` must be >= 1
- `max_length` must be > min_length and <= 255
- `regex_pattern` must be valid regex
- `test_inputs.name` must not be empty
- `random_seed` must be stable (for drift detection)

**Relationships**:
- One TestResource belongs to one ValidationConfiguration
- One TestResource maps to one entry in resourceDefinition.json

**Example**:
```yaml
resource_type: "azurerm_storage_account"
slug: "st"
min_length: 3
max_length: 24
regex_pattern: "^[a-z0-9]+$"
test_inputs:
  name: "testapp"
  prefixes: ["dev"]
  suffixes: ["001"]
  random_length: 3
  random_seed: 12345
```

### 3. NamingConventionTest

Represents a test of a specific naming convention.

**Attributes**:
- `convention`: enum - "cafclassic" | "cafrandom" | "random" | "passthrough"
- `resource_type`: string - Resource type being tested
- `configuration`: map[string]interface{} - Terraform configuration
- `expected_pattern`: string - Expected name pattern (regex)
- `expected_length_min`: integer - Expected minimum length
- `expected_length_max`: integer - Expected maximum length

**Validation Rules**:
- `convention` must be one of the 4 supported conventions
- `configuration` must be valid Terraform HCL
- `expected_pattern` must match generated names
- Generated name length must be within expected bounds

**State Transitions**:
```
initialized → planned → applied → verified | failed
```

**Example**:
```yaml
convention: "cafclassic"
resource_type: "azurerm_storage_account"
configuration: |
  data "azurecaf_name" "test" {
    name          = "testapp"
    resource_type = "azurerm_storage_account"
    prefixes      = ["dev"]
    suffixes      = ["001"]
  }
expected_pattern: "^st-dev-testapp-001$"
expected_length_min: 15
expected_length_max: 24
```

### 4. ValidationResult

Represents the result of an E2E validation run.

**Attributes**:
- `validation_id`: string - Unique identifier (UUID)
- `configuration_name`: string - Reference to ValidationConfiguration
- `status`: enum - "success" | "failure" | "error" | "timeout"
- `start_time`: timestamp - When validation started
- `end_time`: timestamp - When validation completed
- `duration_seconds`: integer - Total execution time
- `stages`: []StageResult - Results for each pipeline stage

**Validation Rules**:
- `validation_id` must be unique
- `status` must be one of the defined enum values
- `end_time` must be after `start_time`
- `duration_seconds` must match end_time - start_time

**Relationships**:
- One ValidationResult has many StageResults
- One ValidationResult has many ResourceTestResults

**Example**:
```yaml
validation_id: "550e8400-e29b-41d4-a716-446655440000"
configuration_name: "quick-pr"
status: "success"
start_time: "2025-10-30T10:00:00Z"
end_time: "2025-10-30T10:07:30Z"
duration_seconds: 450
stages:
  - build
  - plan
  - apply
  - drift_check
```

### 5. StageResult

Represents the result of a single pipeline stage.

**Attributes**:
- `stage_name`: enum - "build" | "plan" | "apply" | "drift_check" | "cleanup"
- `status`: enum - "success" | "failure" | "skipped"
- `start_time`: timestamp - When stage started
- `end_time`: timestamp - When stage completed
- `duration_seconds`: integer - Stage execution time
- `output`: string - Stage output logs
- `error_message`: string - Error message if failed
- `artifacts`: []string - Paths to generated artifacts

**Validation Rules**:
- `stage_name` must be one of the defined pipeline stages
- `status` must be success, failure, or skipped
- If status is failure, `error_message` must not be empty
- `artifacts` must contain valid file paths

**State Transitions**:
```
pending → running → success | failure | skipped
```

**Example**:
```yaml
stage_name: "plan"
status: "success"
start_time: "2025-10-30T10:01:00Z"
end_time: "2025-10-30T10:03:30Z"
duration_seconds: 150
output: "Terraform plan output..."
artifacts:
  - "/tmp/plan-output.txt"
  - "/tmp/terraform.tfplan"
```

### 6. ResourceTestResult

Represents the result of testing a specific resource type.

**Attributes**:
- `resource_type`: string - Azure resource type tested
- `naming_convention`: string - Convention used
- `terraform_version`: string - Terraform version used
- `status`: enum - "pass" | "fail" | "skip"
- `generated_name`: string - Name generated by provider
- `name_length`: integer - Length of generated name
- `validation_passed`: boolean - Whether name passed validation rules
- `plan_succeeded`: boolean - Whether terraform plan succeeded
- `apply_succeeded`: boolean - Whether terraform apply succeeded
- `drift_detected`: boolean - Whether drift was detected on second apply
- `error_message`: string - Error details if failed
- `duration_seconds`: integer - Test execution time

**Validation Rules**:
- `resource_type` must be valid Azure resource type
- `naming_convention` must be one of the 4 supported conventions
- `terraform_version` must match semantic versioning pattern
- If `status` is fail, `error_message` must not be empty
- `generated_name` length must match `name_length`
- If `drift_detected` is true, `status` must be fail

**Metrics**:
- Pass rate per resource type
- Pass rate per naming convention
- Pass rate per Terraform version
- Average test duration per resource type

**Example**:
```yaml
resource_type: "azurerm_storage_account"
naming_convention: "cafclassic"
terraform_version: "1.8.0"
status: "pass"
generated_name: "st-dev-testapp-001"
name_length: 18
validation_passed: true
plan_succeeded: true
apply_succeeded: true
drift_detected: false
duration_seconds: 5
```

### 7. TerraformState

Represents Terraform state for testing.

**Attributes**:
- `version`: integer - Terraform state version
- `terraform_version`: string - Terraform version that created state
- `serial`: integer - State serial number
- `lineage`: string - State lineage identifier
- `resources`: []TerraformResource - Resources in state
- `outputs`: map[string]interface{} - Terraform outputs

**Validation Rules**:
- `version` must be 4 (current Terraform state version)
- `terraform_version` must match semantic versioning
- `serial` increments on each state change
- `lineage` must be stable across applies

**Purpose**:
- Track provider-generated names across applies
- Validate state stability (no drift)
- Verify random seeds are persisted correctly

**Example**:
```json
{
  "version": 4,
  "terraform_version": "1.8.0",
  "serial": 2,
  "lineage": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "resources": [
    {
      "type": "azurecaf_name",
      "name": "test",
      "provider": "provider[\"registry.terraform.io/aztfmod/azurecaf\"]",
      "instances": [
        {
          "attributes": {
            "id": "st-dev-testapp-001",
            "result": "st-dev-testapp-001",
            "random_seed": 12345
          }
        }
      ]
    }
  ]
}
```

## Data Relationships

```
ValidationConfiguration (1) ──has──> (many) TestResource
ValidationConfiguration (1) ──has──> (many) NamingConventionTest
ValidationConfiguration (1) ──produces──> (1) ValidationResult
ValidationResult (1) ──contains──> (many) StageResult
ValidationResult (1) ──contains──> (many) ResourceTestResult
TestResource (1) ──tested_by──> (many) NamingConventionTest
NamingConventionTest (1) ──produces──> (1) TerraformState
ResourceTestResult (many) ──references──> (1) TestResource
```

## Data Storage

### Pipeline Execution (Runtime)

**Storage**: GitHub Actions workflow state and artifacts

- ValidationConfiguration: YAML in `.github/e2e-tests/configs/`
- TestResource: Generated from resourceDefinition.json
- TerraformState: Temporary local backend, deleted after test
- StageResult: Captured in GitHub Actions logs
- Artifacts: Stored in GitHub Actions artifacts (7-day retention)

### Results (Persistent)

**Storage**: GitHub Actions run history and artifacts

- ValidationResult: Derived from workflow run status
- ResourceTestResult: Extracted from test output logs
- Performance metrics: Tracked via GitHub Actions duration
- Trend data: Optional external analytics (future enhancement)

## Data Validation

### Pre-Execution Validation

Before running E2E validation:
1. Validate ValidationConfiguration against schema
2. Verify all resource_types exist in resourceDefinition.json
3. Confirm naming_conventions are supported
4. Check Terraform versions are available
5. Validate test input parameters

### Runtime Validation

During E2E execution:
1. Validate generated names match expected patterns
2. Verify name lengths are within constraints
3. Confirm Terraform state is stable (no drift)
4. Check random seeds are deterministic
5. Validate cleanup completes successfully

### Post-Execution Validation

After E2E completion:
1. Verify all tests have results (no orphans)
2. Confirm stage order is correct (build → plan → apply → drift)
3. Validate duration metrics are reasonable
4. Check artifact files were created
5. Ensure cleanup removed all test resources

## Performance Considerations

### Data Volume

**Per PR Validation** (quick mode):
- 20 resource types × 2 conventions × 2 TF versions = 80 tests
- ~5 seconds per test = 400 seconds / 10 parallel = 40 seconds
- Logs: ~100 KB per test = 8 MB total
- Artifacts: ~500 KB per test = 40 MB total

**Per Main Validation** (comprehensive):
- 395 resource types × 4 conventions × 4 TF versions = 6,320 tests  
- ~5 seconds per test = 31,600 seconds / 50 parallel = 632 seconds (10.5 min)
- Logs: ~100 KB per test = 632 MB total
- Artifacts: ~500 KB per test = 3.2 GB total

### Optimization Strategies

1. **Parallel Execution**: Run up to 50 tests simultaneously
2. **Caching**: Cache provider builds and Terraform binaries
3. **Smart Selection**: Run subset on PRs, full suite on main
4. **Artifact Retention**: 7 days (GitHub default), compress logs
5. **State Cleanup**: Delete temporary state files immediately after test

## Schema Definitions

### ValidationConfiguration Schema

```yaml
type: object
required:
  - name
  - resource_types
  - naming_conventions
  - terraform_versions
  - test_mode
properties:
  name:
    type: string
    pattern: "^[a-z0-9-]+$"
  resource_types:
    type: array
    minItems: 1
    items:
      type: string
      pattern: "^azurerm_[a-z_]+$"
  naming_conventions:
    type: array
    minItems: 1
    items:
      type: string
      enum: ["cafclassic", "cafrandom", "random", "passthrough"]
  terraform_versions:
    type: array
    minItems: 1
    items:
      type: string
      pattern: "^[0-9]+\\.[0-9]+\\.[0-9]+$"
  test_mode:
    type: string
    enum: ["quick", "standard", "comprehensive"]
  max_parallel_jobs:
    type: integer
    minimum: 1
    maximum: 50
  timeout_minutes:
    type: integer
    minimum: 1
    maximum: 180
  cleanup_enabled:
    type: boolean
```

### ResourceTestResult Schema

```yaml
type: object
required:
  - resource_type
  - naming_convention
  - terraform_version
  - status
properties:
  resource_type:
    type: string
  naming_convention:
    type: string
    enum: ["cafclassic", "cafrandom", "random", "passthrough"]
  terraform_version:
    type: string
  status:
    type: string
    enum: ["pass", "fail", "skip"]
  generated_name:
    type: string
  name_length:
    type: integer
  validation_passed:
    type: boolean
  plan_succeeded:
    type: boolean
  apply_succeeded:
    type: boolean
  drift_detected:
    type: boolean
  error_message:
    type: string
  duration_seconds:
    type: integer
    minimum: 0
```

## Usage Examples

### Creating a Quick Validation Configuration

```yaml
# .github/e2e-tests/configs/quick-pr.yaml
name: "quick-pr"
resource_types:
  - "azurerm_storage_account"
  - "azurerm_key_vault"
  - "azurerm_resource_group"
  - "azurerm_virtual_network"
  - "azurerm_subnet"
naming_conventions:
  - "cafclassic"
  - "cafrandom"
terraform_versions:
  - "1.5.7"
  - "1.8.0"
test_mode: "quick"
max_parallel_jobs: 10
timeout_minutes: 15
cleanup_enabled: true
```

### Querying Test Results

```bash
# Get all failed tests from last run
cat validation-result.json | jq '.resource_test_results[] | select(.status == "fail")'

# Get pass rate per naming convention
cat validation-result.json | jq '
  .resource_test_results 
  | group_by(.naming_convention) 
  | map({
      convention: .[0].naming_convention, 
      pass_rate: (map(select(.status == "pass")) | length) / length
    })'

# Find slowest resource types
cat validation-result.json | jq '
  .resource_test_results 
  | sort_by(.duration_seconds) 
  | reverse 
  | .[0:10]'
```
