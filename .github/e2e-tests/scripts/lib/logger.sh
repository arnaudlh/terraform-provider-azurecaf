#!/usr/bin/env bash
# Logging Framework for E2E Validation
# Provides structured logging with levels and timestamps

set -euo pipefail

# Prevent double-sourcing
if [[ -n "${_LOGGER_SH_SOURCED:-}" ]]; then
  return 0
fi
readonly _LOGGER_SH_SOURCED=1

# Log levels
if [[ -z "${LOG_LEVEL_DEBUG:-}" ]]; then
  readonly LOG_LEVEL_DEBUG=0
  readonly LOG_LEVEL_INFO=1
  readonly LOG_LEVEL_WARN=2
  readonly LOG_LEVEL_ERROR=3
fi

# Current log level (default: INFO)
LOG_LEVEL="${LOG_LEVEL:-$LOG_LEVEL_INFO}"

# Log file (optional)
LOG_FILE="${LOG_FILE:-}"

# Colors for terminal output
readonly COLOR_RESET="\033[0m"
readonly COLOR_DEBUG="\033[0;36m"    # Cyan
readonly COLOR_INFO="\033[0;32m"     # Green
readonly COLOR_WARN="\033[0;33m"     # Yellow
readonly COLOR_ERROR="\033[0;31m"    # Red

# Get timestamp
get_timestamp() {
    date '+%Y-%m-%d %H:%M:%S'
}

# Log message with level
# Usage: log <level> <message>
log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp
    timestamp=$(get_timestamp)
    
    # Check if we should log this level
    case "$level" in
        DEBUG)
            [[ $LOG_LEVEL -gt $LOG_LEVEL_DEBUG ]] && return 0
            local color="$COLOR_DEBUG"
            ;;
        INFO)
            [[ $LOG_LEVEL -gt $LOG_LEVEL_INFO ]] && return 0
            local color="$COLOR_INFO"
            ;;
        WARN)
            [[ $LOG_LEVEL -gt $LOG_LEVEL_WARN ]] && return 0
            local color="$COLOR_WARN"
            ;;
        ERROR)
            [[ $LOG_LEVEL -gt $LOG_LEVEL_ERROR ]] && return 0
            local color="$COLOR_ERROR"
            ;;
        *)
            local color="$COLOR_RESET"
            ;;
    esac
    
    # Format log message
    local log_msg="[$timestamp] [$level] $message"
    
    # Output to terminal with color
    if [[ -t 1 ]]; then
        echo -e "${color}${log_msg}${COLOR_RESET}"
    else
        echo "$log_msg"
    fi
    
    # Output to log file if specified
    if [[ -n "$LOG_FILE" ]]; then
        echo "$log_msg" >> "$LOG_FILE"
    fi
}

# Convenience functions
log_debug() {
    log DEBUG "$@"
}

log_info() {
    log INFO "$@"
}

log_warn() {
    log WARN "$@"
}

log_error() {
    log ERROR "$@"
}

# Log command execution
# Usage: log_cmd <command> [args...]
log_cmd() {
    log_debug "Executing: $*"
    "$@"
}

# Log stage start
# Usage: log_stage_start <stage_name>
log_stage_start() {
    local stage="$1"
    log_info "========================================="
    log_info "Stage: $stage - STARTING"
    log_info "========================================="
}

# Log stage end
# Usage: log_stage_end <stage_name> <status>
log_stage_end() {
    local stage="$1"
    local status="$2"
    
    if [[ "$status" == "success" ]]; then
        log_info "========================================="
        log_info "Stage: $stage - COMPLETED SUCCESSFULLY"
        log_info "========================================="
    else
        log_error "========================================="
        log_error "Stage: $stage - FAILED"
        log_error "========================================="
    fi
}

# Log test result
# Usage: log_test_result <test_name> <status> [details]
log_test_result() {
    local test_name="$1"
    local status="$2"
    local details="${3:-}"
    
    if [[ "$status" == "pass" ]]; then
        log_info "✓ $test_name - PASSED"
    elif [[ "$status" == "fail" ]]; then
        log_error "✗ $test_name - FAILED"
        if [[ -n "$details" ]]; then
            log_error "  Details: $details"
        fi
    elif [[ "$status" == "skip" ]]; then
        log_warn "⊘ $test_name - SKIPPED"
        if [[ -n "$details" ]]; then
            log_warn "  Reason: $details"
        fi
    fi
}

# Initialize logging
# Usage: init_logging [log_file]
init_logging() {
    local log_file="${1:-}"
    
    if [[ -n "$log_file" ]]; then
        LOG_FILE="$log_file"
        # Create log directory if needed
        mkdir -p "$(dirname "$LOG_FILE")"
        # Initialize log file
        echo "=== E2E Validation Log ===" > "$LOG_FILE"
        echo "Started: $(get_timestamp)" >> "$LOG_FILE"
        echo "" >> "$LOG_FILE"
        log_info "Logging initialized: $LOG_FILE"
    fi
}

# Export functions
export -f get_timestamp
export -f log
export -f log_debug
export -f log_info
export -f log_warn
export -f log_error
export -f log_cmd
export -f log_stage_start
export -f log_stage_end
export -f log_test_result
export -f init_logging
