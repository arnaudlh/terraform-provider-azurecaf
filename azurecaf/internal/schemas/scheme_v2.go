package schemas

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zclconf/go-cty/cty"
)

func V2() *schema.Resource {
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
			Type: schema.TypeString,
		},
		Optional:     true,
		ForceNew:     false,
		ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
	}

	return &schema.Resource{
		Schema: schema,
		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    V1().Schema,
				Upgrade: ResourceNameStateUpgradeV2,
				Version: 1,
			},
		},
		Create: resourceNameCreate,
		Read:   resourceNameRead,
		Update: resourceNameUpdate,
		Delete: resourceNameDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			if err := ValidateResourceNameInSchema(d); err != nil {
				return fmt.Errorf("resource name validation failed: %v", err)
			}
			return nil
		},
	}
}

func V1() *schema.Resource {
	return &schema.Resource{
		Schema: BaseSchema(),
	}
}

func ResourceNameStateUpgradeV2(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["use_slug"] = true
	return rawState, nil
}
