#!/usr/bin/env bash
# E2E Validation Orchestrator for terraform-provider-azurecaf
# Coordinates build, plan, apply, drift detection, and cleanup stages

set -euo pipefail

# Source required libraries
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LIB_DIR="${SCRIPT_DIR}/lib"

# shellcheck source=lib/logger.sh
source "${LIB_DIR}/logger.sh"
# shellcheck source=lib/error-handler.sh
source "${LIB_DIR}/error-handler.sh"
# shellcheck source=lib/config-loader.sh
source "${LIB_DIR}/config-loader.sh"
# shellcheck source=lib/resource-selector.sh
source "${LIB_DIR}/resource-selector.sh"
# shellcheck source=lib/workspace-manager.sh
source "${LIB_DIR}/workspace-manager.sh"
# shellcheck source=lib/result-tracker.sh
source "${LIB_DIR}/result-tracker.sh"

# Script information
readonly SCRIPT_NAME="run-validation"
readonly SCRIPT_VERSION="1.0.0"

# Paths
readonly PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"
readonly CONFIG_DIR="${SCRIPT_DIR}/../configs"
readonly RESOURCE_DEFINITION="${PROJECT_ROOT}/resourceDefinition.json"

# Default values
DEFAULT_MODE="quick"
DEFAULT_TERRAFORM_VERSION="1.5.7"
DEFAULT_TIMEOUT_MINUTES="15"
DEFAULT_PARALLEL_JOBS="1"

# Configuration
MODE="${MODE:-${DEFAULT_MODE}}"
TERRAFORM_VERSION="${TERRAFORM_VERSION:-${DEFAULT_TERRAFORM_VERSION}}"
TIMEOUT_MINUTES="${TIMEOUT_MINUTES:-${DEFAULT_TIMEOUT_MINUTES}}"
PARALLEL_JOBS="${PARALLEL_JOBS:-${DEFAULT_PARALLEL_JOBS}}"
SKIP_CLEANUP="${SKIP_CLEANUP:-false}"
SKIP_BUILD="${SKIP_BUILD:-false}"
SKIP_APPLY="${SKIP_APPLY:-true}"
PROVIDER_BINARY="${PROVIDER_BINARY:-}"
WORKSPACE_NAME=""
WORKSPACE_PATH=""
TERRAFORM_BIN=""
CONVENTIONS=()
VALIDATION_ID=""

#######################################
# Print usage information
#######################################
usage() {
  cat <<EOF
Usage: ${SCRIPT_NAME}.sh [OPTIONS]

Run E2E validation for terraform-provider-azurecaf.

OPTIONS:
  --mode MODE                  Validation mode: quick, standard, comprehensive (default: quick)
  --conventions CONV           Naming conventions to test (comma-separated): cafclassic,cafrandom,random,passthrough
  --terraform-version VERSION  Terraform version to use (default: 1.5.7)
  --parallel JOBS              Number of parallel test jobs (default: 1)
  --timeout MINUTES            Timeout in minutes (default: 15)
  --skip-cleanup               Skip workspace cleanup after tests
  --skip-build                 Skip provider build (use existing binary)
  --skip-apply                 Skip terraform apply (plan only, default: true)
  --enable-apply               Enable terraform apply (creates resources)
  --provider-binary PATH       Path to pre-built provider binary
  --workspace-name NAME        Custom workspace name (default: auto-generated)
  --debug                      Enable debug logging
  --help                       Show this help message

ENVIRONMENT VARIABLES:
  MODE                  Override default validation mode
  TERRAFORM_VERSION     Override default Terraform version
  TIMEOUT_MINUTES       Override default timeout
  PARALLEL_JOBS         Override default parallel jobs
  SKIP_CLEANUP          Skip cleanup (true/false)
  SKIP_BUILD            Skip build (true/false)
  PROVIDER_BINARY       Path to provider binary
  LOG_LEVEL             Set log level (DEBUG=0, INFO=1, WARN=2, ERROR=3)

VALIDATION MODES:
  quick          - Fast validation with 20 resources (~5-8 min)
  standard       - Standard validation with 100 resources (~10-15 min)
  comprehensive  - Full validation with all resources (~25-30 min)

EXAMPLES:
  # Quick validation (default)
  ${SCRIPT_NAME}.sh

  # Test specific conventions
  ${SCRIPT_NAME}.sh --mode quick --conventions cafclassic,cafrandom

  # Use pre-built provider
  ${SCRIPT_NAME}.sh --skip-build --provider-binary ./terraform-provider-azurecaf

  # Comprehensive validation with debug logging
  ${SCRIPT_NAME}.sh --mode comprehensive --debug

EXIT CODES:
  0   - All tests passed
  1   - Build failure
  2   - Dependency error
  10  - Terraform init failure
  11  - Terraform plan failure
  12  - Plan validation failure
  20  - Terraform apply failure (future)
  30  - Drift detected (future)
  40  - Cleanup failed (warning)
  124 - Timeout
  125 - Configuration error

EOF
}

#######################################
# Parse command-line arguments
#######################################
parse_args() {
  while [[ $# -gt 0 ]]; do
    case $1 in
      --mode)
        MODE="$2"
        shift 2
        ;;
      --conventions)
        IFS=',' read -ra CONVENTIONS <<< "$2"
        shift 2
        ;;
      --terraform-version)
        TERRAFORM_VERSION="$2"
        shift 2
        ;;
      --parallel)
        PARALLEL_JOBS="$2"
        shift 2
        ;;
      --timeout)
        TIMEOUT_MINUTES="$2"
        shift 2
        ;;
      --skip-cleanup)
        SKIP_CLEANUP=true
        shift
        ;;
      --skip-build)
        SKIP_BUILD=true
        shift
        ;;
      --skip-apply)
        SKIP_APPLY=true
        shift
        ;;
      --enable-apply)
        SKIP_APPLY=false
        shift
        ;;
      --provider-binary)
        PROVIDER_BINARY="$2"
        SKIP_BUILD=true
        shift 2
        ;;
      --workspace-name)
        WORKSPACE_NAME="$2"
        shift 2
        ;;
      --debug)
        export LOG_LEVEL=0
        shift
        ;;
      --help)
        usage
        exit 0
        ;;
      *)
        log_error "Unknown option: $1"
        usage
        exit "${EXIT_CONFIG_ERROR}"
        ;;
    esac
  done
  
  # Validate mode
  if [[ ! "${MODE}" =~ ^(quick|standard|comprehensive)$ ]]; then
    log_error "Invalid mode: ${MODE}"
    exit "${EXIT_CONFIG_ERROR}"
  fi
  
  # Set default conventions if not specified
  if [[ ${#CONVENTIONS[@]} -eq 0 ]]; then
    CONVENTIONS=("cafclassic" "cafrandom" "random" "passthrough")
  fi
  
  # Generate workspace name if not provided
  if [[ -z "${WORKSPACE_NAME}" ]]; then
    WORKSPACE_NAME="e2e-${MODE}-$(date +%s)"
  fi
  
  # Generate validation ID
  VALIDATION_ID="validation-${MODE}-$(date +%Y%m%d-%H%M%S)"
}

#######################################
# Check prerequisites
#######################################
check_prerequisites() {
  log_stage_start "Checking Prerequisites"
  
  log_info "PROJECT_ROOT: ${PROJECT_ROOT}"
  log_info "RESOURCE_DEFINITION: ${RESOURCE_DEFINITION}"
  log_info "Checking if ${RESOURCE_DEFINITION} exists..."
  
  require_command "go" "Go compiler required for building provider"
  require_command "jq" "jq required for JSON processing"
  require_file "${RESOURCE_DEFINITION}" "Resource definition file not found"
  
  # Check config directory
  if [[ ! -d "${CONFIG_DIR}/${MODE}" ]]; then
    log_error "Configuration directory not found: ${CONFIG_DIR}/${MODE}"
    exit "${EXIT_CONFIG_ERROR}"
  fi
  
  log_stage_end "Prerequisites Check" 0
}

#######################################
# Build provider binary
#######################################
build_provider() {
  if [[ "${SKIP_BUILD}" == "true" ]]; then
    if [[ -n "${PROVIDER_BINARY}" ]] && [[ -f "${PROVIDER_BINARY}" ]]; then
      log_info "Using pre-built provider: ${PROVIDER_BINARY}"
      return 0
    else
      log_error "Skip build specified but no valid provider binary found"
      exit "${EXIT_BUILD_FAILURE}"
    fi
  fi
  
  log_stage_start "Building Provider"
  record_stage_start "build" "${VALIDATION_ID}"
  
  local build_script="${SCRIPT_DIR}/build-provider.sh"
  require_file "${build_script}" "Build script not found"
  
  pushd "${PROJECT_ROOT}" > /dev/null
  
  if ! bash "${build_script}" > build.log 2>&1; then
    log_error "Provider build failed"
    cat build.log
    record_stage_end "build" "${VALIDATION_ID}" "${EXIT_BUILD_FAILURE}"
    popd > /dev/null
    exit "${EXIT_BUILD_FAILURE}"
  fi
  
  PROVIDER_BINARY="${PROJECT_ROOT}/terraform-provider-azurecaf"
  require_file "${PROVIDER_BINARY}" "Provider binary not found after build"
  
  popd > /dev/null
  
  record_stage_end "build" "${VALIDATION_ID}" 0
  log_stage_end "Provider Build" 0
}

#######################################
# Setup Terraform environment
#######################################
setup_terraform() {
  log_stage_start "Setting Up Terraform"
  record_stage_start "terraform_setup" "${VALIDATION_ID}"
  
  local setup_script="${SCRIPT_DIR}/setup-terraform.sh"
  require_file "${setup_script}" "Setup script not found"
  
  # Run setup script
  local setup_output
  if ! setup_output=$(bash "${setup_script}" \
    --version "${TERRAFORM_VERSION}" \
    --provider-binary "${PROVIDER_BINARY}" \
    --workspace-dir "${WORKSPACE_PATH}" 2>&1); then
    log_error "Terraform setup failed"
    echo "${setup_output}"
    record_stage_end "terraform_setup" "${VALIDATION_ID}" "${EXIT_DEPENDENCY_ERROR}"
    exit "${EXIT_DEPENDENCY_ERROR}"
  fi
  
  # Extract Terraform binary path from setup output
  TERRAFORM_BIN=$(echo "${setup_output}" | jq -r '.terraform_path')
  
  if [[ -z "${TERRAFORM_BIN}" ]] || [[ ! -x "${TERRAFORM_BIN}" ]]; then
    log_error "Terraform binary not found after setup"
    record_stage_end "terraform_setup" "${VALIDATION_ID}" "${EXIT_DEPENDENCY_ERROR}"
    exit "${EXIT_DEPENDENCY_ERROR}"
  fi
  
  log_info "Using Terraform: ${TERRAFORM_BIN}"
  
  record_stage_end "terraform_setup" "${VALIDATION_ID}" 0
  log_stage_end "Terraform Setup" 0
}

#######################################
# Initialize workspace
#######################################
init_workspace() {
  log_stage_start "Initializing Workspace"
  
  pushd "${WORKSPACE_PATH}" > /dev/null
  
  log_info "Running terraform init"
  if ! "${TERRAFORM_BIN}" init -no-color > init.log 2>&1; then
    log_error "Terraform init failed"
    cat init.log
    popd > /dev/null
    exit "${EXIT_PLAN_INIT}"
  fi
  
  log_info "Terraform initialized successfully"
  popd > /dev/null
  
  log_stage_end "Workspace Initialization" 0
}

#######################################
# Validate Terraform configuration
#######################################
validate_config() {
  local convention="$1"
  
  log_info "Validating ${convention} configuration"
  
  pushd "${WORKSPACE_PATH}" > /dev/null
  
  if ! "${TERRAFORM_BIN}" validate -no-color > "validate-${convention}.log" 2>&1; then
    log_error "Terraform validate failed for ${convention}"
    cat "validate-${convention}.log"
    popd > /dev/null
    return "${EXIT_PLAN_VALIDATION}"
  fi
  
  log_info "Configuration ${convention} is valid"
  popd > /dev/null
  
  return 0
}

#######################################
# Run Terraform plan
#######################################
run_plan() {
  local convention="$1"
  
  log_info "Running terraform plan for ${convention}"
  
  pushd "${WORKSPACE_PATH}" > /dev/null
  
  local plan_file="plan-${convention}.tfplan"
  local plan_json="plan-${convention}.json"
  
  # Run plan
  if ! "${TERRAFORM_BIN}" plan -no-color -out="${plan_file}" > "plan-${convention}.log" 2>&1; then
    log_error "Terraform plan failed for ${convention}"
    cat "plan-${convention}.log"
    popd > /dev/null
    return "${EXIT_PLAN_FAILURE}"
  fi
  
  # Convert plan to JSON
  if ! "${TERRAFORM_BIN}" show -json "${plan_file}" > "${plan_json}" 2>&1; then
    log_error "Failed to convert plan to JSON for ${convention}"
    popd > /dev/null
    return "${EXIT_PLAN_FAILURE}"
  fi
  
  log_info "Plan generated successfully for ${convention}"
  popd > /dev/null
  
  return 0
}

#######################################
# Validate generated names
#######################################
validate_names() {
  local convention="$1"
  local plan_json="${WORKSPACE_PATH}/plan-${convention}.json"
  
  log_info "Validating generated names for ${convention}"
  
  # Extract resource names from plan outputs
  local outputs
  outputs=$(jq -r '.planned_values.outputs // {}' "${plan_json}")
  
  if [[ "${outputs}" == "{}" ]]; then
    log_warn "No outputs found in plan for ${convention}"
    return 0
  fi
  
  # Get all_names output if it exists
  local all_names
  all_names=$(echo "${outputs}" | jq -r '.all_names.value // {}')
  
  if [[ "${all_names}" == "{}" ]]; then
    log_warn "No all_names output found for ${convention}"
    return 0
  fi
  
  # Validate each resource name
  local validation_errors=0
  local resource_type
  local generated_name
  local min_length
  local max_length
  local validation_regex
  
  while IFS= read -r resource_type; do
    generated_name=$(echo "${all_names}" | jq -r ".${resource_type}")
    
    if [[ "${generated_name}" == "null" ]] || [[ -z "${generated_name}" ]]; then
      continue
    fi
    
    # Get resource constraints from resourceDefinition.json
    local resource_def
    resource_def=$(jq -r ".[] | select(.name == \"${resource_type}\")" "${RESOURCE_DEFINITION}")
    
    if [[ -z "${resource_def}" ]]; then
      log_warn "No definition found for resource type: ${resource_type}"
      continue
    fi
    
    min_length=$(echo "${resource_def}" | jq -r '.min_length')
    max_length=$(echo "${resource_def}" | jq -r '.max_length')
    validation_regex=$(echo "${resource_def}" | jq -r '.validation_regex')
    
    # Validate length
    local name_length=${#generated_name}
    if [[ ${name_length} -lt ${min_length} ]] || [[ ${name_length} -gt ${max_length} ]]; then
      log_error "Name length validation failed for ${resource_type}: '${generated_name}' (length: ${name_length}, expected: ${min_length}-${max_length})"
      record_test_result "${VALIDATION_ID}" "${resource_type}" "${convention}" "fail" "Length validation failed"
      validation_errors=$((validation_errors + 1))
      continue
    fi
    
    # Validate regex pattern
    if [[ -n "${validation_regex}" ]] && [[ "${validation_regex}" != "null" ]]; then
      if ! echo "${generated_name}" | grep -qE "${validation_regex}"; then
        log_error "Regex validation failed for ${resource_type}: '${generated_name}' (pattern: ${validation_regex})"
        record_test_result "${VALIDATION_ID}" "${resource_type}" "${convention}" "fail" "Regex validation failed"
        validation_errors=$((validation_errors + 1))
        continue
      fi
    fi
    
    log_test_result "pass" "${resource_type}" "${generated_name}"
    record_test_result "${VALIDATION_ID}" "${resource_type}" "${convention}" "pass" ""
    
  done < <(echo "${all_names}" | jq -r 'keys[]')
  
  if [[ ${validation_errors} -gt 0 ]]; then
    log_error "Name validation failed with ${validation_errors} errors for ${convention}"
    return "${EXIT_PLAN_VALIDATION}"
  fi
  
  log_info "All names validated successfully for ${convention}"
  return 0
}

#######################################
# Test single convention
#######################################
test_convention() {
  local convention="$1"
  
  log_stage_start "Testing Convention: ${convention}"
  record_stage_start "plan_${convention}" "${VALIDATION_ID}"
  
  # Copy convention config to workspace
  local config_file="${CONFIG_DIR}/${MODE}/${convention}.tf"
  if [[ ! -f "${config_file}" ]]; then
    log_error "Configuration file not found: ${config_file}"
    record_stage_end "plan_${convention}" "${VALIDATION_ID}" "${EXIT_CONFIG_ERROR}"
    return "${EXIT_CONFIG_ERROR}"
  fi
  
  cp "${config_file}" "${WORKSPACE_PATH}/"
  
  # Validate configuration
  if ! validate_config "${convention}"; then
    log_error "Configuration validation failed for ${convention}"
    record_stage_end "plan_${convention}" "${VALIDATION_ID}" "${EXIT_PLAN_VALIDATION}"
    return "${EXIT_PLAN_VALIDATION}"
  fi
  
  # Run plan
  if ! run_plan "${convention}"; then
    log_error "Plan failed for ${convention}"
    record_stage_end "plan_${convention}" "${VALIDATION_ID}" "${EXIT_PLAN_FAILURE}"
    return "${EXIT_PLAN_FAILURE}"
  fi
  
  # Validate generated names
  if ! validate_names "${convention}"; then
    log_error "Name validation failed for ${convention}"
    record_stage_end "plan_${convention}" "${VALIDATION_ID}" "${EXIT_PLAN_VALIDATION}"
    return "${EXIT_PLAN_VALIDATION}"
  fi
  
  record_stage_end "plan_${convention}" "${VALIDATION_ID}" 0
  log_stage_end "Convention ${convention}" 0
  
  return 0
}

#######################################
# Run plan validation
#######################################
run_plan_validation() {
  log_stage_start "Plan Validation"
  record_stage_start "plan" "${VALIDATION_ID}"
  
  local failed_conventions=0
  
  for convention in "${CONVENTIONS[@]}"; do
    if ! test_convention "${convention}"; then
      log_error "Convention test failed: ${convention}"
      failed_conventions=$((failed_conventions + 1))
    fi
  done
  
  if [[ ${failed_conventions} -gt 0 ]]; then
    log_error "Plan validation failed for ${failed_conventions} conventions"
    record_stage_end "plan" "${VALIDATION_ID}" "${EXIT_PLAN_VALIDATION}"
    return "${EXIT_PLAN_VALIDATION}"
  fi
  
  record_stage_end "plan" "${VALIDATION_ID}" 0
  log_stage_end "Plan Validation" 0
  
  return 0
}

#######################################
# Run Terraform apply
#######################################
run_apply() {
  local convention="$1"
  
  log_info "Running terraform apply for ${convention}"
  
  pushd "${WORKSPACE_PATH}" > /dev/null
  
  local plan_file="plan-${convention}.tfplan"
  
  # Apply the plan
  if ! "${TERRAFORM_BIN}" apply -no-color -auto-approve "${plan_file}" > "apply-${convention}.log" 2>&1; then
    log_error "Terraform apply failed for ${convention}"
    cat "apply-${convention}.log"
    popd > /dev/null
    return "${EXIT_APPLY_FAILURE}"
  fi
  
  log_info "Apply successful for ${convention}"
  popd > /dev/null
  
  return 0
}

#######################################
# Run Terraform destroy
#######################################
run_destroy() {
  local convention="$1"
  
  log_info "Running terraform destroy for ${convention}"
  
  pushd "${WORKSPACE_PATH}" > /dev/null
  
  # Destroy resources
  if ! "${TERRAFORM_BIN}" destroy -no-color -auto-approve > "destroy-${convention}.log" 2>&1; then
    log_error "Terraform destroy failed for ${convention}"
    cat "destroy-${convention}.log"
    popd > /dev/null
    return "${EXIT_CLEANUP_FAILED}"
  fi
  
  log_info "Destroy successful for ${convention}"
  popd > /dev/null
  
  return 0
}

#######################################
# Run apply validation
#######################################
run_apply_validation() {
  if [[ "${SKIP_APPLY}" == "true" ]]; then
    log_info "Apply stage skipped (plan-only mode)"
    return 0
  fi
  
  log_stage_start "Apply Validation"
  record_stage_start "apply" "${VALIDATION_ID}"
  
  local failed_conventions=0
  
  for convention in "${CONVENTIONS[@]}"; do
    log_info "Applying configuration for ${convention}"
    
    if ! run_apply "${convention}"; then
      log_error "Apply failed for ${convention}"
      failed_conventions=$((failed_conventions + 1))
      continue
    fi
    
    # Destroy resources immediately to avoid conflicts
    if ! run_destroy "${convention}"; then
      log_warn "Destroy failed for ${convention} - resources may need manual cleanup"
    fi
  done
  
  if [[ ${failed_conventions} -gt 0 ]]; then
    log_error "Apply validation failed for ${failed_conventions} conventions"
    record_stage_end "apply" "${VALIDATION_ID}" "${EXIT_APPLY_FAILURE}"
    return "${EXIT_APPLY_FAILURE}"
  fi
  
  record_stage_end "apply" "${VALIDATION_ID}" 0
  log_stage_end "Apply Validation" 0
  
  return 0
}

#######################################
# Generate validation report
#######################################
generate_report() {
  log_stage_start "Generating Report"
  
  local report_json="${WORKSPACE_PATH}/validation-report.json"
  
  if ! generate_summary_report "${VALIDATION_ID}" > "${report_json}"; then
    log_error "Failed to generate validation report"
    return 1
  fi
  
  # Print summary
  log_info "Validation Report:"
  cat "${report_json}" | jq '.'
  
  # Calculate pass rate
  local pass_rate
  pass_rate=$(calculate_pass_rate "${VALIDATION_ID}")
  log_info "Overall Pass Rate: ${pass_rate}%"
  
  log_stage_end "Report Generation" 0
  
  return 0
}

#######################################
# Cleanup workspace
#######################################
cleanup() {
  if [[ "${SKIP_CLEANUP}" == "true" ]]; then
    log_info "Cleanup skipped, workspace preserved at: ${WORKSPACE_PATH}"
    return 0
  fi
  
  log_info "Cleaning up workspace"
  
  if [[ -n "${WORKSPACE_PATH}" ]] && [[ -d "${WORKSPACE_PATH}" ]]; then
    clean_workspace "${WORKSPACE_PATH}"
  fi
}

#######################################
# Main execution
#######################################
main() {
  log_info "Starting ${SCRIPT_NAME} v${SCRIPT_VERSION}"
  log_info "Validation ID: ${VALIDATION_ID}"
  log_info "Mode: ${MODE}"
  log_info "Conventions: ${CONVENTIONS[*]}"
  log_info "Terraform Version: ${TERRAFORM_VERSION}"
  
  # Setup trap for cleanup
  trap cleanup EXIT
  
  # Parse arguments
  parse_args "$@"
  
  # Check prerequisites
  check_prerequisites
  
  # Create workspace
  WORKSPACE_PATH=$(create_workspace "${WORKSPACE_NAME}")
  log_info "Workspace created: ${WORKSPACE_PATH}"
  
  # Copy provider config
  copy_config_to_workspace "${CONFIG_DIR}/${MODE}/provider.tf" "${WORKSPACE_PATH}"
  
  # Build provider
  build_provider
  
  # Setup Terraform
  setup_terraform
  
  # Initialize workspace
  init_workspace
  
  # Run plan validation
  local exit_code=0
  if ! run_plan_validation; then
    exit_code="${EXIT_PLAN_VALIDATION}"
  fi
  
  # Run apply validation (if enabled)
  if [[ ${exit_code} -eq 0 ]] && [[ "${SKIP_APPLY}" == "false" ]]; then
    if ! run_apply_validation; then
      exit_code="${EXIT_APPLY_FAILURE}"
    fi
  fi
  
  # Generate report
  generate_report
  
  # Exit with appropriate code
  if [[ ${exit_code} -ne 0 ]]; then
    log_error "Validation failed with exit code: ${exit_code}"
    exit "${exit_code}"
  fi
  
  log_info "${SCRIPT_NAME} completed successfully"
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
