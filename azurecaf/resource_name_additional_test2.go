package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResource_InvalidResource(t *testing.T) {
	resourceType := "invalid_resource_type"
	_, err := getResource(resourceType)
	
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
}

func TestTrimResourceName_WithTrim(t *testing.T) {
	resourceName := "this-is-a-very-long-resource-name-that-needs-to-be-trimmed"
	maxLength := 20
	
	result := trimResourceName(resourceName, maxLength)
	
	if len(result) != maxLength {
		t.Fatalf("Expected length %d, got %d", maxLength, len(result))
	}
	
	if result != resourceName[:maxLength] {
		t.Fatalf("Expected %s, got %s", resourceName[:maxLength], result)
	}
}

func TestValidateResourceType_EdgeCases(t *testing.T) {
	resourceType := ""
	resourceTypes := []string{}
	_, err := validateResourceType(resourceType, resourceTypes)
	
	if err == nil {
		t.Fatal("Expected error for empty resource type but got none")
	}
}

func TestGetResourceName_InvalidResourceType(t *testing.T) {
	resourceType := "invalid_resource_type"
	separator := "-"
	prefixes := []string{"prefix"}
	name := "testname"
	suffixes := []string{"suffix"}
	randomSuffix := "random"
	convention := ConventionCafClassic
	cleanInput := true
	passthrough := false
	useSlug := true
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	_, err := getResourceName(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, namePrecedence)
	
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
}

func TestGetNameResult_InvalidResourceType(t *testing.T) {
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
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("prefixes", []interface{}{"prefix"})
	d.Set("suffixes", []interface{}{"suffix"})
	d.Set("separator", "-")
	d.Set("resource_type", "invalid_resource_type")
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
}

func TestGetResult_InvalidResourceType(t *testing.T) {
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
	d.Set("resource_type", "invalid_resource_type")
	d.Set("max_length", 63)

	err := getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource type but got none")
	}
}

func TestGetResult_MaxLengthLessThanDefault(t *testing.T) {
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
	d.Set("max_length", 10) // Less than default max length

	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if len(result) > 10 {
		t.Fatalf("Expected result length <= 10, got %d", len(result))
	}
}

func TestGetResult_EmptyName(t *testing.T) {
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
	d.Set("name", "")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionCafRandom)
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
}
