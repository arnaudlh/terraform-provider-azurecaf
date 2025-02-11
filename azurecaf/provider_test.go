package azurecaf

import (
	"os"
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	// Initialize provider
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"azurecaf": testAccProvider,
	}

	// Ensure resource definitions are loaded
	if err := os.MkdirAll("resourceDefinitions", 0755); err != nil {
		panic(err)
	}

	// Copy resource definition file from parent directory
	src, err := os.ReadFile("../resourceDefinition.json")
	if err == nil {
		if err := os.WriteFile("resourceDefinitions/resourceDefinition.json", src, 0644); err != nil {
			panic(err)
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	// Configure the provider
	if err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil)); err != nil {
		t.Fatal(err)
	}
}

// Resources are local and no infrastructure is created in the test suite
func testAccCheckResourceDestroy(s *terraform.State) error {
	return nil
}
