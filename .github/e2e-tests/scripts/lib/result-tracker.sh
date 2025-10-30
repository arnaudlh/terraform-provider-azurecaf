#!/usr/bin/env bash
# Stage Result Tracker for E2E Validation
# Tracks and reports results for each pipeline stage

set -euo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source logger
# shellcheck source=./logger.sh
source "${SCRIPT_DIR}/logger.sh"

# Results file
RESULTS_FILE="${RESULTS_FILE:-}"
RESULTS_JSON="${RESULTS_JSON:-}"

# Stage result structure
declare -A STAGE_RESULTS

# Initialize result tracking
# Usage: init_result_tracking <results_file>
init_result_tracking() {
    local results_file="$1"
    RESULTS_FILE="$results_file"
    RESULTS_JSON="${results_file}.json"
    
    # Create results directory
    mkdir -p "$(dirname "$RESULTS_FILE")"
    
    # Initialize results file
    {
        echo "=== E2E Validation Results ==="
        echo "Started: $(date '+%Y-%m-%d %H:%M:%S')"
        echo ""
    } > "$RESULTS_FILE"
    
    # Initialize JSON results
    echo '{"validation_id":"","status":"running","stages":{},"tests":[]}' > "$RESULTS_JSON"
    
    log_info "Result tracking initialized: $RESULTS_FILE"
}

# Record stage start
# Usage: record_stage_start <stage_name>
record_stage_start() {
    local stage_name="$1"
    local start_time
    start_time=$(date +%s)
    
    STAGE_RESULTS["${stage_name}_start"]="$start_time"
    STAGE_RESULTS["${stage_name}_status"]="running"
    
    log_debug "Stage start recorded: $stage_name at $start_time"
}

# Record stage end
# Usage: record_stage_end <stage_name> <status> [error_message]
record_stage_end() {
    local stage_name="$1"
    local status="$2"
    local error_message="${3:-}"
    local end_time
    end_time=$(date +%s)
    
    STAGE_RESULTS["${stage_name}_end"]="$end_time"
    STAGE_RESULTS["${stage_name}_status"]="$status"
    
    if [[ -n "$error_message" ]]; then
        STAGE_RESULTS["${stage_name}_error"]="$error_message"
    fi
    
    # Calculate duration
    local start_time="${STAGE_RESULTS[${stage_name}_start]:-$end_time}"
    local duration=$((end_time - start_time))
    STAGE_RESULTS["${stage_name}_duration"]="$duration"
    
    log_debug "Stage end recorded: $stage_name - $status (${duration}s)"
    
    # Write to results file
    if [[ -n "$RESULTS_FILE" ]]; then
        {
            echo "Stage: $stage_name"
            echo "  Status: $status"
            echo "  Duration: ${duration}s"
            if [[ -n "$error_message" ]]; then
                echo "  Error: $error_message"
            fi
            echo ""
        } >> "$RESULTS_FILE"
    fi
    
    # Update JSON results
    if [[ -n "$RESULTS_JSON" ]]; then
        update_json_stage_result "$stage_name" "$status" "$duration" "$error_message"
    fi
}

# Update JSON stage result
# Usage: update_json_stage_result <stage_name> <status> <duration> <error>
update_json_stage_result() {
    local stage_name="$1"
    local status="$2"
    local duration="$3"
    local error="${4:-}"
    
    if [[ ! -f "$RESULTS_JSON" ]]; then
        return 0
    fi
    
    local json_update
    if [[ -n "$error" ]]; then
        json_update=$(jq --arg stage "$stage_name" \
                        --arg status "$status" \
                        --arg duration "$duration" \
                        --arg error "$error" \
                        '.stages[$stage] = {status: $status, duration: ($duration | tonumber), error: $error}' \
                        "$RESULTS_JSON")
    else
        json_update=$(jq --arg stage "$stage_name" \
                        --arg status "$status" \
                        --arg duration "$duration" \
                        '.stages[$stage] = {status: $status, duration: ($duration | tonumber)}' \
                        "$RESULTS_JSON")
    fi
    
    echo "$json_update" > "$RESULTS_JSON"
}

# Record test result
# Usage: record_test_result <test_name> <status> <resource_type> <convention> [error]
record_test_result() {
    local test_name="$1"
    local status="$2"
    local resource_type="$3"
    local convention="$4"
    local error="${5:-}"
    
    log_debug "Test result recorded: $test_name - $status"
    
    # Write to results file
    if [[ -n "$RESULTS_FILE" ]]; then
        {
            echo "Test: $test_name"
            echo "  Resource: $resource_type"
            echo "  Convention: $convention"
            echo "  Status: $status"
            if [[ -n "$error" ]]; then
                echo "  Error: $error"
            fi
            echo ""
        } >> "$RESULTS_FILE"
    fi
    
    # Update JSON results
    if [[ -n "$RESULTS_JSON" ]]; then
        update_json_test_result "$test_name" "$status" "$resource_type" "$convention" "$error"
    fi
}

# Update JSON test result
# Usage: update_json_test_result <test_name> <status> <resource_type> <convention> <error>
update_json_test_result() {
    local test_name="$1"
    local status="$2"
    local resource_type="$3"
    local convention="$4"
    local error="${5:-}"
    
    if [[ ! -f "$RESULTS_JSON" ]]; then
        return 0
    fi
    
    local test_obj
    if [[ -n "$error" ]]; then
        test_obj=$(jq -n --arg name "$test_name" \
                        --arg status "$status" \
                        --arg resource "$resource_type" \
                        --arg convention "$convention" \
                        --arg error "$error" \
                        '{name: $name, status: $status, resource_type: $resource, convention: $convention, error: $error}')
    else
        test_obj=$(jq -n --arg name "$test_name" \
                        --arg status "$status" \
                        --arg resource "$resource_type" \
                        --arg convention "$convention" \
                        '{name: $name, status: $status, resource_type: $resource, convention: $convention}')
    fi
    
    local json_update
    json_update=$(jq --argjson test "$test_obj" '.tests += [$test]' "$RESULTS_JSON")
    echo "$json_update" > "$RESULTS_JSON"
}

# Get stage status
# Usage: get_stage_status <stage_name>
get_stage_status() {
    local stage_name="$1"
    echo "${STAGE_RESULTS[${stage_name}_status]:-unknown}"
}

# Get stage duration
# Usage: get_stage_duration <stage_name>
get_stage_duration() {
    local stage_name="$1"
    echo "${STAGE_RESULTS[${stage_name}_duration]:-0}"
}

# Calculate pass rate
# Usage: calculate_pass_rate
calculate_pass_rate() {
    if [[ ! -f "$RESULTS_JSON" ]]; then
        echo "0"
        return
    fi
    
    local total
    total=$(jq '.tests | length' "$RESULTS_JSON")
    
    if [[ $total -eq 0 ]]; then
        echo "0"
        return
    fi
    
    local passed
    passed=$(jq '.tests | map(select(.status == "pass")) | length' "$RESULTS_JSON")
    
    echo "scale=2; $passed * 100 / $total" | bc
}

# Generate summary report
# Usage: generate_summary_report
generate_summary_report() {
    log_info "Generating summary report"
    
    if [[ -n "$RESULTS_FILE" ]]; then
        {
            echo ""
            echo "=== Summary ==="
            echo ""
            
            # Stage results
            echo "Stages:"
            for stage in build plan apply drift_check cleanup; do
                local status
                status=$(get_stage_status "$stage")
                if [[ "$status" != "unknown" ]]; then
                    local duration
                    duration=$(get_stage_duration "$stage")
                    echo "  $stage: $status (${duration}s)"
                fi
            done
            
            echo ""
            
            # Test results
            if [[ -f "$RESULTS_JSON" ]]; then
                local total
                total=$(jq '.tests | length' "$RESULTS_JSON")
                local passed
                passed=$(jq '.tests | map(select(.status == "pass")) | length' "$RESULTS_JSON")
                local failed
                failed=$(jq '.tests | map(select(.status == "fail")) | length' "$RESULTS_JSON")
                local pass_rate
                pass_rate=$(calculate_pass_rate)
                
                echo "Tests:"
                echo "  Total: $total"
                echo "  Passed: $passed"
                echo "  Failed: $failed"
                echo "  Pass Rate: ${pass_rate}%"
            fi
            
            echo ""
            echo "Completed: $(date '+%Y-%m-%d %H:%M:%S')"
        } >> "$RESULTS_FILE"
        
        log_info "Summary report written to: $RESULTS_FILE"
    fi
}

# Export functions
export -f init_result_tracking
export -f record_stage_start
export -f record_stage_end
export -f record_test_result
export -f get_stage_status
export -f get_stage_duration
export -f calculate_pass_rate
export -f generate_summary_report
