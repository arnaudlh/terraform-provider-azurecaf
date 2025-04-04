package azurecaf

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceActionCoverage(t *testing.T) {
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

	t.Skip("Skipping test that requires environment variables")
}

func TestResourceNameDeleteCoverage(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"result": {
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("result", "test-result")

	err := resourceNameDelete(d, nil)
	if err != nil {
		t.Fatalf("resourceNameDelete returned error: %v", err)
	}
}

func TestResourceNameReadCoverage2(t *testing.T) {
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

	err := resourceNameRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNameRead returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}

func TestResourceNamingConventionDeleteCoverage(t *testing.T) {
	r := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"result": {
				Type: schema.TypeString,
			},
		},
	}

	d := r.TestResourceData()
	d.Set("name", "testname")
	d.Set("result", "test-result")

	err := resourceNamingConventionDelete(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionDelete returned error: %v", err)
	}
}

func TestResourceNamingConventionReadCoverage2(t *testing.T) {
	t.Skip("Skipping test that requires more complex setup")
}

func TestGetSlugCoverage(t *testing.T) {
	resourceType := "azurerm_resource_group"
	convention := ConventionCafClassic
	slug := getSlug(resourceType, convention)
	
	if slug != "rg" {
		t.Fatalf("Expected slug 'rg' for resource type '%s', got '%s'", resourceType, slug)
	}
	
	resourceType = "nonexistent_resource"
	slug = getSlug(resourceType, convention)
	
	if slug != "" {
		t.Fatalf("Expected empty slug for nonexistent resource type, got '%s'", slug)
	}
}

func TestComposeNameCoverage(t *testing.T) {
	separator := "-"
	prefixes := []string{"prefix1", "prefix2"}
	name := "testname"
	slug := "rg"
	suffixes := []string{"suffix1", "suffix2"}
	randomSuffix := "random"
	maxlength := 63
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result := composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	separator = "-"
	prefixes = []string{}
	name = ""
	slug = ""
	suffixes = []string{}
	randomSuffix = ""
	maxlength = 63
	namePrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}
	
	result = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, maxlength, namePrecedence)
	
	if result != "" {
		t.Fatalf("Expected empty result, got '%s'", result)
	}
}

func TestCleanSliceCoverage(t *testing.T) {
	data := []string{"test1", "test2", "test3"}
	resource := ResourceDefinitions["azurerm_resource_group"]
	result := cleanSlice(data, &resource)
	
	if len(result) != len(data) {
		t.Fatalf("Expected length %d, got %d", len(data), len(result))
	}
	
	for i := range data {
		if result[i] != data[i] {
			t.Fatalf("Expected '%s' at index %d, got '%s'", data[i], i, result[i])
		}
	}
	
	data = []string{"test1!", "test2@", "test3#"}
	resource = ResourceDefinitions["azurerm_resource_group"]
	result = cleanSlice(data, &resource)
	
	if len(result) != len(data) {
		t.Fatalf("Expected length %d, got %d", len(data), len(result))
	}
	
	for i := range data {
		if result[i] == data[i] {
			t.Fatalf("Expected cleaned string at index %d, got '%s'", i, result[i])
		}
	}
}

func TestConcatenateParametersCoverage(t *testing.T) {
	separator := "-"
	prefixes := []string{"prefix1", "prefix2"}
	name := []string{"testname"}
	suffixes := []string{"suffix1", "suffix2"}
	
	result := concatenateParameters(separator, prefixes, name, suffixes)
	
	if !strings.Contains(result, "prefix1") || !strings.Contains(result, "prefix2") || 
	   !strings.Contains(result, "testname") || !strings.Contains(result, "suffix1") || 
	   !strings.Contains(result, "suffix2") {
		t.Fatalf("Expected result to contain all elements, got '%s'", result)
	}
	
	separator = "-"
	prefixes = []string{}
	name = []string{"testname"}
	suffixes = []string{}
	
	result = concatenateParameters(separator, prefixes, name, suffixes)
	expected := "testname"
	
	if result != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, result)
	}
	
	separator = "-"
	prefixes = []string{"prefix1", "prefix2"}
	name = []string{}
	suffixes = []string{"suffix1", "suffix2"}
	
	result = concatenateParameters(separator, prefixes, name, suffixes)
	
	if !strings.Contains(result, "prefix1") || !strings.Contains(result, "prefix2") || 
	   !strings.Contains(result, "suffix1") || !strings.Contains(result, "suffix2") {
		t.Fatalf("Expected result to contain all elements, got '%s'", result)
	}
}
