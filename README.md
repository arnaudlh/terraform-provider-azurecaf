# Azure Terraform SRE - Terraform provider

> :warning: This solution, offered by the Open-Source community, will no longer receive contributions from Microsoft.

This provider implements a set of methodologies for naming convention implementation including the default Microsoft Cloud Adoption Framework for Azure recommendations as per <https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging>.

## Using the Provider

You can simply consume the provider from the Terraform registry from the following URL: [https://registry.terraform.io/providers/aztfmod/azurecaf/latest](https://registry.terraform.io/providers/aztfmod/azurecaf/latest), then add it in your provider declaration as follow:

```hcl
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "1.2.10"
    }
  }
}
```

The azurecaf_name resource allows you to:

* Clean inputs to make sure they remain compliant with the allowed patterns for each Azure resource.
* Generate random characters to append at the end of the resource name.
* Handle prefix, suffixes (either manual or as per the Azure cloud adoption framework resource conventions).
* Allow passthrough mode (simply validate the output).

## Example usage

This example outputs one name, the result of the naming convention query. The result attribute returns the name based on the convention and parameters input.

The example generates a 23 characters name compatible with the specification for an Azure Resource Group
dev-aztfmod-001

```hcl
data "azurecaf_name" "rg_example" {
  name          = "demogroup"
  resource_type = "azurerm_resource_group"
  prefixes      = ["a", "b"]
  suffixes      = ["y", "z"]
  random_length = 5
  clean_input   = true
}

output "rg_example" {
  value = data.azurecaf_name.rg_example.result
}
```

```
data.azurecaf_name.rg_example: Reading...
data.azurecaf_name.rg_example: Read complete after 0s [id=a-b-rg-demogroup-sjdeh-y-z]

Changes to Outputs:
  + rg_example = "a-b-rg-demogroup-sjdeh-y-z"
```

The provider generates a name using the input parameters and automatically appends a prefix (if defined), a caf prefix (resource type) and postfix (if defined) in addition to a generated padding string based on the selected naming convention.

The example above would generate a name using the pattern [prefix]-[cafprefix]-[name]-[postfix]-[5_random_chars]:

## Argument Reference

The following arguments are supported:

* **name** - (optional) the basename of the resource to create, the basename will be sanitized as per supported characters set for each Azure resources.
* **prefixes** (optional) - a list of prefix to append as the first characters of the generated name - prefixes will be separated by the separator character
* **suffixes** (optional) -  a list of additional suffix added after the basename, this is can be used to append resource index (eg. vm-001). Suffixes are separated by the separator character
* **random_length** (optional) - default to ``0`` : configure additional characters to append to the generated resource name. Random characters will remain compliant with the set of allowed characters per resources and will be appended before suffix(ess).
* **random_seed** (optional) - default to ``0`` : Define the seed to be used for random generator. 0 will not be respected and will generate a seed based in the unix time of the generation.
* **resource_type** (optional) -  describes the type of azure resource you are requesting a name from (eg. azure container registry: azurerm_container_registry). See the Resource Type section
* **resource_types** (optional) -  a list of additional resource type should you want to use the same settings for a set of resources
* **separator** (optional) - defaults to ``-``. The separator character to use between prefixes, resource type, name, suffixes, random character
* **clean_input** (optional) - defaults to ``true``. remove any noncompliant character from the name, suffix or prefix.
* **passthrough** (optional) - defaults to ``false``. Enables the passthrough mode - in that case only the clean input option is considered and the prefixes, suffixes, random, and are ignored. The resource prefixe is not added either to the resulting string
* **use_slug** (optional) - defaults to ``true``. If a slug should be added to the name - If you put false no slug (the few letters that identify the resource type) will be added to the name.

## Attributes Reference

The following attributes are exported:

* **id** - The id of the naming convention object
* **result** - The generated named for an Azure Resource based on the input parameter and the selected naming convention
* **results** - The generated name for the Azure resources based in the resource_types list

## Resource types

We define resource types as per [naming-and-tagging](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
The comprehensive list of resource type can be found [here](./docs/resources/azurecaf_name.md)

## Building the provider

Clone repository to: $GOPATH/src/github.com/aztfmod/terraform-provider-azurecaf

```
mkdir -p $GOPATH/src/github.com/aztfmod; cd $GOPATH/src/github.com/aztfmod
git clone https://github.com/aztfmod/terraform-provider-azurecaf.git

```

Enter the provider directory and build the provider

```
cd $GOPATH/src/github.com/aztfmod/terraform-provider-azurecaf
make build

```

## Developing the provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.13+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To display the makefile help run `make` or `make help`.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ make build
...
$ $GOPATH/bin/terraform-provider-azurecaf
...

```

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

to run all tests (unit, integration, and e2e)

```
make test-all
```

## Related repositories

| Repo                                                                                             | Description                                                |
|--------------------------------------------------------------------------------------------------|------------------------------------------------------------|
| [caf-terraform-landingzones](https://github.com/azure/caf-terraform-landingzones)                | landing zones repo with sample and core documentations     |
| [rover](https://github.com/aztfmod/rover)                                                        | devops toolset for operating landing zones                 |
| [azure_caf_provider](https://github.com/aztfmod/terraform-provider-azurecaf)                     | custom provider for naming conventions                     |
| [module](https://registry.terraform.io/modules/aztfmod)                                          | official CAF module available in the Terraform registry    |

## Community

Feel free to open an issue for feature or bug, or to submit a PR.

In case you have any question, you can reach out to tf-landingzones at microsoft dot com.

You can also reach us on [Gitter](https://gitter.im/aztfmod/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Contributing

information about contributing can be found at [CONTRIBUTING.md](.github/CONTRIBUTING.md)

## Resource Status

This is the current comprehensive status of the implemented resources in the provider comparing with the current list of resources in the azurerm terraform provider.

|resource | status |
|---|---|
|azurerm_aadb2c_directory | ✔ |
|azurerm_advanced_threat_protection | ❌ |
|azurerm_advisor_recommendations | ❌ |
|azurerm_analysis_services_server | ✔ |
|azurerm_api_management | ✔ |
|azurerm_api_management_api | ✔ |
|azurerm_api_management_api_diagnostic | ❌ |
|azurerm_api_management_api_operation | ❌ |
|azurerm_api_management_api_operation_policy | ❌ |
|azurerm_api_management_api_operation_tag | ✔ |
|azurerm_api_management_api_policy | ❌ |
|azurerm_api_management_api_schema | ❌ |
|azurerm_api_management_api_version_set | ❌ |
|azurerm_api_management_authorization_server | ❌ |
|azurerm_api_management_backend | ✔ |
|azurerm_api_management_certificate | ✔ |
|azurerm_api_management_custom_domain | ✔ |
|azurerm_api_management_diagnostic | ❌ |
|azurerm_api_management_gateway | ✔ |
|azurerm_api_management_group | ✔ |
|azurerm_api_management_group_user | ✔ |
|azurerm_api_management_identity_provider_aad | ❌ |
|azurerm_api_management_identity_provider_facebook | ❌ |
|azurerm_api_management_identity_provider_google | ❌ |
|azurerm_api_management_identity_provider_microsoft | ❌ |
|azurerm_api_management_identity_provider_twitter | ❌ |
|azurerm_api_management_logger | ✔ |
|azurerm_api_management_named_value | ❌ |
|azurerm_api_management_openid_connect_provider | ❌ |
|azurerm_api_management_product | ❌ |
|azurerm_api_management_product_api | ❌ |
|azurerm_api_management_product_group | ❌ |
|azurerm_api_management_product_policy | ❌ |
|azurerm_api_management_property | ❌ |
|azurerm_api_management_subscription | ❌ |
|azurerm_api_management_user | ✔ |
|azurerm_app_configuration | ✔ |
|azurerm_app_service | ✔ |
|azurerm_app_service_active_slot | ❌ |
|azurerm_app_service_certificate | ❌ |
|azurerm_app_service_certificate_order | ❌ |
|azurerm_app_service_custom_hostname_binding | ❌ |
|azurerm_app_service_environment | ✔ |
|azurerm_app_service_hybrid_connection | ❌ |
|azurerm_app_service_plan | ✔ |
|azurerm_app_service_slot | ❌ |
|azurerm_app_service_slot_virtual_network_swift_connection | ❌ |
|azurerm_app_service_source_control_token | ❌ |
|azurerm_app_service_virtual_network_swift_connection | ❌ |
|azurerm_application_gateway | ✔ |
|azurerm_application_insights | ✔ |
|azurerm_application_insights_analytics_item | ❌ |
|azurerm_application_insights_api_key | ❌ |
|azurerm_application_insights_web_test | ✔ |
|azurerm_application_security_group | ✔ |
|azurerm_attestation | ❌ |
|azurerm_automation_account | ✔ |
|azurerm_automation_certificate | ✔ |
|azurerm_automation_connection | ❌ |
|azurerm_automation_connection_certificate | ❌ |
|azurerm_automation_connection_classic_certificate | ❌ |
|azurerm_automation_connection_service_principal | ❌ |
|azurerm_automation_credential | ✔ |
|azurerm_automation_dsc_configuration | ❌ |
|azurerm_automation_dsc_nodeconfiguration | ❌ |
|azurerm_automation_hybrid_runbook_worker_group | ✔ |
|azurerm_automation_job_schedule | ✔ |
|azurerm_automation_module | ❌ |
|azurerm_automation_runbook | ✔ |
|azurerm_automation_schedule | ✔ |
|azurerm_automation_variable_bool | ❌ |
|azurerm_automation_variable_datetime | ❌ |
|azurerm_automation_variable_int | ❌ |
|azurerm_automation_variable_string | ❌ |
|azurerm_availability_set | ✔ |
|azurerm_backup_container_storage_account | ❌ |
|azurerm_backup_policy_file_share | ❌ |
|azurerm_backup_policy_vm | ❌ |
|azurerm_backup_protected_file_share | ❌ |
|azurerm_backup_protected_vm | ❌ |
|azurerm_bastion_host | ✔ |
|azurerm_batch_account | ✔ |
|azurerm_batch_application | ✔ |
|azurerm_batch_certificate | ✔ |
|azurerm_batch_pool | ✔ |
|azurerm_blueprint_assignment | ❌ |
|azurerm_blueprint_definition | ❌ |
|azurerm_blueprint_published_version | ❌ |
|azurerm_bot_channel_directline | ✔ |
|azurerm_bot_channel_email | ❌ |
|azurerm_bot_channel_ms_teams | ✔ |
|azurerm_bot_channel_slack | ✔ |
|azurerm_bot_channels_registration | ✔ |
|azurerm_bot_connection | ✔ |
|azurerm_bot_web_app | ✔ |
|azurerm_cdn_endpoint | ✔ |
|azurerm_cdn_frontdoor_custom_domain | ✔ |
|azurerm_cdn_frontdoor_endpoint | ✔ |
|azurerm_cdn_frontdoor_firewall_policy | ✔ |
|azurerm_cdn_frontdoor_origin | ✔ |
|azurerm_cdn_frontdoor_origin_group | ✔ |
|azurerm_cdn_frontdoor_profile | ✔ |
|azurerm_cdn_frontdoor_route | ✔ |
|azurerm_cdn_frontdoor_rule | ✔ |
|azurerm_cdn_frontdoor_rule_set | ✔ |
|azurerm_cdn_frontdoor_secret | ✔ |
|azurerm_cdn_frontdoor_security_policy | ✔ |
|azurerm_cdn_profile | ✔ |
|azurerm_client_config | ❌ |
|azurerm_cognitive_account | ✔ |
|azurerm_communication_service | ✔ |
|azurerm_consumption_budget_resource_group | ✔ |
|azurerm_consumption_budget_subscription | ✔ |
|azurerm_container_app | ✔ |
|azurerm_container_app_environment | ✔ |
|azurerm_container_group | ❌ |
|azurerm_container_registry | ✔ |
|azurerm_container_registry_webhook | ✔ |
|azurerm_cosmosdb_account | ✔ |
|azurerm_cosmosdb_cassandra_keyspace | ❌ |
|azurerm_cosmosdb_gremlin_database | ❌ |
|azurerm_cosmosdb_gremlin_graph | ❌ |
|azurerm_cosmosdb_mongo_collection | ❌ |
|azurerm_cosmosdb_mongo_database | ❌ |
|azurerm_cosmosdb_sql_container | ❌ |
|azurerm_cosmosdb_sql_database | ❌ |
|azurerm_cosmosdb_sql_stored_procedure | ❌ |
|azurerm_cosmosdb_table | ❌ |
|azurerm_cost_management_export_resource_group | ❌ |
|azurerm_custom_provider | ✔ |
|azurerm_dashboard | ✔ |
|azurerm_portal_dashboard | ✔ |
|azurerm_data_factory | ✔ |
|azurerm_data_factory_dataset_azure_blob | ✔ |
|azurerm_data_factory_dataset_cosmosdb_sqlapi | ✔ |
|azurerm_data_factory_dataset_delimited_text | ✔ |
|azurerm_data_factory_dataset_http | ✔ |
|azurerm_data_factory_dataset_json | ✔ |
|azurerm_data_factory_dataset_mysql | ✔ |
|azurerm_data_factory_dataset_postgresql | ✔ |
|azurerm_data_factory_dataset_sql_server_table | ✔ |
|azurerm_data_factory_integration_runtime_managed | ✔ |
|azurerm_data_factory_integration_runtime_self_hosted | ❌ |
|azurerm_data_factory_linked_service_azure_blob_storage | ✔ |
|azurerm_data_factory_linked_service_azure_databricks | ✔ |
|azurerm_data_factory_linked_service_azure_file_storage | ❌ |
|azurerm_data_factory_linked_service_azure_function | ✔ |
|azurerm_data_factory_linked_service_azure_sql_database | ✔ |
|azurerm_data_factory_linked_service_cosmosdb | ✔ |
|azurerm_data_factory_linked_service_data_lake_storage_gen2 | ✔ |
|azurerm_data_factory_linked_service_key_vault | ✔ |
|azurerm_data_factory_linked_service_mysql | ✔ |
|azurerm_data_factory_linked_service_postgresql | ✔ |
|azurerm_data_factory_linked_service_sftp | ✔ |
|azurerm_data_factory_linked_service_sql_server | ✔ |
|azurerm_data_factory_linked_service_web | ✔ |
|azurerm_data_factory_pipeline | ✔ |
|azurerm_data_factory_trigger_schedule | ✔ |
|azurerm_data_lake_analytics_account | ✔ |
|azurerm_data_lake_analytics_firewall_rule | ✔ |
|azurerm_data_lake_store | ✔ |
|azurerm_data_lake_store_file | ❌ |
|azurerm_data_lake_store_firewall_rule | ✔ |
|azurerm_data_protection_backup_policy_blob_storage | ✔ |
|azurerm_data_protection_backup_policy_disk | ✔ |
|azurerm_data_protection_backup_policy_postgresql | ✔ |
|azurerm_data_protection_backup_policy_postgresql_flexible_server | ✔ |
|azurerm_data_protection_backup_vault | ✔ |
|azurerm_data_share | ❌ |
|azurerm_data_share_account | ❌ |
|azurerm_data_share_dataset_blob_storage | ❌ |
|azurerm_data_share_dataset_data_lake_gen1 | ❌ |
|azurerm_data_share_dataset_data_lake_gen2 | ❌ |
|azurerm_data_share_dataset_kusto_cluster | ❌ |
|azurerm_data_share_dataset_kusto_database | ❌ |
|azurerm_database_migration_project | ✔ |
|azurerm_database_migration_service | ✔ |
|azurerm_databricks_workspace | ✔ |
|azurerm_dedicated_hardware_security_module | ❌ |
|azurerm_dedicated_host | ✔ |
|azurerm_dedicated_host_group | ✔ |
|azurerm_dev_test_global_vm_shutdown_schedule | ❌ |
|azurerm_dev_test_lab | ✔ |
|azurerm_dev_test_linux_virtual_machine | ✔ |
|azurerm_dev_test_policy | ❌ |
|azurerm_dev_test_schedule | ❌ |
|azurerm_dev_test_virtual_network | ❌ |
|azurerm_dev_test_windows_virtual_machine | ✔ |
|azurerm_devspace_controller | ❌ |
|azurerm_digital_twins_endpoint_eventgrid | ✔ |
|azurerm_digital_twins_endpoint_eventhub | ✔ |
|azurerm_digital_twins_endpoint_servicebus | ✔ |
|azurerm_digital_twins_instance | ✔ |
|azurerm_disk_encryption_set | ✔ |
|azurerm_dns_a_record | ❌ |
