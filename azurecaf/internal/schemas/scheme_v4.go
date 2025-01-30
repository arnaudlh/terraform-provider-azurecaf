package schemas

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func V4() *schema.Resource {
	baseSchema := BaseSchema()
	resourceMapsKeys := getResourceMaps()

	schema := make(map[string]*schema.Schema)
	for k, v := range baseSchema {
		newSchema := *v
		newSchema.ForceNew = false
		schema[k] = &newSchema
	}

	schema["resource_types"] = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
		},
		Optional: true,
		ForceNew: false,
	}

	return &schema.Resource{
		Schema: schema,
		SchemaVersion: 4,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    V3().Schema,
				Upgrade: ResourceNameStateUpgradeV4,
				Version: 3,
			},
		},
		Create: resourceNameCreate,
		Read:   resourceNameRead,
		Update: resourceNameUpdate,
		Delete: resourceNameDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			if err := ValidateResourceNameInSchema(d); err != nil {
				return err
			}
			return nil
		},
	}
}

func ResourceNameStateUpgradeV4(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return rawState, nil
}
