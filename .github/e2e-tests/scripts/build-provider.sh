#!/usr/bin/env bash
# Build Provider Script for E2E Validation
# Compiles the Terraform provider binary and prepares it for testing

set -euo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

# Source library functions
# shellcheck source=./lib/logger.sh
source "${SCRIPT_DIR}/lib/logger.sh"
# shellcheck source=./lib/error-handler.sh
source "${SCRIPT_DIR}/lib/error-handler.sh"

# Configuration
OUTPUT_DIR="${OUTPUT_DIR:-${REPO_ROOT}}"
BINARY_NAME="terraform-provider-azurecaf"
GO_VERSION="${GO_VERSION:-1.24}"

# Initialize logging
LOG_FILE="${LOG_FILE:-/tmp/e2e-build-$(date +%Y%m%d-%H%M%S).log}"
init_logging "$LOG_FILE"

log_info "========================================="
log_info "E2E Validation - Provider Build"
log_info "========================================="

# Verify prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    require_command go
    require_command jq
    
    # Verify Go version
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go version: $go_version"
    
    # Check if go.mod exists
    require_file "${REPO_ROOT}/go.mod"
    
    log_info "Prerequisites check passed"
}

# Clean previous builds
clean_previous_builds() {
    log_info "Cleaning previous builds..."
    
    # Remove previous binary
    if [[ -f "${OUTPUT_DIR}/${BINARY_NAME}" ]]; then
        rm -f "${OUTPUT_DIR}/${BINARY_NAME}"
        log_debug "Removed previous binary"
    fi
    
    # Clean Go build cache (optional, can be skipped for faster builds)
    if [[ "${CLEAN_CACHE:-false}" == "true" ]]; then
        log_info "Cleaning Go build cache..."
        go clean -cache
    fi
}

# Download dependencies
download_dependencies() {
    log_info "Downloading Go dependencies..."
    
    cd "$REPO_ROOT" || exit 1
    
    local start_time
    start_time=$(date +%s)
    
    if go mod download; then
        local end_time
        end_time=$(date +%s)
        local duration=$((end_time - start_time))
        log_info "Dependencies downloaded successfully in ${duration}s"
    else
        handle_error $EXIT_DEPENDENCY_ERROR "Failed to download Go dependencies"
        return $EXIT_DEPENDENCY_ERROR
    fi
}

# Build provider binary
build_provider() {
    log_info "Building provider binary..."
    
    cd "$REPO_ROOT" || exit 1
    
    local start_time
    start_time=$(date +%s)
    
    # Build flags
    local build_flags=(
        -o "${OUTPUT_DIR}/${BINARY_NAME}"
    )
    
    # Add debug flags if requested
    if [[ "${DEBUG_BUILD:-false}" == "true" ]]; then
        build_flags+=(-gcflags="all=-N -l")
    fi
    
    log_debug "Build command: go build ${build_flags[*]}"
    
    if go build "${build_flags[@]}"; then
        local end_time
        end_time=$(date +%s)
        local duration=$((end_time - start_time))
        
        log_info "Provider binary built successfully in ${duration}s"
        echo "$duration" > "${OUTPUT_DIR}/.build-duration"
        
        return 0
    else
        handle_error $EXIT_BUILD_FAILURE "Failed to build provider binary"
        return $EXIT_BUILD_FAILURE
    fi
}

# Verify binary
verify_binary() {
    log_info "Verifying provider binary..."
    
    local binary_path="${OUTPUT_DIR}/${BINARY_NAME}"
    
    # Check if binary exists
    if [[ ! -f "$binary_path" ]]; then
        handle_error $EXIT_BUILD_FAILURE "Provider binary not found: $binary_path"
        return $EXIT_BUILD_FAILURE
    fi
    
    # Check if binary is executable
    if [[ ! -x "$binary_path" ]]; then
        log_warn "Binary is not executable, fixing permissions..."
        chmod +x "$binary_path"
    fi
    
    # Get binary size
    local binary_size
    binary_size=$(stat -f%z "$binary_path" 2>/dev/null || stat -c%s "$binary_path" 2>/dev/null || echo "unknown")
    log_info "Binary size: $binary_size bytes"
    
    # Calculate SHA256 hash
    local binary_hash
    if command -v sha256sum >/dev/null 2>&1; then
        binary_hash=$(sha256sum "$binary_path" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        binary_hash=$(shasum -a 256 "$binary_path" | awk '{print $1}')
    else
        log_warn "SHA256 command not found, skipping hash calculation"
        binary_hash="N/A"
    fi
    
    log_info "Binary SHA256: $binary_hash"
    echo "$binary_hash" > "${OUTPUT_DIR}/.build-hash"
    
    log_info "Binary verification passed"
}

# Generate build metadata
generate_build_metadata() {
    log_info "Generating build metadata..."
    
    local metadata_file="${OUTPUT_DIR}/.build-metadata.json"
    local binary_path="${OUTPUT_DIR}/${BINARY_NAME}"
    
    local build_time
    build_time=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
    
    local build_duration
    build_duration=$(cat "${OUTPUT_DIR}/.build-duration" 2>/dev/null || echo "0")
    
    local binary_hash
    binary_hash=$(cat "${OUTPUT_DIR}/.build-hash" 2>/dev/null || echo "unknown")
    
    local binary_size
    binary_size=$(stat -f%z "$binary_path" 2>/dev/null || stat -c%s "$binary_path" 2>/dev/null || echo "0")
    
    local go_version
    go_version=$(go version | awk '{print $3}')
    
    local git_commit
    git_commit=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
    
    local git_branch
    git_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    
    cat > "$metadata_file" <<EOF
{
  "build_time": "$build_time",
  "build_duration": $build_duration,
  "binary_path": "$binary_path",
  "binary_size": $binary_size,
  "binary_hash": "$binary_hash",
  "go_version": "$go_version",
  "git_commit": "$git_commit",
  "git_branch": "$git_branch"
}
EOF
    
    log_info "Build metadata written to: $metadata_file"
    
    # Display metadata
    log_info "Build Metadata:"
    jq '.' "$metadata_file" || cat "$metadata_file"
}

# Main execution
main() {
    local exit_code=0
    
    # Record start time
    local overall_start
    overall_start=$(date +%s)
    
    # Execute build steps
    check_prerequisites || exit_code=$?
    
    if [[ $exit_code -eq 0 ]]; then
        clean_previous_builds || exit_code=$?
    fi
    
    if [[ $exit_code -eq 0 ]]; then
        download_dependencies || exit_code=$?
    fi
    
    if [[ $exit_code -eq 0 ]]; then
        build_provider || exit_code=$?
    fi
    
    if [[ $exit_code -eq 0 ]]; then
        verify_binary || exit_code=$?
    fi
    
    if [[ $exit_code -eq 0 ]]; then
        generate_build_metadata || exit_code=$?
    fi
    
    # Calculate total duration
    local overall_end
    overall_end=$(date +%s)
    local total_duration=$((overall_end - overall_start))
    
    # Report result
    if [[ $exit_code -eq 0 ]]; then
        log_info "========================================="
        log_info "Provider build completed successfully!"
        log_info "Total duration: ${total_duration}s"
        log_info "Binary: ${OUTPUT_DIR}/${BINARY_NAME}"
        log_info "========================================="
    else
        log_error "========================================="
        log_error "Provider build FAILED"
        log_error "Exit code: $exit_code"
        log_error "Total duration: ${total_duration}s"
        log_error "========================================="
    fi
    
    return $exit_code
}

# Execute main function
main "$@"
