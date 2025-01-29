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
