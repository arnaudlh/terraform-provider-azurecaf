package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceNameDelete_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameDeleteComprehensiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.test_delete_comprehensive",
						"prefix-rg-testname-suffix",
						25,
						"prefix"),
				),
			},
		},
	})
}

func TestAccResourceNamingConventionDelete_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNamingConventionDeleteComprehensiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_naming_convention.test_delete_comprehensive",
						"prefix-rg-testname-suffix",
						80,
						"prefix"),
				),
			},
		},
	})
}

func TestAccResourceNameConcatenateParameters_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameConcatenateParametersComprehensiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.test_concat_comprehensive",
						"prefix1-prefix2-rg-testname-suffix1-suffix2",
						43,
						"prefix1-prefix2"),
				),
			},
		},
	})
}

func TestAccResourceNameStateUpgradeV2_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameStateUpgradeV2ComprehensiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.test_upgrade_comprehensive",
						"prefix-rg-testname-suffix",
						25,
						"prefix"),
				),
			},
		},
	})
}

func TestAccDataNameRead_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNameReadComprehensiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"data.azurecaf_name.test_read_comprehensive",
						"prefix-rg-testname-suffix",
						25,
						"prefix"),
				),
			},
		},
	})
}

func TestAccDataEnvironmentVariableRead_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnvironmentVariableReadComprehensiveConfig,
				Check: resource.TestCheckResourceAttrSet(
					"data.azurecaf_environment_variable.test_read_comprehensive", "value"),
			},
		},
	})
}

func TestAccResourceNamingConventionGetResult_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNamingConventionGetResultComprehensiveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_naming_convention.test_get_result_comprehensive",
						"prefix-rg-testname-suffix",
						80,
						"prefix"),
				),
			},
		},
	})
}

func TestAccResourceNameAllResourceTypes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameAllResourceTypesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.test_all_resource_types",
						"prefix-rg-testname-suffix",
						25,
						"prefix"),
				),
			},
		},
	})
}

func TestAccResourceNamingConventionAllConventions(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNamingConventionAllConventionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_naming_convention.test_all_conventions",
						"prefix-rg-testname-suffix",
						80,
						"prefix"),
				),
			},
		},
	})
}

func TestAccResourceNameMultipleResourceTypes_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameMultipleResourceTypesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"azurecaf_name.test_multiple_resource_types", "results.azurerm_resource_group"),
					resource.TestCheckResourceAttrSet(
						"azurecaf_name.test_multiple_resource_types", "results.azurerm_storage_account"),
				),
			},
		},
	})
}

func TestAccResourceNamePassthrough_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNamePassthroughConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"azurecaf_name.test_passthrough", "result", "testname"),
				),
			},
		},
	})
}

func TestAccResourceNameCleanInput_Comprehensive(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameCleanInputConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"azurecaf_name.test_clean_input", "result"),
				),
			},
		},
	})
}

const testAccResourceNameDeleteComprehensiveConfig = `
resource "azurecaf_name" "test_delete_comprehensive" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccResourceNamingConventionDeleteComprehensiveConfig = `
resource "azurecaf_naming_convention" "test_delete_comprehensive" {
    name            = "testname"
    convention      = "cafrandom"
    resource_type   = "rg"
    prefix          = "prefix"
    postfix         = "suffix"
}
`

const testAccResourceNameConcatenateParametersComprehensiveConfig = `
resource "azurecaf_name" "test_concat_comprehensive" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix1", "prefix2"]
    suffixes        = ["suffix1", "suffix2"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccResourceNameStateUpgradeV2ComprehensiveConfig = `
resource "azurecaf_name" "test_upgrade_comprehensive" {
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

const testAccDataNameReadComprehensiveConfig = `
data "azurecaf_name" "test_read_comprehensive" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccDataEnvironmentVariableReadComprehensiveConfig = `
data "azurecaf_environment_variable" "test_read_comprehensive" {
    name            = "PATH"
    fails_if_empty  = false
}
`

const testAccResourceNamingConventionGetResultComprehensiveConfig = `
resource "azurecaf_naming_convention" "test_get_result_comprehensive" {
    name            = "testname"
    convention      = "cafrandom"
    resource_type   = "rg"
    prefix          = "prefix"
    postfix         = "suffix"
}
`

const testAccResourceNameAllResourceTypesConfig = `
resource "azurecaf_name" "test_all_resource_types" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccResourceNamingConventionAllConventionsConfig = `
resource "azurecaf_naming_convention" "test_all_conventions" {
    name            = "testname"
    convention      = "cafrandom"
    resource_type   = "rg"
    prefix          = "prefix"
    postfix         = "suffix"
}
`

const testAccResourceNameMultipleResourceTypesConfig = `
resource "azurecaf_name" "test_multiple_resource_types" {
    name            = "testname"
    resource_types  = ["azurerm_resource_group", "azurerm_storage_account"]
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`

const testAccResourceNamePassthroughConfig = `
resource "azurecaf_name" "test_passthrough" {
    name            = "testname"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["prefix"]
    suffixes        = ["suffix"]
    separator       = "-"
    clean_input     = true
    passthrough     = true
    random_length   = 0
}
`

const testAccResourceNameCleanInputConfig = `
resource "azurecaf_name" "test_clean_input" {
    name            = "test@name"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["pre@fix"]
    suffixes        = ["suf@fix"]
    separator       = "-"
    clean_input     = true
    random_length   = 0
}
`
