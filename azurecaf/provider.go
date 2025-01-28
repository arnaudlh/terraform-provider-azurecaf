package azurecaf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the provider schema to Terraform.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		ResourcesMap: map[string]*schema.Resource{
			"azurecaf_name": resourceName(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azurecaf_name": dataSourceName(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azurecaf_environment_variable": dataEnvironmentVariable(),
			"azurecaf_name":                 dataName(),
		},
	}
}
