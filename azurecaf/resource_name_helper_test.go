package azurecaf

import (
	"context"
	"testing"
)

func TestHelperFunctions(t *testing.T) {
	resource, err := getResource("azurerm_resource_group")
	if err != nil {
		t.Fatalf("getResource returned error: %v", err)
	}

	result := cleanString("test!@#$%^&*()_+", resource)
	if result == "" {
		t.Fatal("Expected non-empty result from cleanString")
	}

	slice := []string{"test!@#", "test$%^"}
	cleanedSlice := cleanSlice(slice, resource)
	if len(cleanedSlice) != 2 {
		t.Fatalf("Expected 2 items in cleaned slice, got %d", len(cleanedSlice))
	}

	result = concatenateParameters("-", []string{"test1"}, []string{"test2"}, []string{"test3"})
	if result != "test1-test2-test3" {
		t.Fatalf("Expected 'test1-test2-test3', got '%s'", result)
	}

	resource, err = getResource("azurerm_resource_group")
	if err != nil {
		t.Fatalf("getResource returned error: %v", err)
	}
	if resource.ResourceTypeName == "" {
		t.Fatal("Expected non-empty ResourceTypeName")
	}

	slug := getSlug("azurerm_resource_group", "cafrandom")
	if slug == "" {
		t.Fatal("Expected non-empty slug")
	}

	result = trimResourceName("testtesttesttest", 10)
	if len(result) != 10 {
		t.Fatalf("Expected length 10, got %d", len(result))
	}

	interfaceSlice := []interface{}{"test1", "test2"}
	stringSlice := convertInterfaceToString(interfaceSlice)
	if len(stringSlice) != 2 || stringSlice[0] != "test1" || stringSlice[1] != "test2" {
		t.Fatalf("Expected ['test1', 'test2'], got %v", stringSlice)
	}

	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result = composeName("-", []string{"prefix"}, "name", "slug", []string{"suffix"}, "random", 20, namePrecedence)
	if result == "" {
		t.Fatal("Expected non-empty result from composeName")
	}

	valid, err := validateResourceType("azurerm_resource_group", []string{})
	if !valid || err != nil {
		t.Fatalf("validateResourceType returned error: %v", err)
	}
}

func TestStateUpgrade(t *testing.T) {
	state := map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_resource_group",
		"prefixes":      []interface{}{"prefix"},
		"suffixes":      []interface{}{"suffix"},
	}

	newState, err := resourceNameStateUpgradeV2(context.Background(), state, nil)
	if err != nil {
		t.Fatalf("resourceNameStateUpgradeV2 returned error: %v", err)
	}

	if newState["name"] != "test" {
		t.Fatalf("Expected name 'test', got %s", newState["name"])
	}
	if newState["resource_type"] != "azurerm_resource_group" {
		t.Fatalf("Expected resource_type 'azurerm_resource_group', got %s", newState["resource_type"])
	}
}

func TestResourceNameAllResourceTypes(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name resource")
	}

	resourceTypes := []string{
		"azurerm_resource_group",
		"azurerm_virtual_network",
		"azurerm_subnet",
		"azurerm_network_security_group",
		"azurerm_application_security_group",
		"azurerm_network_interface",
		"azurerm_public_ip",
		"azurerm_public_ip_prefix",
		"azurerm_storage_account",
		"azurerm_virtual_machine",
		"azurerm_windows_virtual_machine",
		"azurerm_linux_virtual_machine",
		"azurerm_managed_disk",
		"azurerm_availability_set",
		"azurerm_app_service",
		"azurerm_app_service_plan",
		"azurerm_function_app",
		"azurerm_key_vault",
		"azurerm_kubernetes_cluster",
		"azurerm_container_registry",
		"azurerm_application_gateway",
		"azurerm_firewall",
		"azurerm_sql_server",
		"azurerm_mysql_server",
		"azurerm_postgresql_server",
		"azurerm_cosmosdb_account",
		"azurerm_eventhub_namespace",
		"azurerm_servicebus_namespace",
		"azurerm_databricks_workspace",
		"azurerm_log_analytics_workspace",
	}

	for _, resourceType := range resourceTypes {
		d := nameResource.TestResourceData()
		d.Set("name", "testname")
		d.Set("resource_type", resourceType)
		d.Set("random_length", 5)

		err := nameResource.Create(d, nil)
		if err != nil {
			t.Fatalf("nameResource.Create returned error for resource type %s: %v", resourceType, err)
		}

		result := d.Get("result").(string)
		if result == "" {
			t.Fatalf("Expected non-empty result for resource type %s", resourceType)
		}
	}
}
