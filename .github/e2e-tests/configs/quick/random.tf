# Random Naming Convention Test Configuration
# Tests provider with random convention (no CAF structure)

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    azurecaf = {
      source  = "local/aztfmod/azurecaf"
      version = "9.9.9"
    }
  }
}

# Test representative resource types with random naming

# Storage Resources (random naming most common here)
resource "azurecaf_name" "storage_account" {
  resource_type = "azurerm_storage_account"
  random_length = 12
  random_seed   = 22222
}

resource "azurecaf_name" "key_vault" {
  resource_type = "azurerm_key_vault"
  random_length = 15
  random_seed   = 22223
}

resource "azurecaf_name" "container_registry" {
  resource_type = "azurerm_container_registry"
  random_length = 10
  random_seed   = 22224
}

# Database Resources
resource "azurecaf_name" "sql_server" {
  resource_type = "azurerm_sql_server"
  random_length = 15
  random_seed   = 22225
}

resource "azurecaf_name" "postgresql" {
  resource_type = "azurerm_postgresql_server"
  random_length = 15
  random_seed   = 22226
}

resource "azurecaf_name" "mysql" {
  resource_type = "azurerm_mysql_server"
  random_length = 15
  random_seed   = 22227
}

resource "azurecaf_name" "cosmosdb" {
  resource_type = "azurerm_cosmosdb_account"
  random_length = 15
  random_seed   = 22228
}

# Compute Resources
resource "azurecaf_name" "vm" {
  resource_type = "azurerm_virtual_machine"
  random_length = 10
  random_seed   = 22229
}

resource "azurecaf_name" "aks" {
  resource_type = "azurerm_kubernetes_cluster"
  random_length = 10
  random_seed   = 22230
}

# Application Services
resource "azurecaf_name" "app_service" {
  resource_type = "azurerm_app_service"
  random_length = 12
  random_seed   = 22231
}

resource "azurecaf_name" "function_app" {
  resource_type = "azurerm_function_app"
  random_length = 12
  random_seed   = 22232
}

# Networking Resources
resource "azurecaf_name" "virtual_network" {
  resource_type = "azurerm_virtual_network"
  random_length = 15
  random_seed   = 22233
}

resource "azurecaf_name" "public_ip" {
  resource_type = "azurerm_public_ip"
  random_length = 15
  random_seed   = 22234
}

# Management Resources
resource "azurecaf_name" "resource_group" {
  resource_type = "azurerm_resource_group"
  random_length = 15
  random_seed   = 22235
}

# Outputs to validate generated names
output "storage_account_name" {
  value = azurecaf_name.storage_account.result
}

output "key_vault_name" {
  value = azurecaf_name.key_vault.result
}

output "all_names" {
  value = {
    storage_account    = azurecaf_name.storage_account.result
    key_vault          = azurecaf_name.key_vault.result
    container_registry = azurecaf_name.container_registry.result
    sql_server         = azurecaf_name.sql_server.result
    postgresql         = azurecaf_name.postgresql.result
    mysql              = azurecaf_name.mysql.result
    cosmosdb           = azurecaf_name.cosmosdb.result
    vm                 = azurecaf_name.vm.result
    aks                = azurecaf_name.aks.result
    app_service        = azurecaf_name.app_service.result
    function_app       = azurecaf_name.function_app.result
    virtual_network    = azurecaf_name.virtual_network.result
    public_ip          = azurecaf_name.public_ip.result
    resource_group     = azurecaf_name.resource_group.result
  }
}
