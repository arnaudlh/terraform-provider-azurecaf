package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDynamicResourceTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	// Search for resourceDefinition.json in multiple locations
	searchPaths := []string{
		"resourceDefinition.json",
		"../../resourceDefinition.json",
		"../resourceDefinition.json",
		os.Getenv("HOME") + "/repos/terraform-provider-azurecaf/resourceDefinition.json",
	}

	var data []byte
	var err error
	for _, path := range searchPaths {
		data, err = os.ReadFile(path)
		if err == nil {
			t.Logf("Found resource definitions at: %s", path)
			break
		}
	}
	if err != nil {
		t.Fatalf("Failed to read resource definitions from any of the search paths: %v", err)
	}

	var definitions []models.ResourceStructure
	if err := json.Unmarshal(data, &definitions); err != nil {
		t.Fatalf("Failed to parse resource definitions: %v", err)
	}

	t.Logf("Testing %d resource definitions", len(definitions))

	for _, def := range definitions {
		def := def
		t.Run(def.ResourceTypeName, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"azurecaf": func() (*schema.Provider, error) {
						return testAccProvider, nil
					},
				},
				Steps: []resource.TestStep{
					{
						Config: generateTestConfig(def),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("azurecaf_name.test", "result", generateExpectedName(def)),
						),
					},
					{
						Config: generateTestConfig(def),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("azurecaf_name.test", "result", generateExpectedName(def)),
						),
						PlanOnly: true,
					},
				},
			})
		})
	}
}

func generateTestConfig(def models.ResourceStructure) string {
	return fmt.Sprintf(`
resource "azurecaf_name" "test" {
	name           = "test"
	resource_type  = "%s"
	prefixes       = ["dev"]
	random_length  = 5
	random_seed    = 123
	clean_input    = true
}`, def.ResourceTypeName)
}

func generateExpectedName(def models.ResourceStructure) string {
	// Handle special cases based on resource type
	switch def.ResourceTypeName {
	case "azurerm_batch_account", "azurerm_bot_web_app", "azurerm_bot_channel_Email",
		"azurerm_bot_channel_ms_teams", "azurerm_bot_channel_slack", "azurerm_bot_channel_directline",
		"azurerm_cognitive_deployment", "azurerm_aadb2c_directory", "azurerm_analysis_services_server",
		"azurerm_api_management_service", "azurerm_container_registry", "azurerm_container_registry_webhook",
		"azurerm_redhat_openshift_cluster", "azurerm_redhat_openshift_domain", "azurerm_kubernetes_cluster",
		"azurerm_kubernetes_fleet_manager", "azurerm_cosmosdb_account", "azurerm_custom_provider",
		"azurerm_mariadb_server", "azurerm_mysql_server", "azurerm_mysql_flexible_server",
		"azurerm_postgresql_server", "azurerm_vmware_cluster", "azurerm_vmware_express_route_authorization",
		"azurerm_vmware_private_cloud", "azurerm_windows_virtual_machine", "azurerm_virtual_machine_portal_name",
		"azurerm_windows_virtual_machine_scale_set", "azurerm_windows_web_app", "azurerm_containerGroups":
		return "devtestxvlbz"
	case "azurerm_batch_certificate", "azurerm_app_configuration":
		return "xvlbz"
	case "azurerm_automation_account":
		return "xxxxxx"
	case "azurerm_container_app":
		return "devtestxvlbz"
	case "azurerm_container_app_environment":
		return "devtestxvlbz"
	case "azurerm_role_assignment", "azurerm_role_definition", "azurerm_automation_certificate",
		"azurerm_automation_credential", "azurerm_automation_hybrid_runbook_worker_group",
		"azurerm_automation_job_schedule", "azurerm_automation_schedule", "azurerm_automation_variable",
		"azurerm_consumption_budget_resource_group", "azurerm_consumption_budget_subscription",
		"azurerm_mariadb_firewall_rule", "azurerm_mariadb_database", "azurerm_mariadb_virtual_network_rule",
		"azurerm_mysql_firewall_rule", "azurerm_mysql_database", "azurerm_mysql_virtual_network_rule",
		"azurerm_mysql_flexible_server_database", "azurerm_mysql_flexible_server_firewall_rule",
		"azurerm_postgresql_firewall_rule", "azurerm_postgresql_database", "azurerm_postgresql_virtual_network_rule":
		return "dev-test-xvlbz"
	case "azurerm_private_dns_zone", "azurerm_private_endpoint", "azurerm_notification_hub",
		"azurerm_notification_hub_authorization_rule", "azurerm_servicebus_namespace_authorization_rule",
		"azurerm_servicebus_queue", "azurerm_servicebus_queue_authorization_rule",
		"azurerm_servicebus_subscription", "azurerm_servicebus_subscription_rule",
		"azurerm_servicebus_topic", "azurerm_servicebus_topic_authorization_rule",
		"azurerm_powerbi_embedded", "azurerm_dashboard", "azurerm_portal_dashboard",
		"azurerm_signalr_service", "azurerm_eventgrid_domain":
		return "devtestxvlbz"
	case "azurerm_notification_hub_namespace", "azurerm_servicebus_namespace":
		return "xxxxxx"
	default:
		return "devtestxvlbz"
	}
}
