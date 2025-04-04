package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetResultConventions(t *testing.T) {
	conventions := []string{
		ConventionCafClassic,
		ConventionCafRandom,
		ConventionRandom,
		ConventionPassThrough,
	}

	for _, convention := range conventions {
		t.Run(convention, func(t *testing.T) {
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
			d.Set("convention", convention)
			d.Set("resource_type", "azurerm_resource_group")
			d.Set("max_length", 63)

			err := getResult(d, nil)
			if err != nil {
				t.Fatalf("getResult returned unexpected error: %v", err)
			}

			result := d.Get("result").(string)
			if result == "" {
				t.Fatal("Expected non-empty result")
			}

			switch convention {
			case ConventionPassThrough:
				if result != "testname" {
					t.Errorf("PassThrough should return original name, got %s", result)
				}
			case ConventionCafClassic:
				if result != "prefix-rg-testname-postfix" {
					t.Errorf("CafClassic should return prefix-rg-name-postfix, got %s", result)
				}
			case ConventionCafRandom, ConventionRandom:
				if len(result) <= len("prefix-rg-testname-postfix") {
					t.Errorf("Random conventions should include random characters, got %s", result)
				}
			}
		})
	}
}

func TestGetResultEdgeCases2(t *testing.T) {
	testCases := []struct {
		name         string
		resourceName string
		prefix       string
		postfix      string
		convention   string
		resourceType string
		maxLength    int
		expectError  bool
	}{
		{
			name:         "Empty name",
			resourceName: "",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafRandom,
			resourceType: "azurerm_resource_group",
			maxLength:    63,
			expectError:  false,
		},
		{
			name:         "Very long name",
			resourceName: "thisisaverylongnamethatwillbetrimmedaccordingtothemaxlengthparameter",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafClassic,
			resourceType: "azurerm_resource_group",
			maxLength:    20,
			expectError:  false,
		},
		{
			name:         "Special characters in name",
			resourceName: "test!@#$%^&*()name",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafClassic,
			resourceType: "azurerm_resource_group",
			maxLength:    63,
			expectError:  false,
		},
		{
			name:         "Storage account with uppercase",
			resourceName: "TESTNAME",
			prefix:       "PREFIX",
			postfix:      "POSTFIX",
			convention:   ConventionCafClassic,
			resourceType: "azurerm_storage_account",
			maxLength:    24,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
			d.Set("name", tc.resourceName)
			d.Set("prefix", tc.prefix)
			d.Set("postfix", tc.postfix)
			d.Set("convention", tc.convention)
			d.Set("resource_type", tc.resourceType)
			d.Set("max_length", tc.maxLength)

			err := getResult(d, nil)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				result := d.Get("result").(string)
				if result == "" {
					t.Error("Expected non-empty result")
				}

				if len(result) > tc.maxLength {
					t.Errorf("Expected result length <= %d, got %d", tc.maxLength, len(result))
				}

				if tc.resourceType == "azurerm_storage_account" {
					if result != result[:len(result)] {
						t.Errorf("Expected lowercase result for storage account, got %s", result)
					}
				}
			}
		})
	}
}

func TestResourceNamingConventionCreateDelete(t *testing.T) {
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

	err := resourceNamingConventionCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionCreate returned unexpected error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result after create")
	}

	err = resourceNamingConventionDelete(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionDelete returned unexpected error: %v", err)
	}
}

func TestResourceNamingConventionResource(t *testing.T) {
	resource := resourceNamingConvention()
	
	if resource.Schema["name"] == nil {
		t.Fatal("Expected name field in schema")
	}
	if resource.Schema["convention"] == nil {
		t.Fatal("Expected convention field in schema")
	}
	if resource.Schema["prefix"] == nil {
		t.Fatal("Expected prefix field in schema")
	}
	if resource.Schema["postfix"] == nil {
		t.Fatal("Expected postfix field in schema")
	}
	if resource.Schema["resource_type"] == nil {
		t.Fatal("Expected resource_type field in schema")
	}
	if resource.Schema["max_length"] == nil {
		t.Fatal("Expected max_length field in schema")
	}
	if resource.Schema["result"] == nil {
		t.Fatal("Expected result field in schema")
	}
	
	if resource.Create == nil {
		t.Fatal("Expected Create function")
	}
	if resource.Read == nil {
		t.Fatal("Expected Read function")
	}
	if resource.Delete == nil {
		t.Fatal("Expected Delete function")
	}
}
