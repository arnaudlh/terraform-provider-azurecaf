package schemas

import (
	models "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

func getResourceMaps() []string {
	resourceMapsKeys := make([]string, 0, len(models.ResourceDefinitions))
	for k := range models.ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}
	return resourceMapsKeys
}
