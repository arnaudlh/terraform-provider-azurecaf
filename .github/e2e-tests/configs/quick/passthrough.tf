# Passthrough Naming Convention Test Configuration
# Tests provider with passthrough convention (validation only)

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    azurecaf = {
      source  = "local/aztfmod/azurecaf"
      version = "9.9.9"
    }
  }
}

# Test passthrough mode - provider validates but doesn't modify names

# Storage Resources
data "azurecaf_name" "storage_account" {
  name          = "devtestapp001"
  resource_type = "azurerm_storage_account"
  passthrough   = true
}

data "azurecaf_name" "key_vault" {
  name          = "kv-dev-testapp-001"
  resource_type = "azurerm_key_vault"
  passthrough   = true
}

data "azurecaf_name" "container_registry" {
  name          = "devtestappacr"
  resource_type = "azurerm_container_registry"
  passthrough   = true
}

# Database Resources
data "azurecaf_name" "sql_server" {
  name          = "sql-dev-testapp-001"
  resource_type = "azurerm_sql_server"
  passthrough   = true
}

data "azurecaf_name" "postgresql" {
  name          = "psql-dev-testapp-001"
  resource_type = "azurerm_postgresql_server"
  passthrough   = true
}

data "azurecaf_name" "mysql" {
  name          = "mysql-dev-testapp-001"
  resource_type = "azurerm_mysql_server"
  passthrough   = true
}

data "azurecaf_name" "cosmosdb" {
  name          = "cosmos-dev-testapp-001"
  resource_type = "azurerm_cosmosdb_account"
  passthrough   = true
}

# Compute Resources
data "azurecaf_name" "vm" {
  name          = "vm-dev-testapp-001"
  resource_type = "azurerm_virtual_machine"
  passthrough   = true
}

data "azurecaf_name" "vmss" {
  name          = "vmss-dev-testapp-001"
  resource_type = "azurerm_virtual_machine_scale_set"
  passthrough   = true
}

data "azurecaf_name" "aks" {
  name          = "aks-dev-testapp"
  resource_type = "azurerm_kubernetes_cluster"
  passthrough   = true
}

# Application Services
data "azurecaf_name" "app_service" {
  name          = "app-dev-testapp-001"
  resource_type = "azurerm_app_service"
  passthrough   = true
}

data "azurecaf_name" "function_app" {
  name          = "func-dev-testapp-001"
  resource_type = "azurerm_function_app"
  passthrough   = true
}

data "azurecaf_name" "api_management" {
  name          = "apim-dev-testapp-001"
  resource_type = "azurerm_api_management"
  passthrough   = true
}

# Networking Resources
data "azurecaf_name" "virtual_network" {
  name          = "vnet-dev-testapp-001"
  resource_type = "azurerm_virtual_network"
  passthrough   = true
}

data "azurecaf_name" "subnet" {
  name          = "snet-dev-testapp-001"
  resource_type = "azurerm_subnet"
  passthrough   = true
}

data "azurecaf_name" "nsg" {
  name          = "nsg-dev-testapp-001"
  resource_type = "azurerm_network_security_group"
  passthrough   = true
}

data "azurecaf_name" "public_ip" {
  name          = "pip-dev-testapp-001"
  resource_type = "azurerm_public_ip"
  passthrough   = true
}

# Management Resources
data "azurecaf_name" "resource_group" {
  name          = "rg-dev-testapp-001"
  resource_type = "azurerm_resource_group"
  passthrough   = true
}

# Outputs to validate names pass through unchanged
output "storage_account_name" {
  value = data.azurecaf_name.storage_account.result
}

output "key_vault_name" {
  value = data.azurecaf_name.key_vault.result
}

output "all_names" {
  value = {
    storage_account    = data.azurecaf_name.storage_account.result
    key_vault          = data.azurecaf_name.key_vault.result
    container_registry = data.azurecaf_name.container_registry.result
    sql_server         = data.azurecaf_name.sql_server.result
    postgresql         = data.azurecaf_name.postgresql.result
    mysql              = data.azurecaf_name.mysql.result
    cosmosdb           = data.azurecaf_name.cosmosdb.result
    vm                 = data.azurecaf_name.vm.result
    vmss               = data.azurecaf_name.vmss.result
    aks                = data.azurecaf_name.aks.result
    app_service        = data.azurecaf_name.app_service.result
    function_app       = data.azurecaf_name.function_app.result
    api_management     = data.azurecaf_name.api_management.result
    virtual_network    = data.azurecaf_name.virtual_network.result
    subnet             = data.azurecaf_name.subnet.result
    nsg                = data.azurecaf_name.nsg.result
    public_ip          = data.azurecaf_name.public_ip.result
    resource_group     = data.azurecaf_name.resource_group.result
  }
}
