package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider = azurecaf.Provider()
var testAccProviders = map[string]*schema.Provider{
	"azurecaf": testAccProvider,
}

func init() {
	// Copy resource definition file for tests
	if err := os.MkdirAll("resourceDefinitions", 0755); err != nil {
		panic(err)
	}
	src, err := os.ReadFile("../../resourceDefinition.json")
	if err == nil {
		if err := os.WriteFile("resourceDefinition.json", src, 0644); err != nil {
			panic(err)
		}
	}
}

func TestProvider(t *testing.T) {
	if err := testAccProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// Configure the provider
	if err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil)); err != nil {
		t.Fatal(err)
	}
}
