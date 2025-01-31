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
			"default_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default value to use if the environment variable is not set.",
			},
		},
	}
}

func resourceAction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	value, ok := os.LookupEnv(name)

	if !ok {
		if defaultValue, exists := d.GetOk("default_value"); exists {
			value = defaultValue.(string)
		} else if failsIfEmpty := d.Get("fails_if_empty").(bool); failsIfEmpty {
			return diag.Errorf("Value is not set for environment variable: %s", name)
		} else {
			value = ""
		}
	}

	d.SetId(name)
	if err := d.Set("value", value); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
