package azurecaf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataName() *schema.Resource {
	return &schema.Resource{
		Read:   resourceName().Read,
		Schema: resourceName().Schema,
	}
}
