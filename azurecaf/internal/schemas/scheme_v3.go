package schemas

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const V3 = 3

func ResourceNameStateUpgradeV3(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return rawState, nil
}

func V3_Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prefixes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"suffixes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"random_length": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"random_seed": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clean_input": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
		},
	}
}
