package azurecaf

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnvironmentVariable_basic(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironmentVariableConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azurecaf_environment_variable.test", "value", "test_value"),
				),
			},
		},
	})
}

func TestAccDataSourceEnvironmentVariable_default(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironmentVariableDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azurecaf_environment_variable.test_default", "value", "default_value"),
				),
			},
		},
	})
}

const testAccDataSourceEnvironmentVariableConfig = `
data "azurecaf_environment_variable" "test" {
  name = "TEST_ENV_VAR"
}
`

const testAccDataSourceEnvironmentVariableDefaultConfig = `
data "azurecaf_environment_variable" "test_default" {
  name = "NONEXISTENT_ENV_VAR"
  default_value = "default_value"
}
`
