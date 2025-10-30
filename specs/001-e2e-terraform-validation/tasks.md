# Tasks: End-to-End Terraform Validation Pipeline

**Feature Branch**: `001-e2e-terraform-validation`  
**Date**: 2025-10-30  
**Input**: Design documents from `/specs/001-e2e-terraform-validation/`

## Task Format

Each task follows the format: `- [ ] [TaskID] [P?] [Story?] Description with file path`

- **[P]**: Parallelizable (different files, no blocking dependencies)
- **[Story]**: User story label (US1-US6) from spec.md
- All file paths are relative to repository root

## Overview

This task list implements E2E Terraform validation across 6 user stories:
- **US1 (P1)**: Automated Provider Build Verification
- **US2 (P1)**: Terraform Plan Validation  
- **US3 (P1)**: Terraform Apply and Resource Creation
- **US4 (P1)**: Drift Detection with Second Apply
- **US5 (P2)**: Multi-Environment Testing
- **US6 (P3)**: CI/CD Pipeline Performance

**Total Tasks**: 65  
**Estimated Duration**: 4-5 weeks  
**MVP Scope**: US1 + US2 (build and plan validation)

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Create directory structure and foundational configuration

- [x] T001 Create .github/e2e-tests/ directory structure per plan.md
- [x] T002 Create .github/e2e-tests/configs/ subdirectories (quick/, comprehensive/)
- [x] T003 Create .github/e2e-tests/scripts/ directory for validation scripts
- [x] T004 Create .github/e2e-tests/fixtures/ directory for test fixtures
- [x] T005 [P] Create docs/E2E_VALIDATION.md initial structure

**Duration**: 1-2 hours  
**Checkpoint**: ✅ Directory structure matches plan.md specification

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Create test configuration schema in .github/e2e-tests/configs/schema.json
- [x] T007 Implement validation config loader in .github/e2e-tests/scripts/lib/config-loader.sh
- [x] T008 [P] Create test resource selector in .github/e2e-tests/scripts/lib/resource-selector.sh
- [x] T009 [P] Implement logging framework in .github/e2e-tests/scripts/lib/logger.sh
- [x] T010 [P] Create error handling utilities in .github/e2e-tests/scripts/lib/error-handler.sh
- [x] T011 Create test workspace manager in .github/e2e-tests/scripts/lib/workspace-manager.sh
- [x] T012 Implement stage result tracker in .github/e2e-tests/scripts/lib/result-tracker.sh

**Duration**: 2-3 days  
**Checkpoint**: ✅ Foundation scripts can load configs, create workspaces, track results

---

## Phase 3: User Story 1 - Automated Provider Build Verification (Priority: P1) 🎯 MVP

**Goal**: Automatically build provider binary on every code change and validate compilation success

**Independent Test**: Push code change → CI builds provider → binary artifact available

### Implementation for User Story 1

- [x] T013 [P] [US1] Create build-provider.sh script in .github/e2e-tests/scripts/build-provider.sh
- [x] T014 [P] [US1] Implement Go module caching logic in build-provider.sh
- [x] T015 [P] [US1] Add build duration tracking in build-provider.sh
- [x] T016 [P] [US1] Add binary hash calculation in build-provider.sh
- [x] T017 [US1] Create e2e-validation-pr.yml workflow in .github/workflows/e2e-validation-pr.yml (build stage only)
- [x] T018 [US1] Configure Go setup with caching in e2e-validation-pr.yml
- [x] T019 [US1] Add provider binary artifact upload in e2e-validation-pr.yml
- [x] T020 [US1] Add build failure reporting in e2e-validation-pr.yml

**Duration**: 2-3 days  
**Checkpoint**: CI builds provider binary successfully, artifacts available, build time < 2 min

---

## Phase 4: User Story 2 - Terraform Plan Validation (Priority: P1) 🎯 MVP

**Goal**: Execute terraform plan with provider and validate successful plan generation for all naming conventions

**Independent Test**: Run terraform plan with test config → plan succeeds → names generated correctly

### Implementation for User Story 2

- [x] T021 [P] [US2] Create quick validation configs in .github/e2e-tests/configs/quick/cafclassic.tf
- [x] T022 [P] [US2] Create quick validation configs in .github/e2e-tests/configs/quick/cafrandom.tf
- [x] T023 [P] [US2] Create quick validation configs in .github/e2e-tests/configs/quick/random.tf
- [x] T024 [P] [US2] Create quick validation configs in .github/e2e-tests/configs/quick/passthrough.tf
- [x] T025 [P] [US2] Create provider configuration template in .github/e2e-tests/configs/quick/provider.tf
- [x] T026 [US2] Implement setup-terraform.sh script in .github/e2e-tests/scripts/setup-terraform.sh
- [x] T027 [US2] Add Terraform version management in setup-terraform.sh
- [x] T028 [US2] Add provider installation logic in setup-terraform.sh
- [x] T029 [US2] Create run-validation.sh orchestrator in .github/e2e-tests/scripts/run-validation.sh
- [x] T030 [US2] Implement plan stage in run-validation.sh (terraform init + plan)
- [x] T031 [US2] Add plan output parsing in run-validation.sh
- [x] T032 [US2] Add name validation checks in run-validation.sh
- [x] T033 [US2] Update e2e-validation-pr.yml to add plan stage
- [x] T034 [US2] Add plan artifact upload in e2e-validation-pr.yml

**Duration**: 3-4 days  
**Checkpoint**: Terraform plan succeeds for 20 representative resources, all 4 conventions, plan time < 3 min

---

## Phase 5: User Story 3 - Terraform Apply and Resource Creation (Priority: P1)

**Goal**: Execute terraform apply to validate provider-generated names work with Azure API validation rules

**Independent Test**: Run terraform apply → resources created → Azure accepts all names

### Implementation for User Story 3

- [ ] T035 [US3] Implement apply stage in run-validation.sh (terraform apply)
- [ ] T036 [US3] Add apply success validation in run-validation.sh
- [ ] T037 [US3] Add state file verification in run-validation.sh
- [ ] T038 [US3] Add resource creation tracking in run-validation.sh
- [ ] T039 [US3] Implement apply timeout handling in run-validation.sh
- [ ] T040 [US3] Update e2e-validation-pr.yml to add apply stage
- [ ] T041 [US3] Add state file artifact upload in e2e-validation-pr.yml

**Duration**: 2-3 days  
**Checkpoint**: Terraform apply succeeds for all test resources, state files valid, apply time < 10 min

---

## Phase 6: User Story 4 - Drift Detection with Second Apply (Priority: P1)

**Goal**: Run second terraform apply to verify no drift detected, confirming state management correctness

**Independent Test**: Run apply twice → second apply reports "No changes" → no drift detected

### Implementation for User Story 4

- [ ] T042 [US4] Create check-drift.sh script in .github/e2e-tests/scripts/check-drift.sh
- [ ] T043 [US4] Implement second terraform plan in check-drift.sh
- [ ] T044 [US4] Add drift detection logic in check-drift.sh
- [ ] T045 [US4] Add drift details reporting in check-drift.sh
- [ ] T046 [US4] Integrate drift check into run-validation.sh
- [ ] T047 [US4] Update e2e-validation-pr.yml to add drift check stage
- [ ] T048 [US4] Add drift report artifact upload in e2e-validation-pr.yml
- [ ] T049 [US4] Add random seed stability tests in check-drift.sh

**Duration**: 2-3 days  
**Checkpoint**: Second apply shows zero drift for 100% of resources, random seeds stable

---

## Phase 7: User Story 5 - Multi-Environment Testing (Priority: P2)

**Goal**: Test provider against multiple Terraform versions and all naming conventions for broad compatibility

**Independent Test**: Run validation with TF 1.5.7, 1.6.6, 1.7.5, 1.8.0 → all versions succeed

### Implementation for User Story 5

- [ ] T050 [P] [US5] Create comprehensive validation configs in .github/e2e-tests/configs/comprehensive/compute.tf
- [ ] T051 [P] [US5] Create comprehensive validation configs in .github/e2e-tests/configs/comprehensive/storage.tf
- [ ] T052 [P] [US5] Create comprehensive validation configs in .github/e2e-tests/configs/comprehensive/networking.tf
- [ ] T053 [P] [US5] Create comprehensive validation configs in .github/e2e-tests/configs/comprehensive/databases.tf
- [ ] T054 [P] [US5] Create comprehensive validation configs in .github/e2e-tests/configs/comprehensive/analytics.tf
- [ ] T055 [US5] Create e2e-validation-main.yml workflow in .github/workflows/e2e-validation-main.yml
- [ ] T056 [US5] Add Terraform version matrix strategy in e2e-validation-main.yml
- [ ] T057 [US5] Add naming convention matrix strategy in e2e-validation-main.yml
- [ ] T058 [US5] Configure parallelization (20-50 jobs) in e2e-validation-main.yml
- [ ] T059 [US5] Add comprehensive test reporting in e2e-validation-main.yml

**Duration**: 3-4 days  
**Checkpoint**: All 395 resource types tested, all TF versions work, all conventions pass

---

## Phase 8: User Story 6 - CI/CD Pipeline Performance (Priority: P3)

**Goal**: Optimize pipeline to complete quick validation < 10 min, comprehensive < 30 min

**Independent Test**: Measure pipeline duration → quick < 10 min, comprehensive < 30 min

### Implementation for User Story 6

- [ ] T060 [P] [US6] Implement parallel execution groups in run-validation.sh
- [ ] T061 [P] [US6] Add performance metrics tracking in run-validation.sh
- [ ] T062 [P] [US6] Optimize Terraform binary caching in setup-terraform.sh
- [ ] T063 [US6] Create e2e-validation-manual.yml workflow in .github/workflows/e2e-validation-manual.yml
- [ ] T064 [US6] Add manual workflow inputs (mode, resources, conventions) in e2e-validation-manual.yml
- [ ] T065 [US6] Add performance reporting dashboard to e2e-validation-manual.yml

**Duration**: 2-3 days  
**Checkpoint**: PR validation completes in 7-8 min, main validation in 25-30 min

---

## Phase 9: Cleanup & Reliability

**Purpose**: Resource cleanup and error handling across all user stories

- [ ] T066 Create cleanup-resources.sh script in .github/e2e-tests/scripts/cleanup-resources.sh
- [ ] T067 Implement terraform destroy logic in cleanup-resources.sh
- [ ] T068 Add temporary file cleanup in cleanup-resources.sh
- [ ] T069 Add force cleanup option in cleanup-resources.sh
- [ ] T070 Integrate cleanup into run-validation.sh (always-run stage)
- [ ] T071 Add cleanup failure handling in e2e-validation-pr.yml
- [ ] T072 Add cleanup artifact upload in e2e-validation-pr.yml

**Duration**: 1-2 days  
**Checkpoint**: Cleanup succeeds 95%+ of time, no resource leakage

---

## Phase 10: Documentation & Polish

**Purpose**: Complete documentation and cross-cutting improvements

- [ ] T073 [P] Complete docs/E2E_VALIDATION.md with full usage guide
- [ ] T074 [P] Update README.md CI/CD section with E2E validation info
- [ ] T075 [P] Update CHANGELOG.md with MINOR version bump and feature description
- [ ] T076 [P] Add inline documentation to all scripts (headers, comments)
- [ ] T077 [P] Create examples for manual validation runs in docs/E2E_VALIDATION.md
- [ ] T078 Run complete quickstart.md validation scenarios
- [ ] T079 Validate all GitHub Actions workflows syntax
- [ ] T080 Verify constitution compliance (test coverage >95%, documentation complete)

**Duration**: 2-3 days  
**Checkpoint**: Documentation complete, examples work, ready for review

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational) ← BLOCKING - Must complete before any user story
    ↓
    ├─→ Phase 3 (US1 - Build) ← MVP Start
    ↓
    ├─→ Phase 4 (US2 - Plan) ← MVP Minimum
    ↓
    ├─→ Phase 5 (US3 - Apply)
    ↓
    ├─→ Phase 6 (US4 - Drift)
    ↓
    ├─→ Phase 7 (US5 - Multi-Env)
    ↓
    ├─→ Phase 8 (US6 - Performance)
    ↓
Phase 9 (Cleanup) ← Integrates with all user stories
    ↓
Phase 10 (Documentation)
```

### User Story Dependencies

**Critical Path** (P1 stories):
```
US1 (Build) → US2 (Plan) → US3 (Apply) → US4 (Drift)
```

**Parallel Opportunities** (after foundational):
- US1 can start immediately after Phase 2
- US2 requires US1 complete (needs binary)
- US3 requires US2 complete (needs plan)
- US4 requires US3 complete (needs state)
- US5 can start in parallel with US3/US4 (independent configs)
- US6 optimizes across all stories (start after US1-US4 functional)

### Within Each User Story

**US1 (Build)**:
- T013-T016 (script tasks) can run in parallel
- T017-T020 (workflow tasks) depend on scripts

**US2 (Plan)**:
- T021-T025 (config files) can run in parallel
- T026-T028 (setup script) can run in parallel with configs
- T029-T032 (orchestrator) depend on setup script
- T033-T034 (workflow) depend on orchestrator

**US3 (Apply)**:
- All tasks depend on US2 orchestrator (sequential)

**US4 (Drift)**:
- T042-T045 (drift script) can develop in parallel
- T046-T049 (integration) depend on drift script

**US5 (Multi-Env)**:
- T050-T054 (configs) can run in parallel
- T055-T059 (workflow) depend on configs

**US6 (Performance)**:
- T060-T062 (optimization) can run in parallel
- T063-T065 (manual workflow) can run in parallel

---

## Parallel Execution Examples

### Phase 1 Setup (Parallel)
```bash
Task T001 & Task T002 & Task T003 & Task T004 & Task T005
# All directory creation can happen simultaneously
```

### Phase 2 Foundational (Partial Parallel)
```bash
Task T006  # Config schema first
wait
Task T007 & Task T008 & Task T009 & Task T010  # Scripts in parallel
wait
Task T011 & Task T012  # Workspace and tracking
```

### Phase 3 US1 Build (Parallel Scripts)
```bash
Task T013 & Task T014 & Task T015 & Task T016  # Build script components
wait
Task T017  # Workflow (depends on scripts)
Task T018 & Task T019 & Task T020  # Workflow enhancements in parallel
```

### Phase 4 US2 Plan (Parallel Configs)
```bash
Task T021 & Task T022 & Task T023 & Task T024 & Task T025  # All configs
Task T026 & Task T027 & Task T028  # Setup script in parallel
wait
Task T029 & Task T030 & Task T031 & Task T032  # Orchestrator (sequential)
wait
Task T033 & Task T034  # Workflow updates
```

### Phase 7 US5 Multi-Env (Parallel Configs)
```bash
Task T050 & Task T051 & Task T052 & Task T053 & Task T054  # All comprehensive configs
wait
Task T055 & Task T056 & Task T057 & Task T058 & Task T059  # Workflow
```

### Phase 10 Documentation (All Parallel)
```bash
Task T073 & Task T074 & Task T075 & Task T076 & Task T077  # All docs in parallel
```

---

## Implementation Strategy

### Strategy 1: MVP First (Recommended for Solo Developer)

**Goal**: Get basic E2E validation working quickly

**Phases**:
1. Complete Phase 1 (Setup) - 2 hours
2. Complete Phase 2 (Foundational) - 2-3 days
3. Complete Phase 3 (US1 - Build) - 2-3 days
4. Complete Phase 4 (US2 - Plan) - 3-4 days
5. **STOP and VALIDATE**: Test US1+US2 independently
6. Demo: "Provider builds and terraform plan validates names"

**Timeline**: 1.5 weeks for MVP  
**Deliverable**: PR validation with build + plan (no apply yet)

### Strategy 2: Critical Path (Full P1 Features)

**Goal**: Complete all P1 user stories for production-ready validation

**Phases**:
1. Complete Phases 1-2 (Setup + Foundation) - 3-4 days
2. Complete Phase 3 (US1) - 2-3 days
3. Complete Phase 4 (US2) - 3-4 days
4. Complete Phase 5 (US3) - 2-3 days
5. Complete Phase 6 (US4) - 2-3 days
6. Complete Phase 9 (Cleanup) - 1-2 days
7. Complete Phase 10 (Documentation) - 2-3 days
8. **STOP and VALIDATE**: Test full pipeline end-to-end

**Timeline**: 3-4 weeks for full P1  
**Deliverable**: Complete E2E validation with drift detection

### Strategy 3: Parallel Team (Multi-Developer)

**Goal**: Maximize velocity with parallel development

**Team Assignment**:
- **Developer A**: Phases 1-2 (Foundation) → Phase 3 (US1) → Phase 4 (US2)
- **Developer B**: Phase 7 (US5 - configs) → Phase 8 (US6 - performance)
- **Developer C**: Phase 5 (US3) → Phase 6 (US4) → Phase 9 (Cleanup)
- **Developer D**: Phase 10 (Documentation) throughout

**Coordination Points**:
- Week 1: All on Foundation (Phase 2)
- Week 2: Dev A on US1/US2, Dev B on configs, Dev C waiting
- Week 3: Dev A finishing US2, Dev B on US5, Dev C starts US3/US4
- Week 4: Integration, cleanup, documentation

**Timeline**: 2-3 weeks with team

### Strategy 4: Incremental Delivery

**Goal**: Ship working features incrementally

**Releases**:
- **v1.0 (Week 1.5)**: MVP - Build + Plan validation only
- **v1.1 (Week 2.5)**: Add Apply validation
- **v1.2 (Week 3.5)**: Add Drift detection  
- **v2.0 (Week 4.5)**: Add Multi-environment + Performance

Each release adds value without breaking previous functionality.

---

## Task Summary

**Total Tasks**: 80  
**Tasks by User Story**:
- Setup: 5 tasks
- Foundational: 7 tasks (BLOCKING)
- US1 (Build): 8 tasks
- US2 (Plan): 14 tasks
- US3 (Apply): 7 tasks
- US4 (Drift): 8 tasks
- US5 (Multi-Env): 10 tasks
- US6 (Performance): 6 tasks
- Cleanup: 7 tasks
- Documentation: 8 tasks

**Parallel Opportunities**: 35 tasks marked [P] (44% parallelizable)

**Independent Test Criteria**:
- US1: Provider builds successfully, binary artifact available
- US2: Terraform plan succeeds for 20 resources, all conventions
- US3: Terraform apply creates resources, names accepted by Azure
- US4: Second apply reports zero drift
- US5: All 395 resources tested across 4 TF versions
- US6: Pipeline completes within performance targets

**Suggested MVP Scope**: 
- Phases 1-4 (Setup + Foundation + US1 + US2)
- **Tasks**: T001-T034 (34 tasks)
- **Duration**: 1.5 weeks
- **Deliverable**: PR validation with build + plan

**Full Feature Scope**:
- All phases (1-10)
- **Tasks**: T001-T080 (80 tasks)
- **Duration**: 4-5 weeks solo, 2-3 weeks with team
- **Deliverable**: Complete E2E validation pipeline

---

## Format Validation ✅

**Checklist Format Compliance**:
- ✅ All tasks start with `- [ ]` (markdown checkbox)
- ✅ All tasks have sequential IDs (T001-T080)
- ✅ All tasks have [P] marker where parallelizable
- ✅ All user story tasks have [Story] label (US1-US6)
- ✅ All tasks include file paths in descriptions
- ✅ Setup/Foundational/Polish phases have NO story label
- ✅ User story phases have REQUIRED story labels

**Example Task Formats**:
- ✅ `- [ ] T001 Create .github/e2e-tests/ directory structure per plan.md`
- ✅ `- [ ] T013 [P] [US1] Create build-provider.sh script in .github/e2e-tests/scripts/build-provider.sh`
- ✅ `- [ ] T029 [US2] Create run-validation.sh orchestrator in .github/e2e-tests/scripts/run-validation.sh`

---

## Next Steps

1. Review tasks against spec.md user stories (verify all stories covered)
2. Validate task breakdown against plan.md technical context
3. Confirm file paths match plan.md project structure
4. Choose implementation strategy (MVP first recommended)
5. Begin Phase 1 (Setup) immediately
6. Complete Phase 2 (Foundational) before any user story work
7. Track progress using task checkboxes
8. Update CHANGELOG.md as you complete each phase

**Ready to implement!** 🚀
