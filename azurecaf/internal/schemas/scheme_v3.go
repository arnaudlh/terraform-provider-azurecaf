package schemas

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/go-cty/cty"
)

func V3() *schema.Resource {
	baseSchema := BaseSchema()
	resourceMapsKeys := getResourceMaps()

	schemaMap := make(map[string]*schema.Schema)
	for k, v := range baseSchema {
		newSchema := *v
		newSchema.ForceNew = true
		schemaMap[k] = &newSchema
	}

	schemaMap["resource_types"] = &schema.Schema{
		Type:     schema.TypeList,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
		},
		Optional: true,
		ForceNew: true,
	}

	schemaMap["use_slug"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: true,
		Default:  true,
	}

	return &schema.Resource{
		Schema:         schemaMap,
		SchemaVersion: 3,
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

func ResourceNameStateUpgradeV3(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return nil, nil
	}

	results := rawState["results"]
	content := make(map[string]interface{})
	if results != nil {
		switch v := results.(type) {
		case map[string]interface{}:
			content = v
		}
	}

	resourceType, ok := rawState["resource_type"].(string)
	if !ok {
		resourceType = ""
	}
	result, ok := rawState["result"].(string)
	if !ok {
		result = ""
	}

	if resourceType != "" && result != "" {
		if _, ok := content[resourceType]; !ok {
			content[resourceType] = result
		}
	}

	rawState["results"] = content

	if len(content) > 0 {
		ids := make([]string, 0, len(content))
		for k, v := range content {
			ids = append(ids, fmt.Sprintf("%s\t%s", k, v.(string)))
		}
		rawState["id"] = b64.StdEncoding.EncodeToString([]byte(strings.Join(ids, "\n")))
	}

	return rawState, nil
}
