package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResourceNameComprehensive(t *testing.T) {
	resourceTypeName := "azurerm_resource_group"
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

	result, err := getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned error: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	resourceTypeName = "invalid_resource_type"
	_, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}

	resourceTypeName = "azurerm_resource_group"
	passthrough = true
	result, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned error for passthrough = true: %v", err)
	}
	if result != name {
		t.Fatalf("Expected result to be %s for passthrough = true, got %s", name, result)
	}

	passthrough = false
	useSlug = false
	result, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned error for useSlug = false: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result for useSlug = false")
	}

	useSlug = true
	cleanInput = false
	name = "test!@#$%^&*()_+name"
	result, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned error for cleanInput = false and special characters: %v", err)
	}

	resourceTypeName = "azurerm_storage_account"
	name = "UPPERCASE"
	cleanInput = true
	result, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned error for storage account: %v", err)
	}
	for _, c := range result {
		if c >= 'A' && c <= 'Z' {
			t.Fatalf("Expected lowercase result for storage account, got %s", result)
		}
	}

	resourceTypeName = "azurerm_resource_group"
	name = "thisisaverylongnamethatwillbecutoffbythemaxlengthparameter"
	result, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err != nil {
		t.Fatalf("getResourceName returned error for long name: %v", err)
	}
	resource, _ := getResource(resourceTypeName)
	if len(result) > resource.MaxLength {
		t.Fatalf("Expected result length to be <= %d, got %d", resource.MaxLength, len(result))
	}

	name = "test@name"
	cleanInput = false
	_, err = getResourceName(resourceTypeName, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	if err == nil {
		t.Fatal("Expected error for invalid regex pattern but got none")
	}
}

func TestGetNameResultComprehensive(t *testing.T) {
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

	err := getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	d = r.TestResourceData()
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

	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}

	results := d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"azurerm_storage_account", "azurerm_virtual_network"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult returned error: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	results = d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	d = r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "")
	d.Set("resource_types", []interface{}{})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for empty resource_type and resource_types but got none")
	}
}
