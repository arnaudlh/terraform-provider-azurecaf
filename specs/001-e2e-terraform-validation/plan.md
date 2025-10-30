# Implementation Plan: End-to-End Terraform Validation Pipeline

**Branch**: `001-e2e-terraform-validation` | **Date**: 2025-10-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-e2e-terraform-validation/spec.md`

## Summary

This feature adds comprehensive end-to-end validation to the CI/CD pipeline, automatically building the Terraform provider, executing `terraform plan`, performing `terraform apply` to create test resources, and running a second apply to detect drift. The pipeline validates that provider-generated names comply with Azure API constraints across all naming conventions and representative resource types, ensuring zero regressions in production deployments.

## Technical Context

**Language/Version**: Go 1.24+ (aligned with project go.mod)  
**Primary Dependencies**: 
- Terraform Plugin SDK v2 (existing)
- Terraform CLI (multiple versions for testing: 1.5.x, 1.6.x, 1.7.x, 1.8.x)
- GitHub Actions for CI/CD orchestration
- Go test framework (existing)

**Storage**: 
- Terraform state files (temporary, in-memory or local backend for testing)
- Pipeline artifacts (logs, plan outputs) stored in GitHub Actions artifacts
- Test configurations stored in repository under `.github/e2e-tests/` or `e2e/`

**Testing**: 
- Go test framework for unit tests (existing)
- Terraform acceptance testing framework (TF_ACC=1) for integration tests (existing)
- End-to-end validation scripts for pipeline testing (new)
- Mock Azure provider or dedicated test subscription for apply validation

**Target Platform**: 
- GitHub Actions runners (Linux, potentially Windows/macOS for matrix testing)
- CI/CD environment with Terraform CLI installed

**Project Type**: Single project (Terraform provider with CI/CD enhancements)

**Performance Goals**: 
- Quick validation (PR context): Complete in under 10 minutes
- Comprehensive validation (main branch): Complete in under 30 minutes
- Build phase: Complete in under 2 minutes
- Plan phase: Complete in under 3 minutes
- Apply + drift detection: Complete in under 10 minutes

**Constraints**: 
- Must not incur significant Azure costs (use mocks or minimal test resources)
- Must clean up test resources reliably (95%+ success rate)
- Must provide fast feedback (failures detected within 5 minutes)
- Must work within GitHub Actions free tier limits (for open source)
- Must handle Azure API rate limits and transient failures gracefully

**Scale/Scope**: 
- 395 supported resource types (comprehensive validation on main branch)
- 20+ representative resource types (quick validation on PRs)
- 4 naming conventions to validate (cafclassic, cafrandom, random, passthrough)
- Multiple Terraform versions (1.5.x through 1.8.x)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ I. Azure Naming Convention Compliance (NON-NEGOTIABLE)
**Status**: PASS  
**Rationale**: This feature enhances validation of naming convention compliance by testing provider-generated names against actual Azure API constraints. It validates all 4 naming conventions and ensures names work in real-world scenarios, strengthening compliance rather than violating it.

### ✅ II. Test Coverage Discipline (NON-NEGOTIABLE)
**Status**: PASS  
**Rationale**: This feature adds comprehensive E2E testing on top of existing unit and integration tests. The pipeline itself will include tests for the E2E framework. Expected test coverage for new E2E code: >95%.

### ✅ III. Comprehensive Documentation
**Status**: PASS  
**Rationale**: Will include:
- CHANGELOG.md update (MINOR version - new feature without breaking changes)
- README.md update documenting E2E validation in CI/CD section
- E2E testing documentation in new docs/E2E_VALIDATION.md
- GitHub Actions workflow documentation inline

### ✅ IV. Resource Definition Accuracy
**Status**: PASS  
**Rationale**: This feature validates resource definition accuracy by testing provider-generated names against Azure APIs. Any inaccuracies in resourceDefinition.json will be caught during apply phase, strengthening this principle.

### ✅ V. Backward Compatibility
**Status**: PASS  
**Rationale**: This is a pure CI/CD enhancement that doesn't change provider behavior, APIs, or generated names. Completely backward compatible. Version bump: MINOR (new capability, no breaking changes).

### ✅ VI. Security and Code Quality
**Status**: PASS  
**Rationale**: 
- All new code will pass Go vet and linting
- No new dependencies introduced (uses existing Terraform CLI and GitHub Actions)
- GitHub Actions secrets used for Azure credentials (if real resources used)
- Test cleanup prevents resource leakage
- No file permission changes required

## Project Structure

### Documentation (this feature)

```text
specs/001-e2e-terraform-validation/
├── plan.md              # This file (implementation plan)
├── research.md          # Phase 0: Research findings on E2E testing approaches
├── data-model.md        # Phase 1: Test configuration and validation data models
├── quickstart.md        # Phase 1: Quick start guide for running E2E validation
├── contracts/           # Phase 1: Pipeline stage contracts and interfaces
│   └── validation-pipeline.yaml    # Pipeline stage definitions and contracts
└── tasks.md             # Phase 2: Implementation task list (created by /speckit.tasks)
```

### Source Code (repository root)

```text
.github/
├── workflows/
│   ├── e2e-validation-pr.yml        # Quick E2E validation for pull requests
│   ├── e2e-validation-main.yml      # Comprehensive E2E validation for main branch
│   └── e2e-validation-manual.yml    # Manual trigger with configurable options
└── e2e-tests/                       # E2E test configurations and scripts
    ├── configs/                     # Terraform test configurations
    │   ├── quick/                   # Quick validation configs (20 resources)
    │   │   ├── cafclassic.tf
    │   │   ├── cafrandom.tf
    │   │   ├── random.tf
    │   │   └── passthrough.tf
    │   └── comprehensive/           # Full validation configs (395 resources)
    │       ├── compute.tf
    │       ├── storage.tf
    │       ├── networking.tf
    │       ├── databases.tf
    │       └── ...
    ├── scripts/                     # E2E validation scripts
    │   ├── build-provider.sh        # Build provider binary
    │   ├── setup-terraform.sh       # Setup Terraform test environment
    │   ├── run-validation.sh        # Main validation orchestrator
    │   ├── cleanup-resources.sh     # Resource cleanup script
    │   └── check-drift.sh           # Drift detection script
    └── fixtures/                    # Test fixtures and mocks
        ├── mock-azure-provider/     # Mock Azure provider (if using mocks)
        └── test-state-backend/      # Test state backend configuration

e2e/                                 # E2E testing framework (already exists)
├── README.md                        # E2E framework documentation (update)
├── e2e_comprehensive_test.go        # Comprehensive E2E tests (enhance)
├── e2e_test.go                      # Core E2E tests (enhance)
└── framework/                       # E2E framework utilities
    ├── framework.go                 # Framework implementation
    └── utils.go                     # Utility functions

docs/
└── E2E_VALIDATION.md                # E2E validation documentation (new)

scripts/                             # Existing scripts directory
└── validate-e2e-ci.sh               # CI validation script (enhance existing)
```

**Structure Decision**: 

This feature uses the existing single-project structure with CI/CD enhancements. Key decisions:

1. **GitHub Actions Workflows**: Three separate workflows for different validation contexts (PR quick, main comprehensive, manual testing)

2. **E2E Test Configurations**: Organized under `.github/e2e-tests/` to keep CI-specific files together. Separate directories for quick (subset) and comprehensive (all resources) validation.

3. **Existing E2E Framework**: Enhance the existing `e2e/` directory rather than creating new structure. This framework already supports provider testing.

4. **Script Organization**: Validation scripts under `.github/e2e-tests/scripts/` to co-locate with test configurations.

5. **Documentation**: New E2E_VALIDATION.md in docs/ for comprehensive documentation, update existing e2e/README.md for framework enhancements.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

**No violations**: All constitution checks passed. No complexity justification needed.

## Phase 0: Research

**Status**: ✅ Complete  
**Output**: [research.md](./research.md)

### Key Findings

1. **Industry Standard Approach**: Use `terraform-plugin-testing` framework for E2E validation with real provider binary testing in isolated environments

2. **Unique Provider Advantage**: azurecaf provider only generates names without creating Azure resources, enabling:
   - No Azure subscription required for testing
   - Zero cost E2E validation
   - Fast execution (no API calls to Azure)
   - No quota/rate limit concerns

3. **Performance Strategy**: 
   - Two-tier validation: Quick (PR, 10 min) + Comprehensive (main, 30 min)
   - Aggressive parallelization (20-50 concurrent jobs)
   - Smart caching (Go modules, Terraform plugins)
   - Subset testing for PRs (20 representative resources)

4. **Architecture Decisions**: Four ADRs documented covering testing approach, parallelization, validation tiers, and test configuration management

See [research.md](./research.md) for complete findings, competitor analysis, and architectural decisions.

## Phase 1: Design & Contracts

**Status**: ✅ Complete

### Generated Artifacts

1. **Data Models** ([data-model.md](./data-model.md))
   - ValidationConfiguration: Test suite definition
   - TestResource: Individual resource type testing
   - NamingConventionTest: Convention-specific testing
   - ValidationResult: Overall validation outcomes
   - StageResult: Individual stage outcomes
   - ResourceTestResult: Per-resource test results
   - TerraformState: State management for drift detection

2. **Pipeline Contracts** ([contracts/validation-pipeline.md](./contracts/validation-pipeline.md))
   - Stage 1: Build (provider compilation)
   - Stage 2: Plan (terraform init/plan)
   - Stage 3: Apply (terraform apply)
   - Stage 4: Drift Check (second apply for drift detection)
   - Stage 5: Cleanup (resource and file cleanup)
   - Orchestration rules and dependencies

3. **Quick Start Guide** ([quickstart.md](./quickstart.md))
   - Local development setup
   - CI/CD integration
   - Common use cases
   - Troubleshooting guide

### Agent Context

**Updated**: ✅ GitHub Copilot context updated with technology stack (Go 1.24+)

### Constitution Re-Check (Post-Design)

*Re-evaluation after Phase 1 design completion:*

#### ✅ I. Azure Naming Convention Compliance
**Status**: PASS  
**Design Validation**: Data models and contracts enforce validation of all 4 naming conventions. Pipeline tests provider-generated names against Azure constraints.

#### ✅ II. Test Coverage Discipline
**Status**: PASS  
**Design Validation**: E2E framework includes comprehensive test coverage tracking. Pipeline validates 395 resource types in comprehensive mode.

#### ✅ III. Comprehensive Documentation
**Status**: PASS  
**Design Validation**: Three documentation artifacts created (data-model.md, contracts/validation-pipeline.md, quickstart.md). Full documentation plan in project structure.

#### ✅ IV. Resource Definition Accuracy
**Status**: PASS  
**Design Validation**: Pipeline contracts include validation of resource definition accuracy through apply stage. Any inaccuracies will be caught during testing.

#### ✅ V. Backward Compatibility
**Status**: PASS  
**Design Validation**: Design is purely additive (new workflows, test configs, scripts). No changes to existing provider code or APIs.

#### ✅ VI. Security and Code Quality
**Status**: PASS  
**Design Validation**: Contracts specify error handling, cleanup procedures, and timeout strategies. GitHub Actions secrets for credentials. No elevated permissions required.

**Overall Status**: ✅ ALL GATES PASSED - Design phase complete and constitution-compliant

## Phase 2: Tasks

**Status**: ⏳ Pending  
**Next Step**: Run `/speckit.tasks` command to generate implementation task list

The task list will be created in [tasks.md](./tasks.md) with:
- Detailed implementation tasks for each pipeline stage
- Dependencies and ordering
- Time estimates
- Acceptance criteria

## Summary

**Branch**: `001-e2e-terraform-validation`  
**Implementation Plan**: `/Users/arnaud/Documents/github/arnaudlh/terraform-provider-azurecaf/specs/001-e2e-terraform-validation/plan.md`

### Generated Artifacts

- ✅ [spec.md](./spec.md) - Feature specification
- ✅ [research.md](./research.md) - E2E testing research and ADRs
- ✅ [data-model.md](./data-model.md) - Data models and schemas
- ✅ [contracts/validation-pipeline.md](./contracts/validation-pipeline.md) - Pipeline stage contracts
- ✅ [quickstart.md](./quickstart.md) - Quick start guide
- ✅ Agent context updated (GitHub Copilot)
- ⏳ [tasks.md](./tasks.md) - Implementation tasks (generated by `/speckit.tasks`)

### Constitution Compliance

All 6 constitutional principles validated and passed in both pre-research and post-design checks.

### Next Steps

1. Run `/speckit.tasks` to generate implementation task list
2. Begin implementation following task list
3. Update CHANGELOG.md as changes are made
4. Create pull request when ready for review

---

**Planning Phase Complete** ✅

````
