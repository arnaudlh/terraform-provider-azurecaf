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
		Schema:        schemas.V4_Schema(),
		CustomizeDiff: getDifference,
	}
}

func getDifference(context context.Context, d *schema.ResourceDiff, resource interface{}) error {
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
	randomString := d.Get("random_string").(string)
	randomSuffix := utils.RandSeq(randomLength, randomSeed)
	if len(randomString) > 0 {
		randomSuffix = randomString
	} else {
	if err := d.SetNew("random_string", randomSuffix); err != nil {
		return fmt.Errorf("failed to set random_string")
	}
	}
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return fmt.Errorf("failed to get resource name: %s", err.Error())
	}
	
	results := make(map[string]string)
	results[resourceType] = result
	if !d.GetRawState().IsNull() {
	if err := d.SetNew("result", result); err != nil {
		return fmt.Errorf("failed to set result")
	}
	if err := d.SetNew("results", results); err != nil {
		return fmt.Errorf("failed to set results")
	}
	}

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
	randomSeed := int64(d.Get("random_seed").(int))

	if randomSeed == 0 {
		randomSeed = time.Now().UnixMicro()
		if err := d.Set("random_seed", randomSeed); err != nil {
			return diag.FromErr(err)
		}
	}
	randomString := d.Get("random_string").(string)
	randomSuffix := utils.RandSeq(randomLength, randomSeed)
	if len(randomString) > 0 {
		randomSuffix = randomString
	} else {
		if err := d.Set("random_string", randomSuffix); err != nil {
			return diag.FromErr(err)
		}
	}

	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return diag.FromErr(err)
	}
	
	id := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s\t%s", resourceType, result)))
	if len(result) > 0 {
		if err := d.Set("result", result); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(id)
	return diags
}
