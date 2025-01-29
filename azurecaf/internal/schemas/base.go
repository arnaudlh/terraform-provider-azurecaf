package schemas

import (
	"github.com/arnaudlh/terraform-provider-azurecaf/azurecaf/internal/models"
)

func GetResourceMaps() []string {
	resourceMapsKeys := make([]string, 0, len(models.ResourceDefinitions))
	for k := range models.ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}
	return resourceMapsKeys
}
