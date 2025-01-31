// Package azurecaf provides the Azure CAF naming provider functionality.
package azurecaf

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	models "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataName() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(models.ResourceDefinitions))
	for k := range models.ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}

	return &schema.Resource{
		ReadContext: dataNameRead,
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

func dataNameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	resourceType := d.Get("resource_type").(string)

	// Get resource definition for validation
	resource, ok := models.ResourceDefinitions[resourceType]
	if !ok {
		return diag.Errorf("resource type %s not found", resourceType)
	}

	// Validate name against resource pattern
	validationRegEx, err := regexp.Compile(resource.ValidationRegExp)
	if err != nil {
		return diag.Errorf("invalid regex pattern for %s: %v", resourceType, err)
	}

	// Validate length constraints
	if len(name) > resource.MaxLength {
		return diag.Errorf("name length %d exceeds maximum length %d for resource type %s", len(name), resource.MaxLength, resourceType)
	}

	// Validate case requirements
	if resource.LowerCase && name != strings.ToLower(name) {
		return diag.Errorf("name must be lowercase for resource type %s", resourceType)
	}

	// Validate against regex pattern
	if !validationRegEx.MatchString(name) {
		return diag.Errorf("invalid name for %s: %s does not match pattern %s", resourceType, name, resource.ValidationRegExp)
	}

	id := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s\t%s", resourceType, name)))
	d.SetId(id)

	if err := d.Set("result", name); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
