terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "2.0.0-preview4"
    }
  }
}

provider "azurecaf" {}

# Test common Azure resource types
locals {
  test_cases = {
    rg = {
      name = "myapp"
      type = "azurerm_resource_group"
      prefixes = ["dev"]
      suffixes = ["001"]
    }
    st = {
      name = "data"
      type = "azurerm_storage_account"
      prefixes = ["dev"]
      random_length = 5
      use_slug = true
      separator = ""
    }
    kv = {
      name = "secrets"
      type = "azurerm_key_vault"
      prefixes = ["dev"]
      random_length = 5
    }
    func = {
      name = "api"
      type = "azurerm_function_app"
      prefixes = ["dev"]
      suffixes = ["func"]
    }
     acr = {
      name = "registry"
      type = "azurerm_container_registry"
      prefixes = ["dev"]
      random_length = 4
      separator = ""  # ACR names can't contain hyphens
      use_slug = true
    }
    aks = {
      name = "cluster"
      type = "azurerm_kubernetes_cluster"
      prefixes = ["dev"]
      suffixes = ["aks"]
    }
    app = {
      name = "webapp"
      type = "azurerm_app_service"
      prefixes = ["dev"]
      random_length = 3
    }
    cosmos = {
      name = "db"
      type = "azurerm_cosmosdb_account"
      prefixes = ["dev"]
      random_length = 4
    }
    sql = {
      name = "sqldb"
      type = "azurerm_sql_server"
      prefixes = ["dev"]
      random_length = 4
    }
    vm = {
      name = "server"
      type = "azurerm_virtual_machine"
      prefixes = ["dev"]
      suffixes = ["vm"]
    }
  }
}

# Test resources
resource "azurecaf_name" "test" {
  for_each = local.test_cases
  
  name = each.value.name
  resource_type = each.value.type
  prefixes = try(each.value.prefixes, [])
  suffixes = try(each.value.suffixes, [])
  random_length = try(each.value.random_length, 0)
  random_seed = 12345  # Fixed seed for consistent results
  clean_input = true
  separator = try(each.value.separator, "-")
  use_slug = try(each.value.use_slug, false)
}

# Test data sources
data "azurecaf_name" "test" {
  for_each = local.test_cases
  
  name = each.value.name
  resource_type = each.value.type
  prefixes = try(each.value.prefixes, [])
  suffixes = try(each.value.suffixes, [])
  random_length = try(each.value.random_length, 0)
  random_seed = 12345  # Same seed as resource for consistent results
  clean_input = true
  separator = try(each.value.separator, "-")
  use_slug = try(each.value.use_slug, false)
}

# Outputs for validation
output "resource_results" {
  value = {
    for k, v in azurecaf_name.test : k => v.result
  }
}

output "data_source_results" {
  value = {
    for k, v in data.azurecaf_name.test : k => v.result
  }
}

# Verify data source and resource results match
output "validation" {
  value = {
    for k, v in azurecaf_name.test : k => (
      v.result == data.azurecaf_name.test[k].result ? "PASS" : "FAIL: ${v.result} != ${data.azurecaf_name.test[k].result}"
    )
  }
}

# Additional validation outputs
output "resource_types" {
  value = distinct([for k, v in azurecaf_name.test : v.resource_type])
  description = "List of all resource types tested"
}

output "name_lengths" {
  value = {
    for k, v in azurecaf_name.test : k => length(v.result)
  }
  description = "Length of each generated name to verify against Azure limits"
}
