package main

import (
	"os"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run gen.go

func main() {
	debug := os.Getenv("TF_PROVIDER_DEBUG") != ""
	level := hclog.Info
	if debug {
		level = hclog.Debug
	}
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "terraform-provider-azurecaf",
		Level:  level,
		Output: os.Stderr,
	})

	if debug {
		logger.Debug("Starting provider in debug mode")
	}

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			p := azurecaf.Provider()
			if debug {
				logger.Debug("Provider schema", "schema", p.Schema)
				logger.Debug("Provider resources", "resources", p.ResourcesMap)
				logger.Debug("Provider data sources", "data_sources", p.DataSourcesMap)
			}
			return p
		},
		Debug:  debug,
		Logger: logger,
		// Required for proper plugin protocol negotiation
		ProviderAddr: "registry.terraform.io/aztfmod/azurecaf",
		// Required for proper schema validation
		NoLogOutputOverride: true,
	}

	plugin.Serve(opts)
}
