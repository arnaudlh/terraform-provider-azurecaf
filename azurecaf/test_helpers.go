package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"azurecaf": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// No environment variables required for testing
}

func testAccCheckResourceDestroy(s *terraform.State) error {
	// No real resources to destroy in this provider
	return nil
}
