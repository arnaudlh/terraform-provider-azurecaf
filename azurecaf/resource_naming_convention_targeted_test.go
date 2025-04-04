package azurecaf

import (
	"testing"
)

func TestGetResultTargetedLines(t *testing.T) {
	provider := Provider()

	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}

	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)

	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	d = conventionResource.TestResourceData()
	d.Set("name", "")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	d = conventionResource.TestResourceData()
	d.Set("name", "testnameverylongstring")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("max_length", 25) // Just enough for the name plus separator

	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult returned error: %v", err)
	}

	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	if len(result) > 25 {
		t.Fatalf("Expected result length <= 25, got %d", len(result))
	}

	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionPassThrough)
	d.Set("resource_type", "invalid_resource_type")

	err = getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}

	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionPassThrough)
	d.Set("resource_type", "")

	err = getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for empty resource type, got nil")
	}
}
