package azurecaf

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
)

func dataEnvironmentVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceAction,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the environment variable.",
			},
			"fails_if_empty": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Throws an error if the environment variable is not set (default: false).",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the environment variable.",
				Sensitive:   true,
			},
		},
	}
}

func resourceAction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	value, ok := os.LookupEnv(name)

	if !ok {
		return diag.Errorf("Value is not set for environment variable: %s", name)
	}

	if err := d.Set("value", value); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	
	// Set the ID after all other operations are successful
	d.SetId(name)
	return nil
}
