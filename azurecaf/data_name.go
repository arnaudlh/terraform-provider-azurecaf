package azurecaf

import (
	"context"
	"fmt"

	"github.com/arnaudlh/terraform-provider-azurecaf/azurecaf/internal/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataName() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(getResourceDefinitions()))
	for k := range getResourceDefinitions() {
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
	err := getNameReadResult(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func getNameReadResult(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	prefixes := convertInterfaceToString(d.Get("prefixes").([]interface{}))
	suffixes := convertInterfaceToString(d.Get("suffixes").([]interface{}))
	separator := d.Get("separator").(string)
	resourceType := d.Get("resource_type").(string)
	cleanInput := d.Get("clean_input").(bool)
	passthrough := d.Get("passthrough").(bool)
	useSlug := d.Get("use_slug").(bool)
	randomLength := d.Get("random_length").(int)
	randomSeed := int64(d.Get("random_seed").(int))

	randomSuffix := randSeq(randomLength, randomSeed)

	// Ensure the correct types are passed to getResourceName
	resourceName, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, []string{})
	if err != nil {
		return err
	}
	if err := d.Set("result", resourceName); err != nil {
		return fmt.Errorf("error setting result: %w", err)
	}
	return nil
}

func getResourceDefinitions() map[string]interface{} {
	resourceDefs := make(map[string]interface{})
	resourceMapsKeys := schemas.GetResourceMaps()
	for _, k := range resourceMapsKeys {
		resourceDefs[k] = nil
	}
	return resourceDefs
}
