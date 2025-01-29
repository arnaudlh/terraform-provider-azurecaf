package azurecaf

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"time"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
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
				Required:     true,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew:     true,
				Description:  "The resource type to generate a name for",
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
	if err := validateResourceType(d.Get("resource_type").(string)); err != nil {
		return diag.FromErr(err)
	}

	if err := getNameReadResult(d, meta); err != nil {
		return diag.FromErr(fmt.Errorf("error generating name: %w", err))
	}

	return diag.Diagnostics{}
}

func validateResourceType(resourceType string) error {
	return validateResourceTypes(resourceType, nil)
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

	if randomSeed == 0 {
		randomSeed = time.Now().UnixMicro()
		if err := d.Set("random_seed", randomSeed); err != nil {
			return fmt.Errorf("error setting random_seed: %w", err)
		}
	}

	randomSuffix := randSeq(randomLength, randomSeed)
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return err
	}

	if err := d.Set("result", result); err != nil {
		return fmt.Errorf("error setting result: %w", err)
	}

	// Set the ID after all other operations are successful
	d.SetId(b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s\t%s", resourceType, result))))
	return nil
}
