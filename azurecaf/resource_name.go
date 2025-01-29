package azurecaf

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNameCreate,
		UpdateContext: resourceNameUpdate,
		ReadContext:   resourceNameRead,
		Delete:        schema.RemoveFromState,
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
	resourceTypes := convertInterfaceToString(d.Get("resource_types").([]interface{}))
	cleanInput := d.Get("clean_input").(bool)
	passthrough := d.Get("passthrough").(bool)
	useSlug := d.Get("use_slug").(bool)
	randomLength := d.Get("random_length").(int)
	randomSeed := int64(d.Get("random_seed").(int))
	randomString := d.Get("random_string").(string)
	randomSuffix := randSeq(randomLength, randomSeed)
	if len(randomString) > 0 {
		randomSuffix = randomString
	} else {
		if err := d.SetNew("random_string", randomSuffix); err != nil {
			return err
		}
	}
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, results, _, err :=
		getData(resourceType, resourceTypes, separator,
			prefixes, name, suffixes, randomSuffix,
			cleanInput, passthrough, useSlug, namePrecedence)
	if !d.GetRawState().IsNull() {
		if err := d.SetNew("result", result); err != nil {
			return err
		}
		if err := d.SetNew("results", results); err != nil {
			return err
		}
	}
	return err
}

func resourceNameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return getNameResult(d, meta)
}

func resourceNameUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return getNameResult(d, meta)
}

func resourceNameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
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
	resourceTypes := convertInterfaceToString(d.Get("resource_types").([]interface{}))
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
	randomSuffix := randSeq(randomLength, randomSeed)
	if len(randomString) > 0 {
		randomSuffix = randomString
	} else {
		if err := d.Set("random_string", randomSuffix); err != nil {
			return diag.FromErr(err)
		}
	}

	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	result, results, id, err :=
		getData(resourceType, resourceTypes, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(result) > 0 {
		if err := d.Set("result", result); err != nil {
			return diag.FromErr(err)
		}
	}
	if len(results) > 0 {
		if err := d.Set("results", results); err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	
	// Set the ID after all other operations are successful
	d.SetId(id)
	return nil
}

func getData(resourceType string, resourceTypes []string, separator string, prefixes []string, name string, suffixes []string, randomSuffix string, cleanInput bool, passthrough bool, useSlug bool, namePrecedence []string) (result string, results map[string]string, id string, err error) {
	validationErr := validateResourceTypes(resourceType, resourceTypes)
	if validationErr != nil {
		err = validationErr
		return
	}
	if results == nil {
		results = make(map[string]string)
	}
	ids := []string{}
	if len(resourceType) > 0 {
		var resourceResult string
		resourceResult, err = getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return
		}
		result = resourceResult
		results[resourceType] = resourceResult
		ids = append(ids, fmt.Sprintf("%s\t%s", resourceType, resourceResult))
	}

	for _, resourceTypeName := range resourceTypes {
		var resourceResult string
		resourceResult, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, cleanInput, passthrough, useSlug, namePrecedence)
		if err != nil {
			return
		}
		results[resourceTypeName] = resourceResult
		ids = append(ids, fmt.Sprintf("%s\t%s", resourceTypeName, resourceResult))
	}
	id = b64.StdEncoding.EncodeToString([]byte(strings.Join(ids, "\n")))
	return
}
