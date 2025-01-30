terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "2.0.0-preview5"
    }
  }
}

provider "azurecaf" {}

resource "azurecaf_name" "test" {
  name          = "test"
  resource_type = "azurerm_resource_group"
  prefixes      = ["pr"]
  suffixes      = ["sf"]
  random_length = 5
  random_seed   = 123
  clean_input   = true
  use_slug      = true
}

output "result" {
  value = azurecaf_name.test.result
}
