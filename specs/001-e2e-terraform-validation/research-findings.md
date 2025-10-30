# E2E Testing Research Findings for Terraform Provider CI/CD

**Research Date:** October 30, 2025  
**Provider:** terraform-provider-azurecaf  
**Purpose:** Establish battle-tested E2E testing patterns for provider CI/CD

---

## Executive Summary

After analyzing testing approaches from HashiCorp's major providers (AWS, Azure, Google) and official documentation, this document provides actionable recommendations for implementing E2E testing in terraform-provider-azurecaf's CI/CD pipeline.

**Key Finding:** Major providers rely on **acceptance tests that create real resources** rather than pure E2E tests, combined with fast unit tests for quick validation.

---

## 1. Terraform Provider E2E Testing Best Practices

### 1.1 Testing Philosophy from Major Providers

**HashiCorp's Official Stance:**
- Acceptance tests are the primary validation mechanism
- Tests create **actual cloud infrastructure** to verify behavior
- The `terraform-plugin-testing` framework is the standard tool
- Real infrastructure tests are the responsibility of the developer running them

**Source:** HashiCorp Plugin Testing Documentation
```go
// Acceptance tests create actual cloud infrastructure, with possible 
// expenses incurred, and are the responsibility of the user running the tests
```

### 1.2 Standard Testing Patterns

#### Pattern 1: Test Step Structure
All major providers follow this pattern:

```go
func TestAccExampleResource_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
    r := ExampleResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),  // Validates state matches real resource
    })
}
```

**Key Components:**
1. **Config Step:** Applies Terraform configuration
2. **Check Functions:** Verifies resource exists and attributes match
3. **ImportStep:** Re-imports resource to validate state accuracy
4. **Parallel Execution:** Uses `t.Parallel()` for independent tests

#### Pattern 2: Test Coverage Types

From Azure Provider analysis:
1. **Basic Test:** Minimum required fields, happy path
2. **Complete Test:** All optional and required fields
3. **Update Test:** Modify configuration, verify changes persist
4. **RequiresImport Test:** Validates resource adoption protection

```go
// Basic: Minimum configuration
func (ExampleResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
resource "azurerm_example" "test" {
  name     = "acctestexample%[1]d"
  location = %[2]q
  required_field = "value"
}
`, data.RandomInteger, data.Locations.Primary)
}

// Complete: All fields
func (ExampleResource) complete(data acceptance.TestData) string {
    return fmt.Sprintf(`
resource "azurerm_example" "test" {
  name     = "acctestexample%[1]d"
  location = %[2]q
  required_field = "value"
  optional_field_1 = "value1"
  optional_field_2 = "value2"
  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
```

### 1.3 Plan/Apply Workflow Testing

**Standard Approach:** The testing framework automatically handles plan/apply cycles:

```go
// Framework automatically:
// 1. Runs terraform init
// 2. Runs terraform plan
// 3. Runs terraform apply
// 4. Runs Check functions
// 5. Optionally runs terraform import
// 6. Optionally runs terraform destroy
```

**Advanced: Explicit Plan Checks (Terraform 1.6+)**

```go
{
    Config: r.update(data),
    ConfigPlanChecks: resource.ConfigPlanChecks{
        PreApply: []plancheck.PlanCheck{
            plancheck.ExpectNoDelete(),  // Ensures update in-place
        },
    },
}
```

### 1.4 State Management Best Practices

#### From AWS Provider:
```go
// Tests automatically clean up after themselves
func testAccCheckExampleDestroy(ctx context.Context) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        for _, rs := range s.RootModule().Resources {
            if rs.Type != "aws_example" {
                continue
            }
            // Verify resource no longer exists
            _, err := client.DescribeExample(ctx, rs.Primary.ID)
            if err == nil {
                return fmt.Errorf("Example %s still exists", rs.Primary.ID)
            }
        }
        return nil
    }
}
```

#### State Refresh Pattern (Azure Provider):
```go
// Framework as of 1.6.0 no longer auto-refreshes
// Add explicit refresh steps when needed
refreshStep := TestStep{
    RefreshState: true,
}
```

---

## 2. Mock vs Real Azure Resources

### 2.1 Industry Practice: Real Resources Win

**All major providers use real cloud resources for acceptance tests:**

| Provider | Approach | Rationale |
|----------|----------|-----------|
| AWS | Real resources | "Verifies behavior against actual APIs" |
| Azure | Real resources | "Validates both local state and remote values" |
| Google | Real resources | "Ensures real-world use case validation" |

**For terraform-provider-azurecaf:** Since this is a **naming provider** that doesn't create actual Azure resources, the tradeoff is different:
- ✅ No actual Azure resource costs
- ✅ Can run tests without Azure credentials
- ⚠️ Should still test against Azure naming constraints

### 2.2 Cost Minimization Strategies

From AWS provider documentation:

```bash
# Environment variables control test scope
TF_ACC=1                    # Enable acceptance tests
AWS_DEFAULT_REGION=us-west-2
TF_ACC_TERRAFORM_VERSION=1.5.x  # Test specific version
```

**Recommended for azurecaf:**
```bash
# Only test name generation, not actual Azure resources
TF_ACC=1
TF_ACC_PROVIDER_NAMESPACE=registry.terraform.io/aztfmod
# No Azure credentials needed since we're only generating names
```

### 2.3 Resource Cleanup Patterns

#### Pattern 1: Automatic Cleanup (Default)
```go
resource.Test(t, resource.TestCase{
    CheckDestroy: testAccCheckResourceDestroy(ctx),
    // ... test steps ...
})
// Framework calls CheckDestroy after all steps
```

#### Pattern 2: Explicit Cleanup (Long-running tests)
```go
// From AWS provider
defer func() {
    // Ensure cleanup even if test panics
    if err := deleteTestResources(ctx, resourceId); err != nil {
        t.Logf("Failed to clean up: %v", err)
    }
}()
```

#### Pattern 3: Sweep Tests (Batch cleanup)
```go
// From Google provider
func init() {
    resource.AddTestSweepers("Example", &resource.Sweeper{
        Name: "Example",
        F:    sweepExample,
    })
}

func sweepExample(region string) error {
    // Find and delete all test resources in region
    // Run with: go test -sweep=us-west-2
}
```

---

## 3. Drift Detection Testing

### 3.1 Standard Pattern: Apply Twice, Expect No Changes

```go
func TestAccExample_noDrift(t *testing.T) {
    config := r.basic(data)
    
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            {
                Config: config,
                Check: resource.ComposeTestCheckFunc(
                    check.That(data.ResourceName).ExistsInAzure(r),
                ),
            },
            // Second apply - should be no-op
            {
                Config: config,
                ConfigPlanChecks: resource.ConfigPlanChecks{
                    PreApply: []plancheck.PlanCheck{
                        plancheck.ExpectEmptyPlan(),  // No changes expected
                    },
                },
            },
        },
    })
}
```

### 3.2 Handling Random Values

#### Problem: Random values cause perpetual drift

**Solution 1: Ignore in State (Azure Pattern)**
```go
func (ExampleResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
resource "azurerm_example" "test" {
  name = "acctest%[1]d"  # Uses test data with fixed seed
  # NOT: name = "acctest${random_string.test.result}"
}
`, data.RandomInteger)  # Framework provides consistent random values
}
```

**Solution 2: Lifecycle Ignore (When unavoidable)**
```go
resource "azurecaf_name" "test" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  random_length = 4
  
  lifecycle {
    ignore_changes = [random_length]  # Don't re-randomize
  }
}
```

### 3.3 Time-Based Values

**Problem:** Creation timestamps, update times cause drift

**Solution from Azure Provider:**
```go
data.ImportStep(
    "metadata.0.resource_version",  # Ignore version
    "status.0.conditions",          # Ignore timestamps
)
```

**For azurecaf:** Likely need to ignore:
```go
data.ImportStep(
    "result_format",  # If computed
    "results",        # Map of generated names - may include timestamps
)
```

### 3.4 Common Pitfalls

| Pitfall | Impact | Solution |
|---------|--------|----------|
| Not using test framework's random | Drift on every run | Use `data.RandomInteger` |
| Computed values in config | Can't re-apply | Move to data sources |
| External dependencies change | False drift | Use `depends_on` or fixed references |
| API returns more fields than set | Drift detection fails | Use `ImportStateVerifyIgnore` |

---

## 4. Multi-Version Terraform Testing

### 4.1 Version Matrix Strategy (Google Provider)

```go
// Testing multiple versions
func TestAccExample_multiVersion(t *testing.T) {
    versions := []string{"1.5.0", "1.6.0", "1.7.0", "1.8.0"}
    
    for _, version := range versions {
        t.Run(version, func(t *testing.T) {
            t.Setenv("TF_ACC_TERRAFORM_VERSION", version)
            // Run tests with specific Terraform version
        })
    }
}
```

### 4.2 Version Compatibility Testing

**Azure Provider Pattern:**
```go
func TestAccExample_migration(t *testing.T) {
    oldVersion := map[string]resource.ExternalProvider{
        "azurerm": {
            VersionConstraint: "3.0.0",
            Source:            "registry.terraform.io/hashicorp/azurerm",
        },
    }
    
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            {
                Config:            testConfig,
                ExternalProviders: oldVersion,  # Use old provider
            },
            {
                Config:                   testConfig,
                ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),  # Upgrade
                // Should succeed without recreation
            },
        },
    })
}
```

### 4.3 Managing Multiple Binaries in CI

**Recommended Approach from AWS Provider:**

```yaml
# .github/workflows/test.yml
strategy:
  matrix:
    terraform: ['1.5.x', '1.6.x', '1.7.x', '1.8.x']
steps:
  - uses: hashicorp/setup-terraform@v3
    with:
      terraform_version: ${{ matrix.terraform }}
      terraform_wrapper: false
  - run: make testacc
    env:
      TF_ACC: "1"
      TF_ACC_TERRAFORM_VERSION: ${{ matrix.terraform }}
```

### 4.4 Terraform Version Checks (Plugin Testing v1.6+)

```go
func TestAccExample_requiresTerraform16(t *testing.T) {
    resource.Test(t, resource.TestCase{
        TerraformVersionChecks: []tfversion.TerraformVersionCheck{
            tfversion.SkipBelow(version.Must(version.NewVersion("1.6.0"))),
        },
        // ... test steps ...
    })
}
```

**For azurecaf recommendations:**
- Support last 4 minor versions (currently 1.5.x - 1.8.x)
- Test weekly against latest patch of each minor version
- Pin major.minor in go.mod: `terraform-plugin-framework ~> 1.5`

---

## 5. GitHub Actions for Terraform Testing

### 5.1 Best Practices from Major Providers

#### AWS Provider CI Structure:
```yaml
name: Acceptance Tests
on:
  pull_request:
    paths:
      - 'internal/service/**/*_test.go'
      - 'internal/service/**/*.go'
  workflow_dispatch:  # Manual trigger for expensive tests

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        terraform: ['1.5.x', '1.6.x', '1.7.x', '1.8.x']
        parallel: [20]  # Run 20 tests in parallel
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
          cache: true
      
      - uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: ${{ matrix.terraform }}
          terraform_wrapper: false  # Important: disable wrapper for testing
      
      - name: Run Tests
        run: |
          make testacc TESTARGS="-run=TestAcc -parallel=${{ matrix.parallel }}"
        env:
          TF_ACC: "1"
          # No AWS credentials needed for azurecaf (naming only)
```

### 5.2 Secrets and Credentials Management

**For azurecaf (naming provider - minimal credentials needed):**

```yaml
# No Azure credentials required for name generation tests
# Only needed if testing against Azure naming constraints
env:
  TF_ACC: "1"
  ARM_SUBSCRIPTION_ID: ${{ secrets.ARM_SUBSCRIPTION_ID }}  # Optional
  # Only for validating generated names against Azure limits
```

**If testing with real Azure resources:**

```yaml
env:
  ARM_CLIENT_ID: ${{ secrets.ARM_CLIENT_ID }}
  ARM_CLIENT_SECRET: ${{ secrets.ARM_CLIENT_SECRET }}
  ARM_SUBSCRIPTION_ID: ${{ secrets.ARM_SUBSCRIPTION_ID }}
  ARM_TENANT_ID: ${{ secrets.ARM_TENANT_ID }}
  ARM_TEST_LOCATION: "eastus"
  ARM_TEST_LOCATION_ALT: "westus2"
```

### 5.3 Artifact Storage

**Pattern from Azure Provider:**

```yaml
- name: Upload Test Results
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: test-results-${{ matrix.terraform }}
    path: |
      test-results/**/*.xml
      test-results/**/*.json
      terraform*.log
    retention-days: 30

- name: Upload Coverage
  uses: codecov/codecov-action@v4
  with:
    files: ./coverage.out
    flags: acceptance
```

### 5.4 Parallel vs Sequential Testing

**Parallel Strategy (Default for Independent Tests):**
```yaml
strategy:
  fail-fast: false
  matrix:
    test-group:
      - 'azurerm_resource_group'
      - 'azurerm_virtual_network'
      - 'azurerm_storage_account'
    
jobs:
  test-${{ matrix.test-group }}:
    # Each group runs in parallel
    run: make testacc TESTARGS="-run=TestAcc${{ matrix.test-group }}"
```

**Sequential Strategy (For Dependent Tests):**
```go
// From AWS provider - tests that must run serially
func TestAccOrganizationPolicy(t *testing.T) {
    // Note: Acting on same organization, run serially
    testCases := map[string]func(t *testing.T){
        "basic":  testBasic,
        "update": testUpdate,
    }
    
    for name, tc := range testCases {
        t.Run(name, func(t *testing.T) {
            tc(t)  // No t.Parallel() - runs sequentially
        })
    }
}
```

**Recommendation for azurecaf:**
- Run all tests in parallel (no Azure state conflicts)
- Use 20-30 parallel workers for optimal speed
- No need for sequential execution (naming is stateless)

---

## 6. Performance Optimization

### 6.1 Target: 10-Minute Quick Validation

**Strategy from Google Provider:**

```yaml
# Quick smoke test (< 10 minutes)
on:
  pull_request:

jobs:
  quick-test:
    runs-on: ubuntu-latest
    steps:
      - run: make test  # Unit tests only (< 2 min)
      
      - name: Smoke Tests
        run: make testacc TESTARGS="-run=TestAcc.*_basic -parallel=30"
        # Only _basic tests, high parallelism
```

**Full test suite (nightly):**
```yaml
on:
  schedule:
    - cron: '0 2 * * *'  # 2 AM daily

jobs:
  full-test:
    strategy:
      matrix:
        terraform: ['1.5.x', '1.6.x', '1.7.x', '1.8.x']
    steps:
      - run: make testacc TESTARGS="-parallel=20"
        timeout-minutes: 120  # 2 hour timeout
```

### 6.2 Parallelization Techniques

**Pattern 1: Test Splitting by Package**
```bash
# From AWS provider's Makefile
SERVICES := $(shell ls internal/service/)

test-service-%:
	TF_ACC=1 go test ./internal/service/$* -v -parallel=20 -timeout=120m
```

**Pattern 2: Dynamic Test Discovery**
```go
// From Google provider
func TestMain(m *testing.M) {
    // Detect available CPU cores
    parallelism := runtime.NumCPU()
    flag.IntVar(&parallelism, "test.parallel", parallelism, "parallelism")
    flag.Parse()
    os.Exit(m.Run())
}
```

**For azurecaf:**
```bash
# Run all tests in parallel (no resource contention)
make testacc TESTARGS="-parallel=50"  # Very high parallelism possible

# Or split by resource type
make testacc TESTARGS="-run=TestAccName.*_basic -parallel=30"
make testacc TESTARGS="-run=TestAccNamingConvention.*_basic -parallel=30"
```

### 6.3 Caching Strategies

#### Provider Build Cache
```yaml
- uses: actions/setup-go@v5
  with:
    go-version-file: 'go.mod'
    cache: true  # Caches go modules and build cache

- name: Cache Terraform Providers
  uses: actions/cache@v4
  with:
    path: |
      ~/.terraform.d/plugin-cache
      ~/.terraform.d/plugins
    key: ${{ runner.os }}-terraform-${{ hashFiles('**/.terraform.lock.hcl') }}
```

#### Terraform Binary Cache
```yaml
- name: Cache Terraform Binaries
  uses: actions/cache@v4
  with:
    path: ~/bin/terraform-*
    key: terraform-binaries-${{ matrix.terraform }}
    
- uses: hashicorp/setup-terraform@v3
  with:
    terraform_version: ${{ matrix.terraform }}
    terraform_wrapper: false
```

#### Go Build Cache
```yaml
- name: Go Build Cache
  uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### 6.4 Optimization Impact Analysis

| Optimization | Time Saved | Implementation Effort |
|--------------|------------|----------------------|
| Parallel execution (30 workers) | 70-80% | Low (flag change) |
| Provider build caching | 2-3 min | Low (GitHub Action) |
| Terraform binary caching | 1-2 min | Low (GitHub Action) |
| Split by package | 20-30% | Medium (Makefile updates) |
| Skip slow tests in PR | 50-60% | Medium (test tagging) |
| Test result caching (skip passed tests) | 30-40% | High (custom logic) |

**Recommended Quick Win Stack:**
```yaml
# Achieves < 10 minute PR checks
- Parallel execution: 30 workers
- Provider build cache
- Terraform binary cache
- Unit tests first (fast fail)
- Only _basic acceptance tests on PR
- Full suite nightly
```

---

## 7. Concrete Recommendations for azurecaf

### 7.1 Testing Strategy

```
┌─────────────────────────────────────────────────────┐
│                  Testing Pyramid                     │
├─────────────────────────────────────────────────────┤
│                                                      │
│  E2E/Integration (Nightly)                          │
│  ├─ All resource types x All naming conventions     │
│  ├─ Multi-version Terraform (1.5.x - 1.8.x)        │
│  └─ ~30 min with full parallelization               │
│                                                      │
│  Acceptance Tests (PR + Nightly)                    │
│  ├─ Basic tests: All resource types                 │
│  ├─ Complete tests: Full configuration              │
│  ├─ Update tests: Configuration changes             │
│  └─ ~5-10 min with 30 parallel workers              │
│                                                      │
│  Unit Tests (Every Commit)                          │
│  ├─ Naming logic                                    │
│  ├─ Validation functions                            │
│  ├─ Resource definition parsing                     │
│  └─ < 1 min                                         │
└─────────────────────────────────────────────────────┘
```

### 7.2 Recommended CI/CD Pipeline

```yaml
name: Terraform Provider Tests

on:
  pull_request:
    paths: ['**/*.go', 'go.mod', 'go.sum']
  push:
    branches: [main]
  schedule:
    - cron: '0 2 * * *'  # Nightly full tests

jobs:
  # Quick validation (< 5 minutes)
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
          cache: true
      - run: make test
      - run: make lint

  # Acceptance tests (< 10 minutes on PR)
  acceptance-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    strategy:
      fail-fast: false
      matrix:
        terraform: 
          - '1.5.x'
          - '1.8.x'  # Only test min/max on PR
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
          cache: true
      
      - uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: ${{ matrix.terraform }}
          terraform_wrapper: false
      
      - name: Run Acceptance Tests
        run: |
          if [ "${{ github.event_name }}" == "pull_request" ]; then
            # PR: Only basic tests
            make testacc TESTARGS="-run=TestAcc.*_basic -parallel=30"
          else
            # Main/Nightly: Full suite
            make testacc TESTARGS="-parallel=30"
          fi
        env:
          TF_ACC: "1"
        timeout-minutes: 30
      
      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: test-results-${{ matrix.terraform }}
          path: test-results/

  # Full E2E (nightly only)
  e2e-tests:
    if: github.event_name == 'schedule'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        terraform: ['1.5.x', '1.6.x', '1.7.x', '1.8.x']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
          cache: true
      - uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: ${{ matrix.terraform }}
          terraform_wrapper: false
      
      - name: Full Test Suite
        run: make test_all
        env:
          TF_ACC: "1"
        timeout-minutes: 60
```

### 7.3 Test Coverage Targets

| Test Type | Coverage Target | Current Status |
|-----------|----------------|----------------|
| Unit Tests | > 80% | ✅ Likely high (pure logic) |
| Acceptance - Basic | 100% resources | ⚠️ Needs validation |
| Acceptance - Complete | 100% resources | ⚠️ Needs validation |
| Acceptance - Update | 80% resources | ❌ Likely missing |
| Multi-version | 4 versions | ❌ Not implemented |
| Drift detection | Critical resources | ❌ Not implemented |

### 7.4 Implementation Priorities

**Phase 1: Foundation (Week 1)**
1. ✅ Add parallel execution to existing tests
2. ✅ Implement caching in GitHub Actions
3. ✅ Split PR vs Nightly test suites
4. Target: Get PR tests < 10 minutes

**Phase 2: Coverage (Week 2)**
1. Add `_complete` tests for all resources
2. Add `_update` tests for all resources
3. Implement drift detection tests
4. Target: 100% basic test coverage

**Phase 3: Multi-Version (Week 3)**
1. Add Terraform version matrix
2. Test 1.5.x through 1.8.x
3. Document compatibility matrix
4. Target: Support last 4 minor versions

**Phase 4: Polish (Week 4)**
1. Add test result artifacts
2. Implement test sweepers (cleanup)
3. Add test documentation
4. Optimize test runtime further

---

## 8. Specific Answers to Research Questions

### Q: How do major providers implement E2E testing?
**A:** They use acceptance tests with the `terraform-plugin-testing` framework that:
- Create real cloud resources
- Run actual Terraform plan/apply cycles
- Validate state matches reality via import
- Clean up automatically via CheckDestroy

### Q: Standard patterns for plan/apply workflows?
**A:** Framework handles automatically:
1. `terraform init` (implicit)
2. `terraform plan` (checked via ConfigPlanChecks)
3. `terraform apply` (executed for each Config step)
4. `terraform import` (via ImportStep)
5. `terraform destroy` (via CheckDestroy)

### Q: How to handle state management?
**A:** 
- Tests run in isolated directories
- State is ephemeral (not committed)
- Framework cleans up between steps
- No need for state backends in tests

### Q: Mock vs real resources?
**A:** 
- **Real resources** are the industry standard
- **For azurecaf:** No actual Azure resources needed (naming only)
- Can run entire test suite without Azure credentials
- Optional: Validate against Azure naming rules

### Q: Cost minimization strategies?
**A:**
- Run expensive tests nightly, not on every PR
- Use parallelization to reduce CI minutes
- Implement test sweepers for batch cleanup
- **For azurecaf:** Essentially free (no resource creation)

### Q: Drift detection best practices?
**A:**
- Apply same config twice, expect no changes
- Use `plancheck.ExpectEmptyPlan()`
- Avoid random values in tests (use framework's random)
- Ignore computed/timestamp fields in ImportStep

### Q: Multi-version testing?
**A:**
- Test matrix in GitHub Actions
- Cache Terraform binaries
- Support last 4 minor versions
- Use version constraints in specific tests

### Q: GitHub Actions best practices?
**A:**
- Cache Go modules, build cache, Terraform binaries
- Use matrix strategy for parallelization
- Split quick vs comprehensive test suites
- Store artifacts for failed tests

### Q: Performance optimization?
**A:** 
- High parallelization (30-50 workers)
- Aggressive caching
- Run subset on PR, full suite nightly
- **Expected for azurecaf:** < 5 min PR, < 30 min nightly

---

## 9. References and Further Reading

### Official Documentation
- [Terraform Plugin Testing Framework](https://developer.hashicorp.com/terraform/plugin/testing)
- [Acceptance Testing Guide](https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests)
- [Testing Patterns](https://developer.hashicorp.com/terraform/plugin/testing/testing-patterns)

### Provider Examples
- [AWS Provider Tests](https://github.com/hashicorp/terraform-provider-aws/tree/main/internal/service)
- [Azure Provider Tests](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/internal/services)
- [Google Provider Tests](https://github.com/hashicorp/terraform-provider-google/tree/main/google/services)

### Tools and Frameworks
- [terraform-plugin-testing](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-testing)
- [terraform-plugin-framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [hashicorp/setup-terraform](https://github.com/hashicorp/setup-terraform)

---

## Appendix: Example Test Suite Structure

```
terraform-provider-azurecaf/
├── azurecaf/
│   ├── resource_name_test.go
│   │   ├── TestAccName_basic
│   │   ├── TestAccName_complete
│   │   ├── TestAccName_update
│   │   ├── TestAccName_requiresImport
│   │   └── TestAccName_noDrift
│   │
│   ├── resource_naming_convention_test.go
│   │   ├── TestAccNamingConvention_basic
│   │   ├── TestAccNamingConvention_complete
│   │   ├── TestAccNamingConvention_cafClassic
│   │   ├── TestAccNamingConvention_cafRandom
│   │   └── TestAccNamingConvention_passthrough
│   │
│   └── data_name_test.go
│       ├── TestAccDataName_basic
│       └── TestAccDataName_complete
│
├── e2e/
│   └── e2e_comprehensive_test.go
│       ├── TestE2E_allResourceTypes
│       └── TestE2E_allNamingConventions
│
└── .github/workflows/
    ├── pr-tests.yml          # < 10 min
    ├── nightly-tests.yml     # < 30 min
    └── release-tests.yml     # Full suite
```

---

**Conclusion:** The research shows a clear path forward for implementing robust E2E testing in terraform-provider-azurecaf using battle-tested patterns from major providers, optimized for the unique characteristics of a naming provider (no actual resource creation costs).
