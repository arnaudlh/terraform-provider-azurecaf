package testutils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
)

// ProviderFactories returns a map of providers used for testing
func ProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"azurecaf": func() (*schema.Provider, error) {
			return azurecaf.Provider(), nil
		},
	}
}

// Provider returns an initialized Azure CAF provider for testing
func Provider() *schema.Provider {
	return azurecaf.Provider()
}
