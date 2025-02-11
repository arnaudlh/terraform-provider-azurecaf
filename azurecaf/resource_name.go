package azurecaf

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"time"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/schemas"
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNameCreate,
		UpdateContext: resourceNameUpdate,
		ReadContext:   resourceNameRead,
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			d.SetId("")
			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 4,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    schemas.V2().CoreConfigSchema().ImpliedType(),
				Upgrade: schemas.ResourceNameStateUpgradeV2,
				Version: 2,
			},
			{
				Type:    schemas.V3().CoreConfigSchema().ImpliedType(),
				Upgrade: schemas.ResourceNameStateUpgradeV3,
				Version: 3,
			},
		},
		Schema:        schemas.V4_Schema().Schema,
		CustomizeDiff: getDifference,
	}
}

func getDifference(context context.Context, d *schema.ResourceDiff, resource interface{}) error {
	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}
	prefixesRaw := d.Get("prefixes")
	var prefixes []string
	if prefixesRaw != nil {
		if arr, ok := prefixesRaw.([]interface{}); ok {
			prefixes = convertInterfaceToString(arr)
		}
	}
	suffixesRaw := d.Get("suffixes")
	var suffixes []string
	if suffixesRaw != nil {
		if arr, ok := suffixesRaw.([]interface{}); ok {
			suffixes = convertInterfaceToString(arr)
		}
	}
	separator, ok := d.Get("separator").(string)
	if !ok {
		separator = "-"
	}
	resourceType, ok := d.Get("resource_type").(string)
	if !ok {
		return fmt.Errorf("resource_type must be a string")
	}
	cleanInput, _ := d.Get("clean_input").(bool)
	passthrough, _ := d.Get("passthrough").(bool)
	useSlug, _ := d.Get("use_slug").(bool)
	randomLength, _ := d.Get("random_length").(int)
	// Get or generate random seed
	var randomSeed int64
	if seedRaw, ok := d.GetOk("random_seed"); ok {
		randomSeed = int64(seedRaw.(int))
	} else if d.Id() == "" {
		// Only generate new seed for new resources
		randomSeed = time.Now().UnixNano()
		if err := d.SetNew("random_seed", randomSeed); err != nil {
			return fmt.Errorf("failed to set random_seed: %v", err)
		}
	} else {
		// For existing resources, preserve the current state
		return nil
	}

	// Handle random string generation
	var randomSuffix string
	if randomString, ok := d.GetOk("random_string"); ok && randomString.(string) != "" {
		randomSuffix = randomString.(string)
	} else if randomLength > 0 {
		// Generate random string using seed
		randomSuffix = utils.RandSeq(randomLength, randomSeed)
		if err := d.SetNew("random_string", randomSuffix); err != nil {
			return fmt.Errorf("failed to set random_string: %v", err)
		}
	}

	// Always preserve the seed in state
	if err := d.SetNew("random_seed", randomSeed); err != nil {
		return fmt.Errorf("failed to preserve random_seed: %v", err)
	}
	namePrecedence := []string{"prefixes", "name", "slug", "random", "suffixes"}
	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return fmt.Errorf("failed to get resource name: %s", err.Error())
	}
	
	// Convert existing results to string map
	existingResults := make(map[string]string)
	if raw, ok := d.Get("results").(map[string]interface{}); ok {
		for k, v := range raw {
			if str, ok := v.(string); ok {
				existingResults[k] = str
			}
		}
	}
	existingResults[resourceType] = result
	if err := d.SetNew("result", result); err != nil {
		return fmt.Errorf("failed to set result: %v", err)
	}
	if err := d.SetNew("results", existingResults); err != nil {
		return fmt.Errorf("failed to set results: %v", err)
	}
	// Random string is already set above
	return nil
}

func resourceNameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return getNameResult(d, meta)
}

func resourceNameUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return getNameResult(d, meta)
}

func resourceNameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return getNameResult(d, meta)
}

func convertInterfaceToString(source []interface{}) []string {
	s := make([]string, len(source))
	for i, v := range source {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func getNameResult(d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	prefixes := convertInterfaceToString(d.Get("prefixes").([]interface{}))
	suffixes := convertInterfaceToString(d.Get("suffixes").([]interface{}))
	separator := d.Get("separator").(string)
	resourceType := d.Get("resource_type").(string)
	cleanInput := d.Get("clean_input").(bool)
	passthrough := d.Get("passthrough").(bool)
	useSlug := d.Get("use_slug").(bool)
	randomLength := d.Get("random_length").(int)
	var randomSeed int64
	if seedRaw, ok := d.GetOk("random_seed"); ok {
		randomSeed = int64(seedRaw.(int))
	} else if d.Id() == "" {
		randomSeed = time.Now().UnixNano()
		if err := d.Set("random_seed", randomSeed); err != nil {
			return diag.FromErr(err)
		}
	}
	var randomSuffix string
	if randomString, ok := d.Get("random_string").(string); ok && len(randomString) > 0 {
		randomSuffix = randomString
	} else if randomLength > 0 {
		randomSuffix = utils.RandSeq(randomLength, randomSeed)
		if err := d.Set("random_string", randomSuffix); err != nil {
			return diag.FromErr(err)
		}
	}

	namePrecedence := []string{"prefixes", "name", "slug", "random", "suffixes"}
	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return diag.FromErr(err)
	}
	
	id := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s\t%s", resourceType, result)))
	d.SetId(id)

	// Convert existing results to string map
	existingResults := make(map[string]string)
	if raw, ok := d.Get("results").(map[string]interface{}); ok {
		for k, v := range raw {
			if str, ok := v.(string); ok {
				existingResults[k] = str
			}
		}
	}
	existingResults[resourceType] = result

	if err := d.Set("result", result); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("results", existingResults); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
