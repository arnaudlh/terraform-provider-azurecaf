package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResource_MappedResource(t *testing.T) {
	resourceType := "azurerm_resource_group"
	resource, err := getResource(resourceType)

	if err != nil {
		t.Fatalf("getResource returned error for valid mapped resource: %v", err)
	}

	if resource == nil {
		t.Fatal("Expected resource to be returned for valid mapped resource")
	}
}

func TestTrimResourceName_NoTrim(t *testing.T) {
	resourceName := "short-name"
	maxLength := 20

	result := trimResourceName(resourceName, maxLength)

	if result != resourceName {
		t.Fatalf("Expected %s, got %s", resourceName, result)
	}
}

func TestGetNameResult_MultipleResourceTypes(t *testing.T) {
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
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{"azurerm_resource_group", "azurerm_storage_account"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error for multiple resource types: %v", err)
	}

	results := d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results["azurerm_resource_group"] == "" {
		t.Fatal("Expected non-empty result for azurerm_resource_group")
	}

	if results["azurerm_storage_account"] == "" {
		t.Fatal("Expected non-empty result for azurerm_storage_account")
	}
}

func TestGetResourceName_AllParameters(t *testing.T) {
	resourceType := "azurerm_resource_group"
	separator := "-"
	prefixes := []string{"prefix1", "prefix2"}
	name := "testname"
	suffixes := []string{"suffix1", "suffix2"}
	randomSuffix := "random"
	convention := ConventionCafClassic
	cleanInput := true
	passthrough := false
	useSlug := true
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)

	if err != nil {
		t.Fatalf("getResourceName returned error: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}

func TestGetResourceName_Passthrough(t *testing.T) {
	resourceType := "azurerm_resource_group"
	separator := "-"
	prefixes := []string{"prefix"}
	name := "testname"
	suffixes := []string{"suffix"}
	randomSuffix := "random"
	convention := ConventionCafClassic
	cleanInput := true
	passthrough := true
	useSlug := true
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}

	result, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)

	if err != nil {
		t.Fatalf("getResourceName returned error: %v", err)
	}

	if result != name {
		t.Fatalf("Expected result to be %s, got %s", name, result)
	}
}

func TestGetResult_AllParameters(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"prefix": {
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

	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error with ConventionCafRandom: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result with ConventionCafRandom")
	}

	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error with ConventionRandom: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result with ConventionRandom")
	}

	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefix", "prefix")
	d.Set("postfix", "postfix")
	d.Set("convention", ConventionPassThrough)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error with ConventionPassThrough: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result with ConventionPassThrough")
	}
}
