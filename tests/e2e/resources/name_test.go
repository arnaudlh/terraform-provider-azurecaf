package resources

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
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

func TestResourceName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "azurecaf_name" "rg_basic" {
  name           = "example"
  resource_type  = "azurerm_resource_group"
  random_length  = 5
  random_seed    = 12345
  clean_input    = true
}

resource "azurecaf_name" "rg_with_slug" {
  name           = "example"
  resource_type  = "azurerm_resource_group"
  random_length  = 5
  random_seed    = 12345
  clean_input    = true
  use_slug       = true
}

resource "azurecaf_name" "rg_with_prefixes" {
  name           = "example"
  resource_type  = "azurerm_resource_group"
  prefixes       = ["dev", "test"]
  random_length  = 5
  random_seed    = 12345
  clean_input    = true
}

resource "azurecaf_name" "rg_with_suffixes" {
  name           = "example"
  resource_type  = "azurerm_resource_group"
  suffixes       = ["001", "prod"]
  random_length  = 5
  random_seed    = 12345
  clean_input    = true
}

resource "azurecaf_name" "acr_name" {
  name           = "myregistry"
  resource_type  = "azurerm_container_registry"
  random_length  = 5
  random_seed    = 12345
  clean_input    = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					// Basic resource group name validation
					resource.TestCheckResourceAttr("azurecaf_name.rg_basic", "random_seed", "12345"),
					resource.TestCheckResourceAttrSet("azurecaf_name.rg_basic", "result"),
					resource.TestMatchResourceAttr("azurecaf_name.rg_basic", "result", 
						regexp.MustCompile("^[a-zA-Z0-9-_]+$")),

					// Resource group name with slug validation
					resource.TestCheckResourceAttrSet("azurecaf_name.rg_with_slug", "result"),
					resource.TestMatchResourceAttr("azurecaf_name.rg_with_slug", "result", 
						regexp.MustCompile("^rg-[a-zA-Z0-9-_]+$")),

					// Resource group name with prefixes validation
					resource.TestCheckResourceAttrSet("azurecaf_name.rg_with_prefixes", "result"),
					resource.TestMatchResourceAttr("azurecaf_name.rg_with_prefixes", "result", 
						regexp.MustCompile("^dev-test-[a-zA-Z0-9-_]+$")),

					// Resource group name with suffixes validation
					resource.TestCheckResourceAttrSet("azurecaf_name.rg_with_suffixes", "result"),
					resource.TestMatchResourceAttr("azurecaf_name.rg_with_suffixes", "result", 
						regexp.MustCompile("^[a-zA-Z0-9-_]+-001-prod$")),

					// Container registry name validation (specific format)
					resource.TestCheckResourceAttrSet("azurecaf_name.acr_name", "result"),
					resource.TestMatchResourceAttr("azurecaf_name.acr_name", "result", 
						regexp.MustCompile("^cr-[a-zA-Z0-9-]+-[a-zA-Z0-9]+$")),
				),
			},
		},
	})
}
