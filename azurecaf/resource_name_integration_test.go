//go:build integration

package azurecaf

import (
	"regexp"
	"testing"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceName_CafClassic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.classic_rg",
						"pr1-pr2-rg-myrg-",
						29,
						"pr1-pr2"),
					regexMatch("azurecaf_name.classic_rg", regexp.MustCompile(models.ResourceDefinitions["azurerm_resource_group"].ValidationRegExp), 1),
				),
			},
			{
				Config: testAccResourceNameCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.classic_ca_invalid",
						"ca-myinvalidcaname",
						24,
						""),
					regexMatch("azurecaf_name.classic_ca_invalid", regexp.MustCompile(models.ResourceDefinitions["azurerm_container_app"].ValidationRegExp), 1),
				),
			},
			{
				Config: testAccResourceNameCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.passthrough",
						"passthrough",
						11,
						""),
					regexMatch("azurecaf_name.passthrough", regexp.MustCompile(models.ResourceDefinitions["azurerm_container_app"].ValidationRegExp), 1),
				),
			},
		},
	})
}

func TestAccResourceName_RsvCafClassic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNameCafClassicConfigRsv,
				Check: resource.ComposeTestCheckFunc(
					testAccCafNamingValidation(
						"azurecaf_name.rsv",
						"pr1-rsv-test-gm-su1",
						19,
						""),
					regexMatch("azurecaf_name.rsv", regexp.MustCompile(models.ResourceDefinitions["azurerm_recovery_services_vault"].ValidationRegExp), 1),
				),
			},
		},
	})
}
