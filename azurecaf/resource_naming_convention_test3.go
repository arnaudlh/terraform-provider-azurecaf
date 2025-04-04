package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResultConventionCases(t *testing.T) {
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
	d.Set("name", "")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 0)

	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error for ConventionRandom and empty name: %v", err)
	}

	d = r.TestResourceData()
	d.Set("name", "short")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error for ConventionCafRandom and short name: %v", err)
	}

	d = r.TestResourceData()
	d.Set("name", "thisisaverylongnamethatwillbecutoffbythemaxlengthparameter")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 20)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error for ConventionCafRandom and long name: %v", err)
	}

	result := d.Get("result").(string)
	if len(result) > 20 {
		t.Fatalf("Expected result length to be <= 20, got %d", len(result))
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
		t.Fatalf("getResult returned error for ConventionPassThrough: %v", err)
	}

	d = r.TestResourceData()
	d.Set("name", "UPPERCASE")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_storage_account") // Storage accounts require lowercase
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error for lowercase resource: %v", err)
	}

	result = d.Get("result").(string)
	for _, c := range result {
		if c >= 'A' && c <= 'Z' {
			t.Fatalf("Expected lowercase result for storage account, got %s", result)
		}
	}

	d = r.TestResourceData()
	d.Set("name", "exactlength")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 11) // "exactlength" is 11 characters

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error for name length exactly at max length: %v", err)
	}

	result = d.Get("result").(string)
	if len(result) != 11 {
		t.Fatalf("Expected result length to be 11, got %d", len(result))
	}

	d = r.TestResourceData()
	d.Set("name", "test!@#$%^&*()_+name")
	d.Set("prefix", "")
	d.Set("postfix", "")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("max_length", 63)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error for name with special characters: %v", err)
	}
}
