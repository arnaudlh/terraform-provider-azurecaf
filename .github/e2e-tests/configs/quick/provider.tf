# Provider Configuration for E2E Tests
# Configures local provider for testing

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    azurecaf = {
      source  = "local/aztfmod/azurecaf"
      version = "9.9.9"
    }
  }
}

provider "azurecaf" {
  # No provider configuration needed for naming provider
  # All configuration is done at resource/data source level
}
