# Azure CAF Naming - Terraform provider

[![Build Status](https://github.com/aztfmod/terraform-provider-azurecaf/workflows/Go/badge.svg)](https://github.com/aztfmod/terraform-provider-azurecaf/actions)

The Azure CAF Naming provider is a naming convention repository that helps you naming your resources according to the Cloud Adoption Framework for Azure recommendations: https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging

The provider allows you to:
* Clean inputs to make sure they are compliant with the Azure resources naming convention
* Generate CAF compliant names for a variety of Azure resources
* Ensure generated names are compliant with the resource type they apply to
* Ensure naming convention compliance with the following standards:
  * Cloud Adoption Framework for Azure recommendations
  * Azure resources naming restrictions

## Example usage

```hcl
# Configure Terraform
terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "~> 1.2.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurecaf" {
}

# Generate a CAF compliant name
resource "azurecaf_name" "example" {
  name          = "example"
  resource_type = "azurerm_resource_group"
  prefixes      = ["a", "b"]
  suffixes      = ["y", "z"]
  random_length = 5
  clean_input   = true
}

output "random_result" {
  value = azurecaf_name.example.result
}
```

## Documentation

The documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/aztfmod/azurecaf/latest/docs).

## Supported Azure resources

The provider currently supports the following entities:

| Entity | Resource Type | slug | Min Length | Max Length | Lowercase | Regexp |
|---|---|---|---|---|---|---|
| azure automation account| azurerm_automation_account| aaa| 6| 50| false| ^[a-zA-Z][0-9A-Za-z-]{5,49}$|
| azure container app| azurerm_container_app| ac| 1| 32| true| ^[a-z0-9][a-z0-9-]{0,30}[a-z0-9]$|
| azure container app environment| azurerm_container_app_environment| ace| 1| 60| false| ^[0-9A-Za-z][0-9A-Za-z-]{0,58}[0-9a-zA-Z]$|
| azure container registry| azurerm_container_registry| acr| 5| 50| true| ^[0-9A-Za-z]{5,50}$|
| azure firewall| azurerm_firewall| afw| 1| 80| false| ^[a-zA-Z][0-9A-Za-z_.-]{0,79}$|
| application gateway| azurerm_application_gateway| agw| 1| 80| false| ^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$|
| azure kubernetes service| azurerm_kubernetes_cluster| aks| 1| 63| false| ^[0-9a-zA-Z][0-9A-Za-z_.-]{0,61}[0-9a-zA-Z]$|
| aksdns prefix| aks_dns_prefix| aksdns| 3| 45| false| ^[a-zA-Z][0-9A-Za-z-]{0,43}[0-9a-zA-Z]$|
| aks node pool for Linux| aks_node_pool_linux| aksnpl| 2| 12| true| ^[a-zA-Z][0-9a-z]{0,11}$|
| aks node pool for Windows| aks_node_pool_windows| aksnpw| 2| 6| true| ^[a-zA-Z][0-9a-z]{0,5}$|
| api management| azurerm_api_management| apim| 1| 50| false| ^[a-zA-Z][0-9A-Za-z]{0,49}$|
| web app| azurerm_app_service| app| 2| 60| false| ^[0-9A-Za-z][0-9A-Za-z-]{0,58}[0-9a-zA-Z]$|
| application insights| azurerm_application_insights| appi| 1| 260| false| ^[^%&\\?/. ][^%&\\?/]{0,258}[^%&\\?/. ]$|
| app service environment| azurerm_app_service_environment| ase| 2| 36| false| ^[0-9A-Za-z-]{2,36}$|
| azure site recovery| azurerm_recovery_services_vault| asr| 2| 50| false| ^[a-zA-Z][0-9A-Za-z-]{1,49}$|
| event hub| azurerm_eventhub_namespace| evh| 1| 50| false| ^[a-zA-Z][0-9A-Za-z-]{0,48}[0-9a-zA-Z]$|
| generic| generic| gen| 1| 24| false| ^[0-9a-zA-Z]{1,24}$|
| keyvault| azurerm_key_vault| kv| 3| 24| true| ^[a-zA-Z][0-9A-Za-z-]{0,22}[0-9a-zA-Z]$|
| loganalytics| azurerm_log_analytics_workspace| la| 4| 63| false| ^[0-9a-zA-Z][0-9A-Za-z-]{3,61}[0-9a-zA-Z]$|
| network interface card| azurerm_network_interface| nic| 1| 80| false| ^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$|
| network security group| azurerm_network_security_group| nsg| 1| 80| false| ^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$|
| public ip address| azurerm_public_ip| pip| 1| 80| false| ^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$|
| app service plan| azurerm_app_service_plan| plan| 1| 40| false| ^[0-9A-Za-z-]{1,40}$|
| resource group| azurerm_resource_group| rg| 1| 80| false| ^[-\w\._\(\)]{1,80}$|
| virtual network subnet| azurerm_subnet| snet| 1| 80| false| ^[0-9a-zA-Z][0-9A-Za-z_.-]{0,78}[0-9a-zA-Z_]$|
| azure sql db server| azurerm_sql_server| sql| 1| 63| true| ^[0-9a-z][0-9a-z-]{0,61}[0-9a-z]$|
| azure sql db| azurerm_sql_database| sqldb| 1| 128| false| ^[^<>*%&:\\/?. ][^<>*%&:\\/?]{0,126}[^<>*%&:\\/?. ]$|
| storage account| azurerm_storage_account| st| 3| 24| true| ^[0-9a-z]{3,24}$|
| virtual machine (linux)| azurerm_windows_virtual_machine_linux| vml| 1| 64| false| ^[0-9a-zA-Z][0-9A-Za-z_-]{0,62}[0-9a-zA-Z_]$|
| virtual machine (windows)| azurerm_windows_virtual_machine_windows| vmw| 1| 15| false| ^[0-9a-zA-Z][0-9A-Za-z_-]{0,13}[0-9a-zA-Z_]$|
| virtual network| azurerm_virtual_network| vnet| 2| 64| false| ^[0-9a-zA-Z][0-9A-Za-z_.-]{0,62}[0-9a-zA-Z_]$|

## Testing

Running the acceptance test suite requires does not require an Azure subscription.

to run the unit test:

```
make unittest
```

to run the integration test

```
make test
```

to run the end-to-end tests

```
make e2etest
```

to run all tests

```
make test-all
```

## Contributing

Contributions are welcome! Please see the [contributing guide](CONTRIBUTING.md) for more information.

## License

[MIT](LICENSE)
