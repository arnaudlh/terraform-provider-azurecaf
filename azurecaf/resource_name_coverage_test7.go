package azurecaf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceNameReadCoverage7(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"separator": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clean_input": {
				Type: schema.TypeBool,
			},
			"passthrough": {
				Type: schema.TypeBool,
			},
			"use_slug": {
				Type: schema.TypeBool,
			},
			"random_length": {
				Type: schema.TypeInt,
			},
			"random_seed": {
				Type: schema.TypeInt,
			},
			"result": {
				Type: schema.TypeString,
			},
			"results": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := resourceNameRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNameRead returned error: %v", err)
	}
}

func TestResourceNamingConventionReadCoverage7(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefix": {
				Type: schema.TypeString,
			},
			"postfix": {
				Type: schema.TypeString,
			},
			"convention": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"max_length": {
				Type: schema.TypeInt,
			},
			"result": {
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err := resourceNamingConventionRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionRead returned error: %v", err)
	}
}

func TestDataNameReadCoverage7(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"separator": {
				Type: schema.TypeString,
			},
			"resource_type": {
				Type: schema.TypeString,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clean_input": {
				Type: schema.TypeBool,
			},
			"passthrough": {
				Type: schema.TypeBool,
			},
			"use_slug": {
				Type: schema.TypeBool,
			},
			"random_length": {
				Type: schema.TypeInt,
			},
			"random_seed": {
				Type: schema.TypeInt,
			},
			"result": {
				Type: schema.TypeString,
			},
			"results": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	diags := dataNameRead(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("dataNameRead returned error: %v", diags)
	}
}

func TestDataEnvironmentVariableReadCoverage7(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"fails_if_empty": {
				Type: schema.TypeBool,
			},
			"value": {
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "PATH")
	d.Set("fails_if_empty", false)

	diags := resourceAction(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("resourceAction returned error: %v", diags)
	}
}

func TestGetSlugEdgeCases7(t *testing.T) {
	resourceType := ""
	convention := ConventionCafClassic
	
	slug := getSlug(resourceType, convention)
	
	if slug != "" {
		t.Fatalf("Expected empty slug for empty resource type, got '%s'", slug)
	}
	
	resourceType = "invalid_resource_type"
	
	slug = getSlug(resourceType, convention)
	
	if slug != "" {
		t.Fatalf("Expected empty slug for invalid resource type, got '%s'", slug)
	}
	
	resourceType = "azurerm_resource_group"
	
	slug1 := getSlug(resourceType, ConventionCafClassic)
	slug2 := getSlug(resourceType, ConventionCafRandom)
	slug3 := getSlug(resourceType, ConventionPassThrough)
	
	if slug1 != slug2 || slug1 != slug3 {
		t.Fatal("Expected same slug for different conventions")
	}
}

func TestComposeNameEdgeCases7(t *testing.T) {
	separator := "-"
	prefixes := []string{}
	name := ""
	slug := ""
	suffixes := []string{}
	randomSuffix := ""
	maxlength := 63
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result := composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result != "" {
		t.Fatalf("Expected empty result for empty parameters, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{}
	name = ""
	slug = ""
	suffixes = []string{}
	randomSuffix = "random"
	maxlength = 63
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result != randomSuffix {
		t.Fatalf("Expected result to be random suffix, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{"prefix1", "prefix2"}
	name = ""
	slug = ""
	suffixes = []string{}
	randomSuffix = ""
	maxlength = 63
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result != "prefix1-prefix2" {
		t.Fatalf("Expected result to be concatenated prefixes, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{}
	name = ""
	slug = ""
	suffixes = []string{"suffix1", "suffix2"}
	randomSuffix = ""
	maxlength = 63
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result != "suffix1-suffix2" {
		t.Fatalf("Expected result to be concatenated suffixes, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{"prefix1", "prefix2"}
	name = "testname"
	slug = "rg"
	suffixes = []string{"suffix1", "suffix2"}
	randomSuffix = "random"
	maxlength = 10
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if len(result) > maxlength {
		t.Fatalf("Expected result length to be <= %d, got %d", maxlength, len(result))
	}
}

func TestConcatenateParametersEdgeCases7(t *testing.T) {
	separator := "-"
	prefixes := []string{}
	name := []string{}
	suffixes := []string{}
	
	result := concatenateParameters(separator, prefixes, name, suffixes)
	
	if result != "" {
		t.Fatalf("Expected empty result for empty parameters, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{}
	name = []string{"testname"}
	suffixes = []string{}
	
	result = concatenateParameters(separator, prefixes, name, suffixes)
	
	if result != "testname" {
		t.Fatalf("Expected result to be name, got '%s'", result)
	}
	
	separator = "_"
	prefixes = []string{"prefix1", "prefix2"}
	name = []string{"testname"}
	suffixes = []string{"suffix1", "suffix2"}
	
	result = concatenateParameters(separator, prefixes, name, suffixes)
	
	if result != "prefix1_prefix2_testname_suffix1_suffix2" {
		t.Fatalf("Expected result with underscore separator, got '%s'", result)
	}
}

func TestCleanSliceEdgeCases7(t *testing.T) {
	data := []string{}
	resource := ResourceDefinitions["azurerm_resource_group"]
	
	result := cleanSlice(data, &resource)
	
	if len(result) != 0 {
		t.Fatalf("Expected empty result for empty slice, got %v", result)
	}
	
	data = []string{"test1", "test2"}
	
	result = cleanSlice(data, nil)
	
	if len(result) != len(data) {
		t.Fatalf("Expected same length for nil resource, got %d", len(result))
	}
	
	for i := range data {
		if result[i] != data[i] {
			t.Fatalf("Expected unchanged string at index %d, got '%s'", i, result[i])
		}
	}
}

func TestCleanStringEdgeCases7(t *testing.T) {
	name := ""
	resource := ResourceDefinitions["azurerm_resource_group"]
	
	result := cleanString(name, &resource)
	
	if result != "" {
		t.Fatalf("Expected empty result for empty string, got '%s'", result)
	}
	
	name = "test!@#$%^&*()_+name"
	
	result = cleanString(name, nil)
	
	if result != name {
		t.Fatalf("Expected unchanged string for nil resource, got '%s'", result)
	}
}

func TestConvertInterfaceToStringEdgeCases7(t *testing.T) {
	source := []interface{}{}
	
	result := convertInterfaceToString(source)
	
	if len(result) != 0 {
		t.Fatalf("Expected empty result for empty slice, got %v", result)
	}
	
	var nilSource []interface{}
	
	result = convertInterfaceToString(nilSource)
	
	if len(result) != 0 {
		t.Fatalf("Expected empty result for nil, got %v", result)
	}
	
	source = []interface{}{"test1", 2, true, nil, 3.14}
	
	result = convertInterfaceToString(source)
	
	if len(result) != len(source) {
		t.Fatalf("Expected length %d, got %d", len(source), len(result))
	}
	
	if result[3] != "<nil>" {
		t.Fatalf("Expected '<nil>' for nil, got '%s'", result[3])
	}
}

func TestRandSeqEdgeCases7(t *testing.T) {
	length := 0
	
	result := randSeq(length, nil)
	
	if result != "" {
		t.Fatalf("Expected empty result for zero length, got '%s'", result)
	}
	
	length = -1
	
	result = randSeq(length, nil)
	
	if result != "" {
		t.Fatalf("Expected empty result for negative length, got '%s'", result)
	}
	
	length = 10
	seed1 := int64(123)
	seed2 := int64(123)
	
	result1 := randSeq(length, &seed1)
	result2 := randSeq(length, &seed2)
	
	if result1 != result2 {
		t.Fatalf("Expected same result for same seed, got '%s' and '%s'", result1, result2)
	}
	
	seed1 = int64(123)
	seed2 = int64(456)
	
	result1 = randSeq(length, &seed1)
	result2 = randSeq(length, &seed2)
	
	if result1 == result2 {
		t.Fatalf("Expected different results for different seeds, got '%s' for both", result1)
	}
}
