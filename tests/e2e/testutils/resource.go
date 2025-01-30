package testutils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ResourceTest runs a basic test for a resource configuration
func ResourceTest(t *testing.T, name string, config string) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		ProviderFactories: ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}
