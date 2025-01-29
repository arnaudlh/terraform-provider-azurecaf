package azurecaf

import (
	"context"
	"fmt"
	"sort"
	"strings"
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
	sort.Strings(resourceMapsKeys)

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
			"random_string": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
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

func expandStringList(input []interface{}) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}
	for _, v := range input {
		if v == nil {
			continue
		}
		output = append(output, v.(string))
	}
	return output
}

func validateResourceType(resourceType string) error {
	if _, ok := models.ResourceDefinitions[resourceType]; !ok {
		return fmt.Errorf("resource_type %q is not supported", resourceType)
	}
	return nil
}

func getNameReadResult(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	prefixes := expandStringList(d.Get("prefixes").([]interface{}))
	suffixes := expandStringList(d.Get("suffixes").([]interface{}))
	separator := d.Get("separator").(string)
	resourceType := d.Get("resource_type").(string)
	cleanInput := d.Get("clean_input").(bool)
	passthrough := d.Get("passthrough").(bool)
	useSlug := d.Get("use_slug").(bool)
	randomLength := d.Get("random_length").(int)
	randomSeed := int64(d.Get("random_seed").(int))

	// Use existing random string or generate a new one with the provided seed
	randomString := d.Get("random_string").(string)
	if randomString == "" && randomLength > 0 {
		if randomSeed == 0 {
			randomSeed = time.Now().UnixMicro()
			if err := d.Set("random_seed", randomSeed); err != nil {
				return fmt.Errorf("error setting random_seed: %w", err)
			}
		}
		randomString = randSeq(randomLength, randomSeed)
		if err := d.Set("random_string", randomString); err != nil {
			return fmt.Errorf("error setting random_string: %w", err)
		}
	}

	namePrecedence := []string{"name", "random", "suffixes", "prefixes"}
	slug := ""
	if useSlug {
		slug = getSlug(resourceType)
	}
	result, _, id, err := getData(resourceType, nil, separator, prefixes, name, suffixes, randomString, cleanInput, passthrough, false, namePrecedence)
	if err != nil {
		return fmt.Errorf("error generating name: %w", err)
	}

	// Add slug after name generation if useSlug is true
	if useSlug && len(slug) > 0 {
		if resourceType == "azurerm_storage_account" {
			// For storage accounts, handle without separators
			result = strings.ReplaceAll(result, separator, "")
			if len(prefixes) > 0 {
				// Insert slug after prefixes
				prefix := strings.Join(prefixes, "")
				rest := strings.TrimPrefix(result, prefix)
				result = prefix + slug + rest
			} else {
				result = slug + result
			}
			// For storage accounts, ensure prefix + st + rest format
			if strings.HasPrefix(result, slug) {
				result = strings.TrimPrefix(result, slug)
				if len(prefixes) > 0 {
					result = strings.Join(prefixes, "") + slug + result
				} else {
					result = slug + result
				}
			}
		} else {
			parts := strings.Split(result, separator)
			for i, part := range parts {
				if part == name {
					parts = append(parts[:i+1], parts[i:]...)
					parts[i] = slug
					break
				}
			}
			result = strings.Join(parts, separator)
		}
	}
	if err != nil {
		return fmt.Errorf("error generating name: %w", err)
	}

	if len(result) > 0 {
		if err := d.Set("result", result); err != nil {
			return fmt.Errorf("error setting result: %w", err)
		}
	}

	d.SetId(id)
	return nil
}
