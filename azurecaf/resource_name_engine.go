package azurecaf

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/models"
)

func cleanSlice(names []string, resourceDefinition *models.ResourceStructure) []string {
	result := make([]string, len(names))
	for i, name := range names {
		if name == "" {
			result[i] = name
			continue
		}
		cleaned := cleanString(name, resourceDefinition)
		if cleaned == "" {
			result[i] = name
		} else {
			result[i] = cleaned
		}
	}
	return result
}

func cleanString(name string, resourceDefinition *models.ResourceStructure) string {
	if name == "" || resourceDefinition == nil || resourceDefinition.RegEx == "" {
		return name
	}

	// Special handling for Azure Container Registry which doesn't allow hyphens
	if strings.Contains(resourceDefinition.RegEx, "^[a-zA-Z0-9]{1,63}$") {
		return strings.ReplaceAll(name, "-", "")
	}

	pattern := resourceDefinition.RegEx
	start := strings.Index(pattern, "[")
	end := strings.Index(pattern, "]")
	if start == -1 || end == -1 || start >= end {
		return name
	}

	allowedChars := pattern[start+1 : end]
	invalidCharsPattern := fmt.Sprintf("[^%s]", allowedChars)
	re, err := regexp.Compile(invalidCharsPattern)
	if err != nil {
		return name
	}

	cleaned := re.ReplaceAllString(name, "")
	if cleaned == "" {
		return name
	}
	return cleaned
}



func getResource(resourceType string) (*models.ResourceStructure, error) {
	if resourceKey, existing := models.ResourceMaps[resourceType]; existing {
		resourceType = resourceKey
	}
	if resource, resourceFound := models.ResourceDefinitions[resourceType]; resourceFound {
		return &resource, nil
	}
	return nil, fmt.Errorf("invalid resource type %s", resourceType)
}

// Retrieve the resource slug / shortname based on the resourceType and the selected convention
func getSlug(resourceType string) string {
	if val, ok := models.ResourceDefinitions[resourceType]; ok {
		return val.CafPrefix
	}
	return ""
}

func trimResourceName(resourceName string, maxLength int) string {
	var length int = len(resourceName)

	if length > maxLength {
		length = maxLength
	}

	return string(resourceName[0:length])
}

func composeName(separator string,
	prefixes []string,
	name string,
	slug string,
	suffixes []string,
	randomSuffix string,
	maxlength int,
	namePrecedence []string) string {
	contents := []string{}
	currentlength := 0

	for _, precedence := range namePrecedence {
		initialized := 0
		if len(contents) > 0 {
			initialized = len(separator)
		}

		switch precedence {
		case "prefixes":
			for _, prefix := range prefixes {
				if len(prefix) > 0 && currentlength+len(prefix)+initialized <= maxlength {
					contents = append(contents, prefix)
					currentlength = currentlength + len(prefix) + initialized
				}
			}
		case "name":
			if len(name) > 0 && currentlength+len(name)+initialized <= maxlength {
				contents = append(contents, name)
				currentlength = currentlength + len(name) + initialized
			}
		case "slug":
			if len(slug) > 0 && currentlength+len(slug)+initialized <= maxlength {
				contents = append(contents, slug)
				currentlength = currentlength + len(slug) + initialized
			}
		case "random":
			if len(randomSuffix) > 0 && currentlength+len(randomSuffix)+initialized <= maxlength {
				contents = append(contents, randomSuffix)
				currentlength = currentlength + len(randomSuffix) + initialized
			}
		case "suffixes":
			for _, suffix := range suffixes {
				if len(suffix) > 0 && currentlength+len(suffix)+initialized <= maxlength {
					contents = append(contents, suffix)
					currentlength = currentlength + len(suffix) + initialized
				}
			}
		}
	}

	content := strings.Join(contents, separator)
	return content
}

// validateResourceType is implemented in data_name.go

func getResourceName(resourceTypeName string, separator string,
	prefixes []string,
	name string,
	suffixes []string,
	randomSuffix string,
	cleanInput bool,
	passthrough bool,
	useSlug bool,
	namePrecedence []string) (string, error) {

	resource, err := getResource(resourceTypeName)
	if err != nil {
		return "", err
	}
	validationRegEx, err := regexp.Compile(resource.ValidationRegExp)
	if err != nil {
		return "", err
	}

	slug := ""
	if useSlug {
		slug = getSlug(resourceTypeName)
	}

	// Set default name precedence if not provided
	if len(namePrecedence) == 0 {
		namePrecedence = []string{"prefixes", "slug", "name", "random", "suffixes"}
	}

	if cleanInput {
		prefixes = cleanSlice(prefixes, resource)
		suffixes = cleanSlice(suffixes, resource)
		name = cleanString(name, resource)
		separator = cleanString(separator, resource)
		randomSuffix = cleanString(randomSuffix, resource)
	}

	var resourceName string

	if passthrough {
		resourceName = name
	} else {
		resourceName = composeName(separator, prefixes, name, slug, suffixes, randomSuffix, resource.MaxLength, namePrecedence)
	}
	resourceName = trimResourceName(resourceName, resource.MaxLength)

	// Handle resources that require alphanumeric-only names (no hyphens)
	if strings.Contains(resource.ValidationRegExp, "^[a-zA-Z0-9]") && !strings.Contains(resource.ValidationRegExp, "-") {
		resourceName = strings.ReplaceAll(resourceName, "-", "")
	}

	if resource.LowerCase {
		resourceName = strings.ToLower(resourceName)
	}

	if !validationRegEx.MatchString(resourceName) {
		return "", fmt.Errorf("invalid name for CAF naming %s %s, the pattern %s doesn't match %s", resource.ResourceTypeName, name, resource.ValidationRegExp, resourceName)
	}

	return resourceName, nil
}
