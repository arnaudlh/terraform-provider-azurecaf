terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "2.0.0-preview5"
    }
  }

  # Use local provider build
  provider_installation {
    filesystem_mirror {
      path = "/home/runner/.terraform.d/plugins"
    }
    direct {
      exclude = ["registry.terraform.io/aztfmod/azurecaf"]
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
  default_value = "default"
}
