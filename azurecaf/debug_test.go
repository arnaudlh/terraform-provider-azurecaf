package azurecaf

import (
	"fmt"
	"testing"
)

func TestDebugGetResource(t *testing.T) {
	resource, err := getResource("azurerm_resource_group")
	if err != nil {
		t.Fatalf("getResource returned error: %v", err)
	}
	fmt.Printf("ResourceTypeName: %s\n", resource.ResourceTypeName)
	fmt.Printf("CafPrefix: %s\n", resource.CafPrefix)
	fmt.Printf("MinLength: %d\n", resource.MinLength)
	fmt.Printf("MaxLength: %d\n", resource.MaxLength)
}
