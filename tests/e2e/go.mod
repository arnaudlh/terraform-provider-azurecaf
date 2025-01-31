module github.com/aztfmod/terraform-provider-azurecaf/tests/e2e

go 1.22.0

require (
	github.com/aztfmod/terraform-provider-azurecaf v0.0.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.35.0
)

replace github.com/aztfmod/terraform-provider-azurecaf => ../..
