#!/usr/bin/env bash
# Error Handling Utilities for E2E Validation
# Provides error handling, exit codes, and recovery mechanisms

set -euo pipefail

# Prevent double-sourcing
if [[ -n "${_ERROR_HANDLER_SH_SOURCED:-}" ]]; then
  return 0
fi
readonly _ERROR_HANDLER_SH_SOURCED=1

# Exit codes (matching contracts/validation-pipeline.md)
if [[ -z "${EXIT_SUCCESS:-}" ]]; then
  readonly EXIT_SUCCESS=0
  readonly EXIT_BUILD_FAILURE=1
  readonly EXIT_DEPENDENCY_ERROR=2
  readonly EXIT_INIT_ERROR=10
  readonly EXIT_PLAN_ERROR=11
  readonly EXIT_VALIDATION_ERROR=12
  readonly EXIT_APPLY_ERROR=20
  readonly EXIT_RESOURCE_ERROR=21
  readonly EXIT_STATE_ERROR=22
  readonly EXIT_DRIFT_DETECTED=30
  readonly EXIT_DRIFT_PLAN_ERROR=31
  readonly EXIT_CLEANUP_ERROR=40
  readonly EXIT_RESOURCE_STUCK=41
  readonly EXIT_TIMEOUT=124
  readonly EXIT_CONFIG_ERROR=125
  readonly EXIT_PIPELINE_ERROR=126
fi

# Source logger if available
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ -f "${SCRIPT_DIR}/logger.sh" ]]; then
    # shellcheck source=./logger.sh
    source "${SCRIPT_DIR}/logger.sh"
else
    # Fallback logging functions
    log_error() { echo "ERROR: $*" >&2; }
    log_warn() { echo "WARN: $*" >&2; }
    log_info() { echo "INFO: $*"; }
fi

# Error context stack
declare -a ERROR_CONTEXT_STACK=()

# Push error context
# Usage: push_error_context <context>
push_error_context() {
    local context="$1"
    ERROR_CONTEXT_STACK+=("$context")
}

# Pop error context
# Usage: pop_error_context
pop_error_context() {
    if [[ ${#ERROR_CONTEXT_STACK[@]} -gt 0 ]]; then
        unset 'ERROR_CONTEXT_STACK[${#ERROR_CONTEXT_STACK[@]}-1]'
    fi
}

# Get current error context
# Usage: get_error_context
get_error_context() {
    if [[ ${#ERROR_CONTEXT_STACK[@]} -gt 0 ]]; then
        echo "${ERROR_CONTEXT_STACK[@]}"
    else
        echo "unknown"
    fi
}

# Handle error with context
# Usage: handle_error <exit_code> <error_message>
handle_error() {
    local exit_code="$1"
    local error_message="$2"
    local context
    context=$(get_error_context)
    
    log_error "Error in context: $context"
    log_error "Error message: $error_message"
    log_error "Exit code: $exit_code"
    
    # Get exit code description
    local exit_desc
    exit_desc=$(get_exit_code_description "$exit_code")
    log_error "Exit code description: $exit_desc"
    
    return "$exit_code"
}

# Get exit code description
# Usage: get_exit_code_description <exit_code>
get_exit_code_description() {
    local code="$1"
    
    case "$code" in
        $EXIT_SUCCESS) echo "Success" ;;
        $EXIT_BUILD_FAILURE) echo "Build failure" ;;
        $EXIT_DEPENDENCY_ERROR) echo "Dependency error" ;;
        $EXIT_INIT_ERROR) echo "Init error" ;;
        $EXIT_PLAN_ERROR) echo "Plan error" ;;
        $EXIT_VALIDATION_ERROR) echo "Validation error" ;;
        $EXIT_APPLY_ERROR) echo "Apply error" ;;
        $EXIT_RESOURCE_ERROR) echo "Resource error" ;;
        $EXIT_STATE_ERROR) echo "State error" ;;
        $EXIT_DRIFT_DETECTED) echo "Drift detected (CRITICAL)" ;;
        $EXIT_DRIFT_PLAN_ERROR) echo "Drift check plan error" ;;
        $EXIT_CLEANUP_ERROR) echo "Cleanup error (WARNING)" ;;
        $EXIT_RESOURCE_STUCK) echo "Resource stuck (WARNING)" ;;
        $EXIT_TIMEOUT) echo "Timeout" ;;
        $EXIT_CONFIG_ERROR) echo "Configuration error" ;;
        $EXIT_PIPELINE_ERROR) echo "Pipeline error" ;;
        *) echo "Unknown error code: $code" ;;
    esac
}

# Retry command with exponential backoff
# Usage: retry_command <max_attempts> <command> [args...]
retry_command() {
    local max_attempts="$1"
    shift
    local attempt=1
    local delay=1
    
    while [[ $attempt -le $max_attempts ]]; do
        log_info "Attempt $attempt/$max_attempts: $*"
        
        if "$@"; then
            log_info "Command succeeded on attempt $attempt"
            return 0
        fi
        
        if [[ $attempt -lt $max_attempts ]]; then
            log_warn "Command failed on attempt $attempt, retrying in ${delay}s..."
            sleep "$delay"
            delay=$((delay * 2))  # Exponential backoff
        fi
        
        attempt=$((attempt + 1))
    done
    
    log_error "Command failed after $max_attempts attempts"
    return 1
}

# Execute command with timeout
# Usage: timeout_command <timeout_seconds> <command> [args...]
timeout_command() {
    local timeout_seconds="$1"
    shift
    
    log_debug "Executing with timeout ${timeout_seconds}s: $*"
    
    if timeout "$timeout_seconds" "$@"; then
        return 0
    else
        local exit_code=$?
        if [[ $exit_code -eq 124 ]]; then
            log_error "Command timed out after ${timeout_seconds}s"
            return $EXIT_TIMEOUT
        else
            log_error "Command failed with exit code: $exit_code"
            return "$exit_code"
        fi
    fi
}

# Check command exists
# Usage: require_command <command>
require_command() {
    local cmd="$1"
    
    if ! command -v "$cmd" >/dev/null 2>&1; then
        handle_error $EXIT_DEPENDENCY_ERROR "Required command not found: $cmd"
        return $EXIT_DEPENDENCY_ERROR
    fi
}

# Check file exists
# Usage: require_file <file_path>
require_file() {
    local file_path="$1"
    
    if [[ ! -f "$file_path" ]]; then
        handle_error $EXIT_CONFIG_ERROR "Required file not found: $file_path"
        return $EXIT_CONFIG_ERROR
    fi
}

# Validate environment variable is set
# Usage: require_env_var <var_name>
require_env_var() {
    local var_name="$1"
    
    if [[ -z "${!var_name:-}" ]]; then
        handle_error $EXIT_CONFIG_ERROR "Required environment variable not set: $var_name"
        return $EXIT_CONFIG_ERROR
    fi
}

# Trap errors and cleanup
# Usage: setup_error_trap <cleanup_function>
setup_error_trap() {
    local cleanup_function="$1"
    
    trap "$cleanup_function" ERR EXIT
}

# Export exit codes
export EXIT_SUCCESS
export EXIT_BUILD_FAILURE
export EXIT_DEPENDENCY_ERROR
export EXIT_INIT_ERROR
export EXIT_PLAN_ERROR
export EXIT_VALIDATION_ERROR
export EXIT_APPLY_ERROR
export EXIT_RESOURCE_ERROR
export EXIT_STATE_ERROR
export EXIT_DRIFT_DETECTED
export EXIT_DRIFT_PLAN_ERROR
export EXIT_CLEANUP_ERROR
export EXIT_RESOURCE_STUCK
export EXIT_TIMEOUT
export EXIT_CONFIG_ERROR
export EXIT_PIPELINE_ERROR

# Export functions
export -f push_error_context
export -f pop_error_context
export -f get_error_context
export -f handle_error
export -f get_exit_code_description
export -f retry_command
export -f timeout_command
export -f require_command
export -f require_file
export -f require_env_var
export -f setup_error_trap
