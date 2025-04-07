package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceNameDelete(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameDeleteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.test_delete",
						"prefix-rg-testname-suffix",
						25,
						"prefix"),
				),
			},
		},
	})
}

func TestAccResourceNamingConventionDelete(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNamingConventionDeleteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_naming_convention.test_delete",
						"prefix-rg-testname-suffix",
						80,
						"prefix"),
				),
			},
		},
	})
}

func TestAccResourceNameConcatenateParameters(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameConcatenateParametersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.test_concat",
						"prefix1-prefix2-rg-testname-suffix1-suffix2",
						43,
						"prefix1-prefix2"),
				),
			},
		},
	})
}

func TestAccResourceNameStateUpgradeV2New(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameStateUpgradeV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.test_upgrade",
						"prefix-rg-testname-suffix",
						25,
						"prefix"),
				),
			},
		},
	})
}

func TestAccDataNameRead(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNameReadConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"data.azurecaf_name.test_read",
						"prefix-rg-testname-suffix",
						25,
						"prefix"),
				),
			},
		},
	})
}

func TestAccDataEnvironmentVariableRead(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnvironmentVariableReadConfig,
				Check: resource.TestCheckResourceAttrSet(
					"data.azurecaf_environment_variable.test_read", "value"),
			},
		},
	})
}

func TestAccResourceNamingConventionGetResult(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNamingConventionGetResultConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_naming_convention.test_get_result",
						"prefix-rg-testname-suffix",
						80,
						"prefix"),
				),
			},
		},
	})
}

const testAccResourceNameDeleteConfig = `
resource "azurecaf_name" "test_delete" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccResourceNamingConventionDeleteConfig = `
resource "azurecaf_naming_convention" "test_delete" {
    name            = "testname"
    convention      = "cafrandom"
    resource_type   = "rg"
    prefix          = "prefix"
    postfix         = "suffix"
}
`

const testAccResourceNameConcatenateParametersConfig = `
resource "azurecaf_name" "test_concat" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix1", "prefix2"]
    suffixes        = ["suffix1", "suffix2"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccResourceNameStateUpgradeV2Config = `
resource "azurecaf_name" "test_upgrade" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
    use_slug        = true
}
`

const testAccDataNameReadConfig = `
data "azurecaf_name" "test_read" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccDataEnvironmentVariableReadConfig = `
data "azurecaf_environment_variable" "test_read" {
    name            = "PATH"
    fails_if_empty  = false
}
`

const testAccResourceNamingConventionGetResultConfig = `
resource "azurecaf_naming_convention" "test_get_result" {
    name            = "testname"
    convention      = "cafrandom"
    resource_type   = "rg"
    prefix          = "prefix"
    postfix         = "suffix"
}
`
