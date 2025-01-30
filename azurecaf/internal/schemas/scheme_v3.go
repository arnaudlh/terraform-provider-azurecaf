package schemas

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func V3() *schema.Resource {
	resourceMapsKeys := getResourceMaps()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"random_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"results": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"separator": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "-",
			},
			"clean_input": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew:     true,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				},
				Optional: true,
				ForceNew: true,
			},
			"random_seed": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"use_slug": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
		},
	}
}

func ResourceNameStateUpgradeV3(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return nil, nil
	}

	// Handle results field
	results := rawState["results"]
	content := make(map[string]interface{})
	if results != nil {
		switch v := results.(type) {
		case map[string]interface{}:
			content = v
		}
	}

	// Handle resource_type and result fields safely
	resourceType, ok := rawState["resource_type"].(string)
	if !ok {
		resourceType = ""
	}
	result, ok := rawState["result"].(string)
	if !ok {
		result = ""
	}

	// Only update content if we have valid resource_type and result
	if resourceType != "" && result != "" {
		if _, ok := content[resourceType]; !ok {
			content[resourceType] = result
		}
	}

	rawState["results"] = content

	// Generate ID only if we have content
	if len(content) > 0 {
		ids := make([]string, 0, len(content))
		for k, v := range content {
			ids = append(ids, fmt.Sprintf("%s\t%s", k, v.(string)))
		}
		rawState["id"] = b64.StdEncoding.EncodeToString([]byte(strings.Join(ids, "\n")))
	}

	return rawState, nil
}
