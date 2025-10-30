#!/usr/bin/env bash
# Test Workspace Manager for E2E Validation
# Manages temporary test workspaces and state files

set -euo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source logger and error handler
# shellcheck source=./logger.sh
source "${SCRIPT_DIR}/logger.sh"
# shellcheck source=./error-handler.sh
source "${SCRIPT_DIR}/error-handler.sh"

# Base directory for test workspaces
WORKSPACE_BASE_DIR="${WORKSPACE_BASE_DIR:-/tmp/e2e-validation}"

# Create test workspace
# Usage: create_workspace <workspace_name>
create_workspace() {
    local workspace_name="$1"
    local timestamp
    timestamp=$(date +%Y%m%d-%H%M%S)
    local workspace_dir="${WORKSPACE_BASE_DIR}-${timestamp}/${workspace_name}"
    
    log_debug "Creating test workspace: $workspace_dir"
    
    if ! mkdir -p "$workspace_dir"; then
        handle_error $EXIT_PIPELINE_ERROR "Failed to create workspace directory: $workspace_dir"
        return $EXIT_PIPELINE_ERROR
    fi
    
    # Create subdirectories
    mkdir -p "$workspace_dir"/{configs,state,plans,logs}
    
    log_info "Test workspace created: $workspace_dir"
    echo "$workspace_dir"
}

# Initialize Terraform workspace
# Usage: init_terraform_workspace <workspace_dir>
init_terraform_workspace() {
    local workspace_dir="$1"
    
    log_debug "Initializing Terraform workspace: $workspace_dir"
    
    cd "$workspace_dir" || {
        handle_error $EXIT_PIPELINE_ERROR "Failed to change to workspace directory: $workspace_dir"
        return $EXIT_PIPELINE_ERROR
    }
    
    # Create .terraform directory
    mkdir -p .terraform
    
    log_info "Terraform workspace initialized: $workspace_dir"
}

# Copy configuration to workspace
# Usage: copy_config_to_workspace <config_file> <workspace_dir>
copy_config_to_workspace() {
    local config_file="$1"
    local workspace_dir="$2"
    
    log_debug "Copying configuration to workspace: $config_file -> $workspace_dir"
    
    if [[ ! -f "$config_file" ]]; then
        handle_error $EXIT_CONFIG_ERROR "Configuration file not found: $config_file"
        return $EXIT_CONFIG_ERROR
    fi
    
    if ! cp "$config_file" "$workspace_dir/"; then
        handle_error $EXIT_PIPELINE_ERROR "Failed to copy configuration: $config_file"
        return $EXIT_PIPELINE_ERROR
    fi
    
    log_debug "Configuration copied successfully"
}

# Copy provider binary to workspace
# Usage: copy_provider_to_workspace <provider_binary> <workspace_dir>
copy_provider_to_workspace() {
    local provider_binary="$1"
    local workspace_dir="$2"
    local provider_dir="${workspace_dir}/.terraform/plugins/local/aztfmod/azurecaf/9.9.9"
    
    log_debug "Installing provider binary to workspace"
    
    if [[ ! -f "$provider_binary" ]]; then
        handle_error $EXIT_BUILD_FAILURE "Provider binary not found: $provider_binary"
        return $EXIT_BUILD_FAILURE
    fi
    
    # Determine OS and architecture
    local os_arch
    os_arch=$(go env GOOS)_$(go env GOARCH)
    
    # Create provider directory
    mkdir -p "${provider_dir}/${os_arch}"
    
    # Copy provider binary
    if ! cp "$provider_binary" "${provider_dir}/${os_arch}/terraform-provider-azurecaf"; then
        handle_error $EXIT_PIPELINE_ERROR "Failed to install provider binary"
        return $EXIT_PIPELINE_ERROR
    fi
    
    # Make executable
    chmod +x "${provider_dir}/${os_arch}/terraform-provider-azurecaf"
    
    log_info "Provider binary installed: ${provider_dir}/${os_arch}/"
}

# Get workspace state file path
# Usage: get_state_file_path <workspace_dir>
get_state_file_path() {
    local workspace_dir="$1"
    echo "${workspace_dir}/terraform.tfstate"
}

# Get workspace plan file path
# Usage: get_plan_file_path <workspace_dir>
get_plan_file_path() {
    local workspace_dir="$1"
    echo "${workspace_dir}/terraform.tfplan"
}

# Check if workspace has state
# Usage: has_state <workspace_dir>
has_state() {
    local workspace_dir="$1"
    local state_file
    state_file=$(get_state_file_path "$workspace_dir")
    
    [[ -f "$state_file" ]]
}

# Get state file size
# Usage: get_state_size <workspace_dir>
get_state_size() {
    local workspace_dir="$1"
    local state_file
    state_file=$(get_state_file_path "$workspace_dir")
    
    if [[ -f "$state_file" ]]; then
        stat -f%z "$state_file" 2>/dev/null || stat -c%s "$state_file" 2>/dev/null
    else
        echo "0"
    fi
}

# Clean workspace
# Usage: clean_workspace <workspace_dir> [force]
clean_workspace() {
    local workspace_dir="$1"
    local force="${2:-false}"
    
    log_debug "Cleaning workspace: $workspace_dir"
    
    if [[ ! -d "$workspace_dir" ]]; then
        log_warn "Workspace directory does not exist: $workspace_dir"
        return 0
    fi
    
    # Remove .terraform directory
    if [[ -d "${workspace_dir}/.terraform" ]]; then
        rm -rf "${workspace_dir}/.terraform"
    fi
    
    # Remove state files
    rm -f "${workspace_dir}"/terraform.tfstate*
    
    # Remove plan files
    rm -f "${workspace_dir}"/terraform.tfplan
    
    # Remove logs
    rm -f "${workspace_dir}"/logs/*.log
    
    if [[ "$force" == "true" ]]; then
        log_warn "Force cleaning entire workspace directory"
        rm -rf "$workspace_dir"
    fi
    
    log_info "Workspace cleaned: $workspace_dir"
}

# List all workspaces
# Usage: list_workspaces
list_workspaces() {
    if [[ -d "$WORKSPACE_BASE_DIR" ]]; then
        find "${WORKSPACE_BASE_DIR}"* -maxdepth 0 -type d 2>/dev/null || true
    fi
}

# Clean old workspaces
# Usage: clean_old_workspaces <days>
clean_old_workspaces() {
    local days="${1:-7}"
    
    log_info "Cleaning workspaces older than $days days"
    
    if [[ -d "$(dirname "$WORKSPACE_BASE_DIR")" ]]; then
        find "$(dirname "$WORKSPACE_BASE_DIR")" -maxdepth 1 -type d -name "e2e-validation-*" -mtime +"$days" -exec rm -rf {} \; 2>/dev/null || true
    fi
    
    log_info "Old workspaces cleaned"
}

# Export functions
export -f create_workspace
export -f init_terraform_workspace
export -f copy_config_to_workspace
export -f copy_provider_to_workspace
export -f get_state_file_path
export -f get_plan_file_path
export -f has_state
export -f get_state_size
export -f clean_workspace
export -f list_workspaces
export -f clean_old_workspaces
