package azurecaf

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceNamingConventionCreate(t *testing.T) {
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
	d.Set("resource_type", "afw")
	d.Set("max_length", 63)

	err := resourceNamingConventionCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionCreate returned unexpected error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result after create")
	}
}

func TestResourceNamingConventionRead(t *testing.T) {
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
	d.Set("resource_type", "afw")
	d.Set("max_length", 63)

	err := resourceNamingConventionRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionRead returned unexpected error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result after read")
	}
}

func TestResourceNamingConventionDeleteCore(t *testing.T) {
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
	d.Set("resource_type", "afw")
	d.Set("max_length", 63)

	err := resourceNamingConventionDelete(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionDelete returned unexpected error: %v", err)
	}
}

func TestGetResult_AllConventions(t *testing.T) {
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
			
			if convention != ConventionPassThrough {
				d.Set("prefix", "prefix")
				d.Set("postfix", "postfix")
			}
			
			d.Set("convention", convention)
			d.Set("resource_type", "afw")
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
				if result != "prefix-afw-testname-postfix" {
					t.Errorf("CafClassic should return prefix-afw-name-postfix, got %s", result)
				}
			case ConventionCafRandom, ConventionRandom:
				if len(result) <= len("prefix-afw-testname-postfix") {
					t.Errorf("Random conventions should include random characters, got %s", result)
				}
			}
		})
	}
}

func TestGetResult_EdgeCases(t *testing.T) {
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
			resourceType: "afw",
			maxLength:    63,
			expectError:  false,
		},
		{
			name:         "Very long name",
			resourceName: "thisisaverylongnamethatwillbetrimmedaccordingtothemaxlengthparameter",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafClassic,
			resourceType: "afw",
			maxLength:    20,
			expectError:  false,
		},
		{
			name:         "Special characters in name",
			resourceName: "test!@#$%^&*()name",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafClassic,
			resourceType: "afw",
			maxLength:    63,
			expectError:  false,
		},
		{
			name:         "Storage account with uppercase",
			resourceName: "TESTNAME",
			prefix:       "PREFIX",
			postfix:      "POSTFIX",
			convention:   ConventionCafClassic,
			resourceType: "st",
			maxLength:    24,
			expectError:  false,
		},
		{
			name:         "Invalid resource type",
			resourceName: "testname",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafClassic,
			resourceType: "invalid_resource_type",
			maxLength:    63,
			expectError:  true,
		},
		{
			name:         "Empty name with ConventionRandom",
			resourceName: "",
			prefix:       "",
			postfix:      "",
			convention:   ConventionRandom,
			resourceType: "afw",
			maxLength:    63,
			expectError:  false,
		},
		{
			name:         "Name with max length exactly matching resource max length",
			resourceName: "testname",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafRandom,
			resourceType: "afw",
			maxLength:    0, // Use resource default max length
			expectError:  false,
		},
		{
			name:         "Name with random suffix that needs last char replacement",
			resourceName: "test",
			prefix:       "prefix",
			postfix:      "postfix",
			convention:   ConventionCafRandom,
			resourceType: "afw",
			maxLength:    63,
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

				if tc.maxLength > 0 && len(result) > tc.maxLength {
					t.Errorf("Expected result length <= %d, got %d", tc.maxLength, len(result))
				}

				if tc.resourceType == "st" {
					if result != strings.ToLower(result) {
						t.Errorf("Expected lowercase result for storage account, got %s", result)
					}
				}
				
				if d.Id() == "" {
					t.Error("Expected ID to be set")
				}
				
				if tc.convention == ConventionRandom {
					if strings.Contains(result, tc.resourceName) && tc.resourceName != "" {
						t.Errorf("Random convention should not include original name, got %s", result)
					}
				}
				
				if tc.convention == ConventionPassThrough {
					if result != tc.resourceName {
						t.Errorf("PassThrough convention should preserve original name, expected %s, got %s", tc.resourceName, result)
					}
				}
			}
		})
	}
}
