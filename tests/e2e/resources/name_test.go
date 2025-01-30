package e2e

import (
	"testing"

	"github.com/aztfmod/terraform-provider-azurecaf/tests/e2e/testutils"
)

func TestResourceName(t *testing.T) {
	config := `
resource "azurecaf_name" "test" {
  name           = "example"
  resource_type  = "azurerm_resource_group"
  random_length  = 5
  clean_input    = true
}
`
	testutils.ResourceTest(t, "azurecaf_name.test", config)
}
