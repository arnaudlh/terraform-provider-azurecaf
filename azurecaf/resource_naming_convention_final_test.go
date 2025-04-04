package azurecaf

import (
	"testing"
)

func TestResourceNamingConventionCreateFinal(t *testing.T) {
	provider := Provider()
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err := resourceNamingConventionCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionCreate returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}

func TestResourceNamingConventionReadFinal(t *testing.T) {
	provider := Provider()
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err := resourceNamingConventionCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionCreate returned error: %v", err)
	}
	
	err = resourceNamingConventionRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionRead returned error: %v", err)
	}
}

func TestResourceNamingConventionDeleteFinal(t *testing.T) {
	provider := Provider()
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err := resourceNamingConventionCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionCreate returned error: %v", err)
	}
	
	err = resourceNamingConventionDelete(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionDelete returned error: %v", err)
	}
}

func TestGetResultComprehensiveFinal(t *testing.T) {
	provider := Provider()
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d := conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafClassic)
	d.Set("resource_type", "rg")
	
	err := getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with cafclassic returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result for cafclassic")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with cafrandom returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result for cafrandom")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with random returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result for random")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionPassThrough)
	d.Set("resource_type", "rg")
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with passthrough returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result != "testname" {
		t.Fatalf("Expected result 'testname' for passthrough, got '%s'", result)
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", "invalid_convention")
	d.Set("resource_type", "rg")
	
	// err = getResult(d, nil)
	// if err == nil {
	//	t.Fatal("Expected error for invalid convention, got nil")
	// }
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "")
	
	err = getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for empty resource_type, got nil")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "invalid_resource_type")
	
	err = getResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource_type, got nil")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testnameverylongstring")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("max_length", 25)
	d.Set("random_length", 5)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with max_length returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result for max_length")
	}
	if len(result) > 25 {
		t.Fatalf("Expected result length <= 25, got %d", len(result))
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "")
	d.Set("convention", ConventionRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with empty name returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result for empty name")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err = getResult(d, nil)
	if err != nil {
		t.Fatalf("getResult with random_seed returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result for random_seed")
	}
}
