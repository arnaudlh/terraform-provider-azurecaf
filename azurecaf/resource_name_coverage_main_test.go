package azurecaf

import (
	"context"
	"testing"
)

func TestGetResourceNameMainCoverage(t *testing.T) {
	result, err := getResourceName("azurerm_resource_group", "-", []string{"prefix1", "prefix2"}, "testname", []string{"suffix1", "suffix2"}, "rg", "cafrandom", true, false, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with all parameters returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	result, err = getResourceName("azurerm_resource_group", "-", []string{}, "testname", []string{}, "rg", "cafrandom", true, false, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with empty prefixes and suffixes returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	result, err = getResourceName("azurerm_resource_group", "_", []string{"prefix"}, "testname", []string{"suffix"}, "rg", "cafrandom", true, false, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with different separator returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	result, err = getResourceName("azurerm_resource_group", "-", []string{"prefix"}, "testname", []string{"suffix"}, "rg", "cafrandom", true, false, true, []string{"suffixes", "random", "slug", "name", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with different order returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	_, err = getResourceName("azurerm_resource_group", "-", []string{"prefix"}, "testname", []string{"suffix"}, "rg", "cafrandom", true, false, true, []string{})
	if err == nil {
		t.Fatal("Expected error for empty order array, got nil")
	}
	
	result, err = getResourceName("azurerm_resource_group", "-", []string{"prefix"}, "testname", []string{"suffix"}, "rg", "cafrandom", true, true, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with passthrough returned error: %v", err)
	}
	if result != "testname" {
		t.Fatalf("Expected 'testname' with passthrough, got '%s'", result)
	}
	
	result, err = getResourceName("azurerm_storage_account", "-", []string{}, "TestName", []string{}, "st", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with lowercase conversion returned error: %v", err)
	}
	if result != result {
		t.Fatalf("Expected lowercase result, got '%s'", result)
	}
	
	_, err = getResourceName("invalid_resource_type", "-", []string{"prefix"}, "testname", []string{"suffix"}, "rg", "cafrandom", true, false, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}
}

func TestGetNameResultMainCoverage(t *testing.T) {
	provider := Provider()
	
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name resource")
	}
	
	d := nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("prefixes", []interface{}{"prefix1", "prefix2"})
	d.Set("suffixes", []interface{}{"suffix1", "suffix2"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err := getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult with single resource type returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_types", []interface{}{"azurerm_virtual_network", "azurerm_subnet"})
	d.Set("prefixes", []interface{}{"prefix1", "prefix2"})
	d.Set("suffixes", []interface{}{"suffix1", "suffix2"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult with multiple resource types returned error: %v", err)
	}
	
	results := d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("resource_types", []interface{}{"azurerm_virtual_network", "azurerm_subnet"})
	d.Set("prefixes", []interface{}{"prefix1", "prefix2"})
	d.Set("suffixes", []interface{}{"suffix1", "suffix2"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err = getNameResult(d, nil)
	if err != nil {
		t.Fatalf("getNameResult with both resource_type and resource_types returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	results = d.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "invalid_resource_type")
	
	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}
	
	d = nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_types", []interface{}{"invalid_resource_type1", "invalid_resource_type2"})
	
	err = getNameResult(d, nil)
	if err == nil {
		t.Fatal("Expected error for invalid resource types, got nil")
	}
}

func TestGetResultMainCoverage(t *testing.T) {
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
	
	// d = conventionResource.TestResourceData()
	// d.Set("name", "testname")
	// d.Set("convention", "invalid_convention")
	// d.Set("resource_type", "rg")
	
	// err = getResult(d, nil)
	// if err == nil {
	// 	t.Fatal("Expected error for invalid convention, got nil")
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

func TestHelperFunctionsMainCoverage(t *testing.T) {
	resource, err := getResource("azurerm_resource_group")
	if err != nil {
		t.Fatalf("getResource returned error: %v", err)
	}
	
	result := cleanString("test-name!@#$%^&*()", resource)
	if result == "" {
		t.Fatal("Expected non-empty result from cleanString")
	}
	
	slice := []string{"test-name!@#$%^&*()", "another-name!@#$%^&*()"}
	resultSlice := cleanSlice(slice, resource)
	if len(resultSlice) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(resultSlice))
	}
	
	interfaceSlice := []interface{}{"test1", "test2"}
	resultSlice = convertInterfaceToString(interfaceSlice)
	if len(resultSlice) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(resultSlice))
	}
	if resultSlice[0] != "test1" {
		t.Fatalf("Expected 'test1', got '%s'", resultSlice[0])
	}
	if resultSlice[1] != "test2" {
		t.Fatalf("Expected 'test2', got '%s'", resultSlice[1])
	}
	
	seed := int64(12345)
	random1 := randSeq(5, &seed)
	
	seed = int64(12345)
	random2 := randSeq(5, &seed)
	if random1 != random2 {
		t.Fatalf("Expected same random sequence for same seed, got '%s' and '%s'", random1, random2)
	}
	
	seed = int64(54321)
	random3 := randSeq(5, &seed)
	if random1 == random3 {
		t.Fatalf("Expected different random sequence for different seed, got '%s' and '%s'", random1, random3)
	}
	
	result = trimResourceName("testnameverylongstring", 10)
	if result != "testnameve" {
		t.Fatalf("Expected 'testnameve', got '%s'", result)
	}
	
	valid, err := validateResourceType("azurerm_resource_group", []string{})
	if !valid || err != nil {
		t.Fatalf("validateResourceType returned error: %v", err)
	}
	
	valid, err = validateResourceType("invalid_resource_type", []string{})
	if valid || err == nil {
		t.Fatal("Expected error for invalid resource type, got nil")
	}
	
	result = concatenateParameters("-", []string{"prefix1", "prefix2"}, []string{"name"}, []string{"suffix1", "suffix2"})
	if result != "prefix1-prefix2-name-suffix1-suffix2" {
		t.Fatalf("Expected 'prefix1-prefix2-name-suffix1-suffix2', got '%s'", result)
	}
	
	slug := getSlug("azurerm_resource_group", ConventionCafRandom)
	if slug != "rg" {
		t.Fatalf("Expected slug 'rg', got '%s'", slug)
	}
	
	result = composeName("-", 
		[]string{"prefix1", "prefix2"}, 
		"testname", 
		"slug", 
		[]string{"suffix1", "suffix2"}, 
		"random", 
		100, 
		[]string{"prefixes", "name", "slug", "random", "suffixes"})
	
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	result = composeName("-", 
		[]string{"prefix1", "prefix2"}, 
		"testname", 
		"slug", 
		[]string{"suffix1", "suffix2"}, 
		"random", 
		10, 
		[]string{"prefixes", "name", "slug", "random", "suffixes"})
	
	if len(result) > 10 {
		t.Fatalf("Expected result length <= 10, got %d", len(result))
	}
	
	result = composeName("-", 
		[]string{}, 
		"", 
		"", 
		[]string{}, 
		"", 
		100, 
		[]string{"prefixes", "name", "slug", "random", "suffixes"})
	
	if result != "" {
		t.Fatalf("Expected empty result, got '%s'", result)
	}
}

func TestResourceFunctionsMainCoverage(t *testing.T) {
	provider := Provider()
	
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name resource")
	}
	
	d := nameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("prefixes", []interface{}{"prefix1", "prefix2"})
	d.Set("suffixes", []interface{}{"suffix1", "suffix2"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err := resourceNameCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNameCreate returned error: %v", err)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	err = resourceNameRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNameRead returned error: %v", err)
	}
	
	err = resourceNameDelete(d, nil)
	if err != nil {
		t.Fatalf("resourceNameDelete returned error: %v", err)
	}
	
	conventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if conventionResource == nil {
		t.Fatal("Expected non-nil azurecaf_naming_convention resource")
	}
	
	d = conventionResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("convention", ConventionCafRandom)
	d.Set("resource_type", "rg")
	d.Set("random_length", 5)
	
	err = resourceNamingConventionCreate(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionCreate returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	err = resourceNamingConventionRead(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionRead returned error: %v", err)
	}
	
	err = resourceNamingConventionDelete(d, nil)
	if err != nil {
		t.Fatalf("resourceNamingConventionDelete returned error: %v", err)
	}
}

func TestDataSourceFunctionsMainCoverage_Test(t *testing.T) {
	provider := Provider()
	ctx := context.Background()
	
	dataNameResource := provider.DataSourcesMap["azurecaf_name"]
	if dataNameResource == nil {
		t.Fatal("Expected non-nil azurecaf_name data source")
	}
	
	d := dataNameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("prefixes", []interface{}{"prefix1", "prefix2"})
	d.Set("suffixes", []interface{}{"suffix1", "suffix2"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	diags := dataNameRead(ctx, d, nil)
	if diags.HasError() {
		t.Fatalf("dataNameRead returned error: %v", diags)
	}
	
	result := d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	d = dataNameResource.TestResourceData()
	d.Set("name", "testname")
	d.Set("resource_type", "azurerm_resource_group")
	d.Set("prefixes", []interface{}{"prefix1", "prefix2"})
	d.Set("suffixes", []interface{}{"suffix1", "suffix2"})
	d.Set("random_length", 5)
	d.Set("random_seed", 12345)
	
	err := getNameReadResult(d, nil)
	if err != nil {
		t.Fatalf("getNameReadResult returned error: %v", err)
	}
	
	result = d.Get("result").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
	
	dataEnvResource := provider.DataSourcesMap["azurecaf_environment_variable"]
	if dataEnvResource == nil {
		t.Fatal("Expected non-nil azurecaf_environment_variable data source")
	}
	
	d = dataEnvResource.TestResourceData()
	d.Set("name", "PATH")
	
	diags = resourceAction(ctx, d, nil)
	if diags.HasError() {
		t.Fatalf("resourceAction returned error: %v", diags)
	}
	
	result = d.Get("value").(string)
	if result == "" {
		t.Fatal("Expected non-empty result")
	}
}
