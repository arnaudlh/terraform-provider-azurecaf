package azurecaf

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataEnvironmentVariable(t *testing.T) {
	resource := dataEnvironmentVariable()
	if resource.Schema["name"] == nil {
		t.Fatal("Expected name attribute in schema")
	}
	if resource.Schema["fails_if_empty"] == nil {
		t.Fatal("Expected fails_if_empty attribute in schema")
	}
	if resource.Schema["value"] == nil {
		t.Fatal("Expected value attribute in schema")
	}
}

func TestResourceAction(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"fails_if_empty": {
				Type: schema.TypeBool,
			},
			"value": {
				Type: schema.TypeString,
			},
		},
	}

	testVarName := "TEST_ENV_VAR_FOR_AZURECAF"
	testVarValue := "test_value"
	os.Setenv(testVarName, testVarValue)
	defer os.Unsetenv(testVarName)

	d := r.TestResourceData()
	d.Set("name", testVarName)
	d.Set("fails_if_empty", true)

	diags := resourceAction(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("resourceAction returned error for existing env var: %v", diags)
	}

	value := d.Get("value").(string)
	if value != testVarValue {
		t.Fatalf("Expected value %s, got %s", testVarValue, value)
	}

	nonExistingVar := "NON_EXISTING_ENV_VAR"
	os.Unsetenv(nonExistingVar) // Make sure it doesn't exist

	d = r.TestResourceData()
	d.Set("name", nonExistingVar)
	d.Set("fails_if_empty", true)

	diags = resourceAction(context.Background(), d, nil)
	if !diags.HasError() {
		t.Fatal("Expected error for non-existing env var but got none")
	}
}
