terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "2.0.0-preview5"
    }
  }
}

provider "azurecaf" {}

# Test resource name generation
resource "azurecaf_name" "test_st" {
  name          = "teststore"
  resource_type = "azurerm_storage_account"
  random_length = 5
  clean_input   = true
}

# Test data source
data "azurecaf_name" "test_st_data" {
  name          = "teststore"
  resource_type = "azurerm_storage_account"
  random_length = 5
  clean_input   = true
}

# Validate that resource and data source produce same output with same input
output "resource_name" {
  value = azurecaf_name.test_st.result
}

output "data_source_name" {
  value = data.azurecaf_name.test_st_data.result
}

# Test environment variable data source
data "azurecaf_environment_variable" "test_env" {
  name = "TEST_ENV_VAR"
  fails_if_empty = false
}

# Compute & Containers
resource "azurecaf_name" "test_aks" {
  name          = "testcluster"
  resource_type = "azurerm_kubernetes_cluster"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_aks_data" {
  name          = "testcluster"
  resource_type = "azurerm_kubernetes_cluster"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_acr" {
  name          = "testregistry"
  resource_type = "azurerm_container_registry"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_acr_data" {
  name          = "testregistry"
  resource_type = "azurerm_container_registry"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_vm" {
  name          = "testvm"
  resource_type = "azurerm_linux_virtual_machine"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_vm_data" {
  name          = "testvm"
  resource_type = "azurerm_linux_virtual_machine"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_function" {
  name          = "testfunction"
  resource_type = "azurerm_function_app"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_function_data" {
  name          = "testfunction"
  resource_type = "azurerm_function_app"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_app_service" {
  name          = "testapp"
  resource_type = "azurerm_app_service"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_app_service_data" {
  name          = "testapp"
  resource_type = "azurerm_app_service"
  random_length = 5
  clean_input   = true
}

# Storage & Data
resource "azurecaf_name" "test_cosmos" {
  name          = "testcosmos"
  resource_type = "azurerm_cosmosdb_account"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_cosmos_data" {
  name          = "testcosmos"
  resource_type = "azurerm_cosmosdb_account"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_redis" {
  name          = "testredis"
  resource_type = "azurerm_redis_cache"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_redis_data" {
  name          = "testredis"
  resource_type = "azurerm_redis_cache"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_servicebus" {
  name          = "testservicebus"
  resource_type = "azurerm_servicebus_namespace"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_servicebus_data" {
  name          = "testservicebus"
  resource_type = "azurerm_servicebus_namespace"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_sql" {
  name          = "testsql"
  resource_type = "azurerm_sql_server"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_sql_data" {
  name          = "testsql"
  resource_type = "azurerm_sql_server"
  random_length = 5
  clean_input   = true
}

# Networking
resource "azurecaf_name" "test_vnet" {
  name          = "testvnet"
  resource_type = "azurerm_virtual_network"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_vnet_data" {
  name          = "testvnet"
  resource_type = "azurerm_virtual_network"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_nsg" {
  name          = "testnsg"
  resource_type = "azurerm_network_security_group"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_nsg_data" {
  name          = "testnsg"
  resource_type = "azurerm_network_security_group"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_appgw" {
  name          = "testappgw"
  resource_type = "azurerm_application_gateway"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_appgw_data" {
  name          = "testappgw"
  resource_type = "azurerm_application_gateway"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_privatedns" {
  name          = "testprivatedns"
  resource_type = "azurerm_private_dns_zone"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_privatedns_data" {
  name          = "testprivatedns"
  resource_type = "azurerm_private_dns_zone"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_firewall" {
  name          = "testfw"
  resource_type = "azurerm_firewall"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_firewall_data" {
  name          = "testfw"
  resource_type = "azurerm_firewall"
  random_length = 5
  clean_input   = true
}

# Management & Security
resource "azurecaf_name" "test_rg" {
  name          = "testrg"
  resource_type = "azurerm_resource_group"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_rg_data" {
  name          = "testrg"
  resource_type = "azurerm_resource_group"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_kv" {
  name          = "testkv"
  resource_type = "azurerm_key_vault"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_kv_data" {
  name          = "testkv"
  resource_type = "azurerm_key_vault"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_log" {
  name          = "testlog"
  resource_type = "azurerm_log_analytics_workspace"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_log_data" {
  name          = "testlog"
  resource_type = "azurerm_log_analytics_workspace"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_monitor" {
  name          = "testmonitor"
  resource_type = "azurerm_monitor_action_group"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_monitor_data" {
  name          = "testmonitor"
  resource_type = "azurerm_monitor_action_group"
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_automation" {
  name          = "testauto"
  resource_type = "azurerm_automation_account"
  random_length = 5
  clean_input   = true
}

data "azurecaf_name" "test_automation_data" {
  name          = "testauto"
  resource_type = "azurerm_automation_account"
  random_length = 5
  clean_input   = true
}

# Outputs to verify resource and data source results match
output "compute_aks_resource" { value = azurecaf_name.test_aks.result }
output "compute_aks_data" { value = data.azurecaf_name.test_aks_data.result }
output "compute_acr_resource" { value = azurecaf_name.test_acr.result }
output "compute_acr_data" { value = data.azurecaf_name.test_acr_data.result }
output "compute_vm_resource" { value = azurecaf_name.test_vm.result }
output "compute_vm_data" { value = data.azurecaf_name.test_vm_data.result }
output "compute_function_resource" { value = azurecaf_name.test_function.result }
output "compute_function_data" { value = data.azurecaf_name.test_function_data.result }
output "compute_app_service_resource" { value = azurecaf_name.test_app_service.result }
output "compute_app_service_data" { value = data.azurecaf_name.test_app_service_data.result }

output "storage_cosmos_resource" { value = azurecaf_name.test_cosmos.result }
output "storage_cosmos_data" { value = data.azurecaf_name.test_cosmos_data.result }
output "storage_redis_resource" { value = azurecaf_name.test_redis.result }
output "storage_redis_data" { value = data.azurecaf_name.test_redis_data.result }
output "storage_servicebus_resource" { value = azurecaf_name.test_servicebus.result }
output "storage_servicebus_data" { value = data.azurecaf_name.test_servicebus_data.result }
output "storage_sql_resource" { value = azurecaf_name.test_sql.result }
output "storage_sql_data" { value = data.azurecaf_name.test_sql_data.result }

output "network_vnet_resource" { value = azurecaf_name.test_vnet.result }
output "network_vnet_data" { value = data.azurecaf_name.test_vnet_data.result }
output "network_nsg_resource" { value = azurecaf_name.test_nsg.result }
output "network_nsg_data" { value = data.azurecaf_name.test_nsg_data.result }
output "network_appgw_resource" { value = azurecaf_name.test_appgw.result }
output "network_appgw_data" { value = data.azurecaf_name.test_appgw_data.result }
output "network_privatedns_resource" { value = azurecaf_name.test_privatedns.result }
output "network_privatedns_data" { value = data.azurecaf_name.test_privatedns_data.result }
output "network_firewall_resource" { value = azurecaf_name.test_firewall.result }
output "network_firewall_data" { value = data.azurecaf_name.test_firewall_data.result }

output "mgmt_rg_resource" { value = azurecaf_name.test_rg.result }
output "mgmt_rg_data" { value = data.azurecaf_name.test_rg_data.result }
output "mgmt_kv_resource" { value = azurecaf_name.test_kv.result }
output "mgmt_kv_data" { value = data.azurecaf_name.test_kv_data.result }
output "mgmt_log_resource" { value = azurecaf_name.test_log.result }
output "mgmt_log_data" { value = data.azurecaf_name.test_log_data.result }
output "mgmt_monitor_resource" { value = azurecaf_name.test_monitor.result }
output "mgmt_monitor_data" { value = data.azurecaf_name.test_monitor_data.result }
output "mgmt_automation_resource" { value = azurecaf_name.test_automation.result }
output "mgmt_automation_data" { value = data.azurecaf_name.test_automation_data.result }
