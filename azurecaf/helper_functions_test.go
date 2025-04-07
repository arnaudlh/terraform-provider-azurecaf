package azurecaf

import (
	"testing"
)

func TestHelperFunctionsUnique(t *testing.T) {
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
}

func TestGetResourceNameComprehensiveUnique(t *testing.T) {
	result, err := getResourceName("azurerm_resource_group", "-", []string{"prefix1", "prefix2"}, "testname", []string{"suffix1", "suffix2"}, "slug", "cafrandom", true, false, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with all parameters returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	result, err = getResourceName("azurerm_resource_group", "-", []string{}, "testname", []string{}, "slug", "cafrandom", true, false, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with empty prefixes and suffixes returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	result, err = getResourceName("azurerm_resource_group", "_", []string{"prefix"}, "testname", []string{"suffix"}, "slug", "cafrandom", true, false, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with different separator returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	result, err = getResourceName("azurerm_resource_group", "-", []string{"prefix"}, "testname", []string{"suffix"}, "slug", "cafrandom", true, false, true, []string{"suffixes", "random", "slug", "name", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with different order returned error: %v", err)
	}
	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	_, err = getResourceName("azurerm_resource_group", "-", []string{"prefix"}, "testname", []string{"suffix"}, "slug", "cafrandom", true, false, true, []string{})
	if err == nil {
		t.Fatal("Expected error for empty order array, got nil")
	}

	result, err = getResourceName("azurerm_resource_group", "-", []string{"prefix"}, "testname", []string{"suffix"}, "slug", "cafrandom", true, true, true, []string{"prefixes", "name", "slug", "random", "suffixes"})
	if err != nil {
		t.Fatalf("getResourceName with passthrough returned error: %v", err)
	}
	if result != "testname" {
		t.Fatalf("Expected 'testname' with passthrough, got '%s'", result)
	}

	result, err = getResourceName("azurerm_storage_account", "-", []string{}, "TestName", []string{}, "", "cafrandom", true, false, true, []string{"name", "slug", "random", "suffixes", "prefixes"})
	if err != nil {
		t.Fatalf("getResourceName with lowercase conversion returned error: %v", err)
	}
	if result != result {
		t.Fatalf("Expected lowercase result, got '%s'", result)
	}
}

func TestComposeNameComprehensiveUnique(t *testing.T) {
	result := composeName("-",
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
		[]string{"prefix1", "prefix2"},
		"testname",
		"slug",
		[]string{"suffix1", "suffix2"},
		"random",
		100,
		[]string{"suffixes", "random", "slug", "name", "prefixes"})

	if result == "" {
		t.Fatal("Expected non-empty result")
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
