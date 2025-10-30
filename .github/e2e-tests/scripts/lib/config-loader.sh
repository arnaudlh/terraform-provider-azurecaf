#!/usr/bin/env bash
# Configuration Loader for E2E Validation
# Loads and validates test configurations

set -euo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="${SCRIPT_DIR}/../../configs"
SCHEMA_FILE="${CONFIG_DIR}/schema.json"

# Load configuration file
# Usage: load_config <config_file>
load_config() {
    local config_file="$1"
    
    if [[ ! -f "$config_file" ]]; then
        echo "ERROR: Configuration file not found: $config_file" >&2
        return 1
    fi
    
    # Validate JSON syntax
    if ! jq empty "$config_file" 2>/dev/null; then
        echo "ERROR: Invalid JSON in configuration file: $config_file" >&2
        return 1
    fi
    
    # Export config to stdout
    cat "$config_file"
}

# Validate configuration against schema
# Usage: validate_config <config_json>
validate_config() {
    local config_json="$1"
    
    if [[ ! -f "$SCHEMA_FILE" ]]; then
        echo "WARNING: Schema file not found, skipping validation: $SCHEMA_FILE" >&2
        return 0
    fi
    
    # Basic validation checks
    local name
    name=$(echo "$config_json" | jq -r '.name // empty')
    if [[ -z "$name" ]]; then
        echo "ERROR: Configuration missing required field: name" >&2
        return 1
    fi
    
    local resource_types
    resource_types=$(echo "$config_json" | jq -r '.resource_types | length')
    if [[ "$resource_types" -eq 0 ]]; then
        echo "ERROR: Configuration must specify at least one resource_type" >&2
        return 1
    fi
    
    local conventions
    conventions=$(echo "$config_json" | jq -r '.naming_conventions | length')
    if [[ "$conventions" -eq 0 ]]; then
        echo "ERROR: Configuration must specify at least one naming_convention" >&2
        return 1
    fi
    
    local tf_versions
    tf_versions=$(echo "$config_json" | jq -r '.terraform_versions | length')
    if [[ "$tf_versions" -eq 0 ]]; then
        echo "ERROR: Configuration must specify at least one terraform_version" >&2
        return 1
    fi
    
    echo "Configuration validation passed" >&2
    return 0
}

# Get configuration value
# Usage: get_config_value <config_json> <json_path>
get_config_value() {
    local config_json="$1"
    local json_path="$2"
    
    echo "$config_json" | jq -r "$json_path"
}

# Get configuration value with default
# Usage: get_config_value_or_default <config_json> <json_path> <default>
get_config_value_or_default() {
    local config_json="$1"
    local json_path="$2"
    local default="$3"
    
    local value
    value=$(echo "$config_json" | jq -r "$json_path // empty")
    
    if [[ -z "$value" || "$value" == "null" ]]; then
        echo "$default"
    else
        echo "$value"
    fi
}

# Export functions
export -f load_config
export -f validate_config
export -f get_config_value
export -f get_config_value_or_default
