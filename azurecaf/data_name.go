package azurecaf

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceName() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNameRead,
		Schema:     resourceName().Schema,
	}
}

func dataSourceNameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := getNameResult(d, m); err != nil {
		return diag.FromErr(fmt.Errorf("failed to get name result: %w", err))
	}

	result := d.Get("result").(string)
	d.SetId(result)
	return diags
}
