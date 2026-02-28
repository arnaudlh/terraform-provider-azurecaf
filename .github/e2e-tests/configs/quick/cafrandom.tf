# CAF Random Naming Convention Test Configuration
# Tests provider with cafrandom convention for representative Azure resource types

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    azurecaf = {
      source  = "local/aztfmod/azurecaf"
      version = "9.9.9"
    }
  }
}

# Test representative resource types with random suffixes

# Compute Resources
resource "azurecaf_name" "vm" {
  name          = "testapp"
  resource_type = "azurerm_virtual_machine"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12345
}

resource "azurecaf_name" "vmss" {
  name          = "testapp"
  resource_type = "azurerm_virtual_machine_scale_set"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12346
}

resource "azurecaf_name" "aks" {
  name          = "testapp"
  resource_type = "azurerm_kubernetes_cluster"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12347
}

# Storage Resources
resource "azurecaf_name" "storage_account" {
  name          = "testapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12348
}

resource "azurecaf_name" "key_vault" {
  name          = "testapp"
  resource_type = "azurerm_key_vault"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12349
}

# Networking Resources
resource "azurecaf_name" "virtual_network" {
  name          = "testapp"
  resource_type = "azurerm_virtual_network"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12350
}

resource "azurecaf_name" "subnet" {
  name          = "testapp"
  resource_type = "azurerm_subnet"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12351
}

resource "azurecaf_name" "nsg" {
  name          = "testapp"
  resource_type = "azurerm_network_security_group"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12352
}

resource "azurecaf_name" "public_ip" {
  name          = "testapp"
  resource_type = "azurerm_public_ip"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12353
}

# Database Resources
resource "azurecaf_name" "sql_server" {
  name          = "testapp"
  resource_type = "azurerm_sql_server"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12354
}

resource "azurecaf_name" "postgresql" {
  name          = "testapp"
  resource_type = "azurerm_postgresql_server"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12355
}

resource "azurecaf_name" "mysql" {
  name          = "testapp"
  resource_type = "azurerm_mysql_server"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12356
}

resource "azurecaf_name" "cosmosdb" {
  name          = "testapp"
  resource_type = "azurerm_cosmosdb_account"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12357
}

# Application Services
resource "azurecaf_name" "container_registry" {
  name          = "testapp"
  resource_type = "azurerm_container_registry"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12358
}

resource "azurecaf_name" "app_service" {
  name          = "testapp"
  resource_type = "azurerm_app_service"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12359
}

resource "azurecaf_name" "function_app" {
  name          = "testapp"
  resource_type = "azurerm_function_app"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12360
}

resource "azurecaf_name" "application_gateway" {
  name          = "testapp"
  resource_type = "azurerm_application_gateway"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12361
}

resource "azurecaf_name" "api_management" {
  name          = "testapp"
  resource_type = "azurerm_api_management"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12362
}

# Management Resources
resource "azurecaf_name" "resource_group" {
  name          = "testapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  random_length = 3
  random_seed   = 12363
}

resource "azurecaf_name" "storage_blob" {
  name          = "testapp"
  resource_type = "azurerm_storage_blob"
  random_length = 8
  random_seed   = 12364
}

# Outputs to validate generated names
output "vm_name" {
  value = azurecaf_name.vm.result
}

output "storage_account_name" {
  value = azurecaf_name.storage_account.result
}

output "key_vault_name" {
  value = azurecaf_name.key_vault.result
}

output "all_names" {
  value = {
    vm                    = azurecaf_name.vm.result
    vmss                  = azurecaf_name.vmss.result
    aks                   = azurecaf_name.aks.result
    storage_account       = azurecaf_name.storage_account.result
    key_vault             = azurecaf_name.key_vault.result
    virtual_network       = azurecaf_name.virtual_network.result
    subnet                = azurecaf_name.subnet.result
    nsg                   = azurecaf_name.nsg.result
    public_ip             = azurecaf_name.public_ip.result
    sql_server            = azurecaf_name.sql_server.result
    postgresql            = azurecaf_name.postgresql.result
    mysql                 = azurecaf_name.mysql.result
    cosmosdb              = azurecaf_name.cosmosdb.result
    container_registry    = azurecaf_name.container_registry.result
    app_service           = azurecaf_name.app_service.result
    function_app          = azurecaf_name.function_app.result
    application_gateway   = azurecaf_name.application_gateway.result
    api_management        = azurecaf_name.api_management.result
    resource_group        = azurecaf_name.resource_group.result
    storage_blob          = azurecaf_name.storage_blob.result
  }
}
