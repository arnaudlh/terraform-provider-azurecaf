package schemas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func V4_Schema() *schema.Resource {
	baseSchema := BaseSchema()
	resourceMapsKeys := getResourceMaps()

	schemaMap := make(map[string]*schema.Schema)
	for k, v := range baseSchema {
		newSchema := *v
		newSchema.ForceNew = false
		schemaMap[k] = &newSchema
	}

	schemaMap["resource_types"] = &schema.Schema{
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
		},
		Optional: true,
		ForceNew: false,
	}

	return &schema.Resource{
		Schema:         schemaMap,
		SchemaVersion: 4,
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func ResourceNameStateUpgradeV4(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return rawState, nil
}
