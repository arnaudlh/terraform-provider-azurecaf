package azurecaf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceNameReadCoverage5(t *testing.T) {
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
	d.Set("result", "prefix-testname-suffix-12345")

	err := resourceNameRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNameRead returned error: %v", err)
	}
}

func TestResourceNamingConventionReadCoverage5(t *testing.T) {
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
	d.Set("result", "prefix-testname-postfix")

	err := resourceNamingConventionRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionRead returned error: %v", err)
	}
}

func TestDataNameReadCoverage5(t *testing.T) {
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
	d.Set("result", "prefix-testname-suffix-12345")

	diags := dataNameRead(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("dataNameRead returned error: %v", diags)
	}
}

func TestGetNameReadResultCoverage5(t *testing.T) {
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
	d.Set("resource_types", []interface{}{"azurerm_storage_account"})
	d.Set("clean_input", true)
	d.Set("passthrough", false)
	d.Set("use_slug", true)
	d.Set("random_length", 5)
	d.Set("random_seed", 123)

	err := getNameReadResult(d, nil)
	if err != nil {
		t.Fatalf("getNameReadResult returned error: %v", err)
	}
}

func TestDataEnvironmentVariableReadCoverage5(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"environment_variable": {
				Type: schema.TypeString,
			},
			"value": {
				Type: schema.TypeString,
			},
			"fails_if_empty": {
				Type: schema.TypeBool,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("environment_variable", "PATH")
	d.Set("fails_if_empty", false)
	d.Set("value", "/usr/local/bin:/usr/bin")

	diags := resourceAction(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("resourceAction returned error: %v", diags)
	}
}

func TestCleanInputCoverage5(t *testing.T) {
	name := "test!@#$%^&*()_+name"
	resource := ResourceDefinitions["azurerm_resource_group"]
	
	result := cleanString(name, &resource)
	
	if result == name {
		t.Fatalf("Expected cleaned string, got '%s'", result)
	}
}

func TestCleanStringCoverage5(t *testing.T) {
	name := "test!@#$%^&*()_+name"
	resource := ResourceDefinitions["azurerm_resource_group"]
	
	result := cleanString(name, &resource)
	
	if result == name {
		t.Fatalf("Expected cleaned string, got '%s'", result)
	}
}

func TestConvertInterfaceToStringCoverage5(t *testing.T) {
	source := []interface{}{"test1", 2, true, 3.14}
	
	result := convertInterfaceToString(source)
	
	if len(result) != len(source) {
		t.Fatalf("Expected length %d, got %d", len(source), len(result))
	}
}

func TestRandSeqCoverage5(t *testing.T) {
	length := 10
	
	result := randSeq(length, nil)
	
	if len(result) != length {
		t.Fatalf("Expected length %d, got %d", length, len(result))
	}
}
