package e2e

import (
	"os"
	"testing"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider = azurecaf.Provider()
var testAccProviders = map[string]*schema.Provider{
	"azurecaf": testAccProvider,
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Fatal("TF_ACC must be set for acceptance tests")
	}
}
