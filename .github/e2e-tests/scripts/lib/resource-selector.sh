#!/usr/bin/env bash
# Test Resource Selector for E2E Validation
# Selects and filters test resources based on criteria

set -euo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Get resource types from resourceDefinition.json
# Usage: get_all_resource_types
get_all_resource_types() {
    local resource_def_file="${SCRIPT_DIR}/../../../../resourceDefinition.json"
    
    if [[ ! -f "$resource_def_file" ]]; then
        echo "ERROR: resourceDefinition.json not found" >&2
        return 1
    fi
    
    # Extract resource types from JSON
    jq -r 'keys[]' "$resource_def_file"
}

# Get representative resource types for quick validation (20 resources)
# Usage: get_quick_resource_types
get_quick_resource_types() {
    cat <<EOF
azurerm_resource_group
azurerm_storage_account
azurerm_key_vault
azurerm_virtual_machine
azurerm_virtual_machine_scale_set
azurerm_kubernetes_cluster
azurerm_virtual_network
azurerm_subnet
azurerm_network_security_group
azurerm_public_ip
azurerm_sql_server
azurerm_postgresql_server
azurerm_mysql_server
azurerm_cosmosdb_account
azurerm_container_registry
azurerm_app_service
azurerm_function_app
azurerm_application_gateway
azurerm_api_management
azurerm_storage_blob
EOF
}

# Get standard resource types for medium validation (100 resources)
# Usage: get_standard_resource_types
get_standard_resource_types() {
    # Return all compute, storage, networking, database resources
    get_all_resource_types | grep -E "^azurerm_(virtual_machine|storage_|key_vault|network_|subnet|sql_|postgresql_|mysql_|cosmosdb_|kubernetes_)"
}

# Filter resources by category
# Usage: filter_by_category <category>
# Categories: compute, storage, networking, database, analytics, security, management
filter_by_category() {
    local category="$1"
    
    case "$category" in
        compute)
            grep -E "^azurerm_(virtual_machine|vmss|kubernetes_cluster|container_|batch_)"
            ;;
        storage)
            grep -E "^azurerm_(storage_|key_vault)"
            ;;
        networking)
            grep -E "^azurerm_(virtual_network|subnet|network_|public_ip|application_gateway|load_balancer|firewall)"
            ;;
        database)
            grep -E "^azurerm_(sql_|postgresql_|mysql_|mariadb_|cosmosdb_|redis_)"
            ;;
        analytics)
            grep -E "^azurerm_(synapse_|data_factory|databricks_|hdinsight_|stream_analytics)"
            ;;
        security)
            grep -E "^azurerm_(key_vault|security_center_|sentinel_)"
            ;;
        management)
            grep -E "^azurerm_(resource_group|management_group|subscription)"
            ;;
        *)
            echo "ERROR: Unknown category: $category" >&2
            return 1
            ;;
    esac
}

# Filter resources by name pattern
# Usage: filter_by_pattern <pattern>
filter_by_pattern() {
    local pattern="$1"
    grep -E "$pattern"
}

# Select random subset of resources
# Usage: select_random_subset <count>
select_random_subset() {
    local count="$1"
    shuf | head -n "$count"
}

# Validate resource type exists in resourceDefinition.json
# Usage: validate_resource_type <resource_type>
validate_resource_type() {
    local resource_type="$1"
    local resource_def_file="${SCRIPT_DIR}/../../../../resourceDefinition.json"
    
    if ! jq -e --arg rt "$resource_type" '.[$rt]' "$resource_def_file" >/dev/null 2>&1; then
        echo "ERROR: Resource type not found in resourceDefinition.json: $resource_type" >&2
        return 1
    fi
    
    return 0
}

# Get resource metadata from resourceDefinition.json
# Usage: get_resource_metadata <resource_type>
get_resource_metadata() {
    local resource_type="$1"
    local resource_def_file="${SCRIPT_DIR}/../../../../resourceDefinition.json"
    
    jq --arg rt "$resource_type" '.[$rt]' "$resource_def_file"
}

# Export functions
export -f get_all_resource_types
export -f get_quick_resource_types
export -f get_standard_resource_types
export -f filter_by_category
export -f filter_by_pattern
export -f select_random_subset
export -f validate_resource_type
export -f get_resource_metadata
