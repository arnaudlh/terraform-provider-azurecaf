package e2e

import (
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"azurecaf": providerserver.NewProtocol6WithError(azurecaf.New()),
}
