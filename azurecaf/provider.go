package azurecaf

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the provider schema to Terraform.
func Provider() *schema.Provider {
	log.Printf("[INFO] Initializing Azure CAF provider")

	// Set debug logging
	debug := os.Getenv("TF_PROVIDER_DEBUG") != ""
	if debug {
		log.Printf("[DEBUG] Provider debug mode enabled")
		log.Printf("[DEBUG] Provider version: %s", os.Getenv("TF_PROVIDER_VERSION"))
	}

	// Initialize provider schema with proper validation
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"random_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Default random string length for all resources",
			},
			"random_seed": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Random seed for consistent resource naming",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"azurecaf_name": resourceName(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azurecaf_name":                dataName(),
			"azurecaf_environment_variable": dataEnvironmentVariable(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	// Validate the complete provider schema
	if err := p.InternalValidate(); err != nil {
		log.Printf("[ERROR] Failed to validate provider schema: %v", err)
		panic(err)
	}

	log.Printf("[DEBUG] Provider schema initialized with resources: %v", p.ResourcesMap)
	log.Printf("[DEBUG] Provider schema initialized with data sources: %v", p.DataSourcesMap)
	return p
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[DEBUG] Configuring Azure CAF provider")

	// Get provider binary directory for resource definitions
	ex, err := os.Executable()
	if err != nil {
		log.Printf("[WARN] Could not determine provider path: %v", err)
	} else {
		resourceDefPath := filepath.Join(filepath.Dir(ex), "resourceDefinition.json")
		log.Printf("[DEBUG] Looking for resource definitions at: %s", resourceDefPath)
		if _, err := os.Stat(resourceDefPath); err == nil {
			log.Printf("[DEBUG] Found resource definitions at: %s", resourceDefPath)
		}
	}

	config := &ProviderConfig{
		Version: os.Getenv("TF_PROVIDER_VERSION"),
	}

	if randomSeed, ok := d.GetOk("random_seed"); ok {
		config.RandomSeed = randomSeed.(int)
	}

	if randomLength, ok := d.GetOk("random_length"); ok {
		config.RandomLength = randomLength.(int)
	}

	log.Printf("[DEBUG] Provider configuration: %+v", config)
	return config, nil
}

// ProviderConfig holds the provider configuration
type ProviderConfig struct {
	Version      string
	RandomSeed   int
	RandomLength int
}
