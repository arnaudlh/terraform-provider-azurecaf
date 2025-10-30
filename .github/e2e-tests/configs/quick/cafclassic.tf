# CAF Classic Naming Convention Test Configuration
# Tests provider with cafclassic convention for representative Azure resource types

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    azurecaf = {
      source  = "local/aztfmod/azurecaf"
      version = "9.9.9"
    }
  }
}

# Test representative resource types across categories

# Compute Resources
data "azurecaf_name" "vm" {
  name          = "testapp"
  resource_type = "azurerm_virtual_machine"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "vmss" {
  name          = "testapp"
  resource_type = "azurerm_virtual_machine_scale_set"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "aks" {
  name          = "testapp"
  resource_type = "azurerm_kubernetes_cluster"
  prefixes      = ["dev"]
  random_seed   = 12345
}

# Storage Resources
data "azurecaf_name" "storage_account" {
  name          = "testapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "key_vault" {
  name          = "testapp"
  resource_type = "azurerm_key_vault"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

# Networking Resources
data "azurecaf_name" "virtual_network" {
  name          = "testapp"
  resource_type = "azurerm_virtual_network"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "subnet" {
  name          = "testapp"
  resource_type = "azurerm_subnet"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "nsg" {
  name          = "testapp"
  resource_type = "azurerm_network_security_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "public_ip" {
  name          = "testapp"
  resource_type = "azurerm_public_ip"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

# Database Resources
data "azurecaf_name" "sql_server" {
  name          = "testapp"
  resource_type = "azurerm_sql_server"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "postgresql" {
  name          = "testapp"
  resource_type = "azurerm_postgresql_server"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "mysql" {
  name          = "testapp"
  resource_type = "azurerm_mysql_server"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "cosmosdb" {
  name          = "testapp"
  resource_type = "azurerm_cosmosdb_account"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

# Application Services
data "azurecaf_name" "container_registry" {
  name          = "testapp"
  resource_type = "azurerm_container_registry"
  prefixes      = ["dev"]
  random_seed   = 12345
}

data "azurecaf_name" "app_service" {
  name          = "testapp"
  resource_type = "azurerm_app_service"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "function_app" {
  name          = "testapp"
  resource_type = "azurerm_function_app"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "application_gateway" {
  name          = "testapp"
  resource_type = "azurerm_application_gateway"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "api_management" {
  name          = "testapp"
  resource_type = "azurerm_api_management"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

# Management Resources
data "azurecaf_name" "resource_group" {
  name          = "testapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_seed   = 12345
}

data "azurecaf_name" "storage_blob" {
  name          = "testapp"
  resource_type = "azurerm_storage_blob"
  prefixes      = ["dev"]
  random_seed   = 12345
}

# Outputs to validate generated names
output "vm_name" {
  value = data.azurecaf_name.vm.result
}

output "storage_account_name" {
  value = data.azurecaf_name.storage_account.result
}

output "key_vault_name" {
  value = data.azurecaf_name.key_vault.result
}

output "all_names" {
  value = {
    vm                    = data.azurecaf_name.vm.result
    vmss                  = data.azurecaf_name.vmss.result
    aks                   = data.azurecaf_name.aks.result
    storage_account       = data.azurecaf_name.storage_account.result
    key_vault             = data.azurecaf_name.key_vault.result
    virtual_network       = data.azurecaf_name.virtual_network.result
    subnet                = data.azurecaf_name.subnet.result
    nsg                   = data.azurecaf_name.nsg.result
    public_ip             = data.azurecaf_name.public_ip.result
    sql_server            = data.azurecaf_name.sql_server.result
    postgresql            = data.azurecaf_name.postgresql.result
    mysql                 = data.azurecaf_name.mysql.result
    cosmosdb              = data.azurecaf_name.cosmosdb.result
    container_registry    = data.azurecaf_name.container_registry.result
    app_service           = data.azurecaf_name.app_service.result
    function_app          = data.azurecaf_name.function_app.result
    application_gateway   = data.azurecaf_name.application_gateway.result
    api_management        = data.azurecaf_name.api_management.result
    resource_group        = data.azurecaf_name.resource_group.result
    storage_blob          = data.azurecaf_name.storage_blob.result
  }
}
