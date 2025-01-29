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
	
	// Special case: if name is longer than maxlength, just return it truncated
	if len(name) >= maxlength {
		return name[:maxlength]
	}

	var contents []string

	// Helper to calculate total length with separators
	calcTotalLength := func(components []string) int {
		if len(components) == 0 {
			return 0
		}
		total := 0
		for _, c := range components {
			if len(c) > 0 {
				total += len(c)
			}
		}
		return total + (len(components)-1)*len(separator)
	}

	// Helper to add a component
	addComponent := func(component string) bool {
		if len(component) == 0 {
			return true
		}
		newComponents := append([]string{}, contents...)
		newComponents = append(newComponents, component)
		if calcTotalLength(newComponents) > maxlength {
			return false
		}
		contents = newComponents
		return true
	}

	// Process components in order specified by namePrecedence
	for _, component := range namePrecedence {
		switch component {
		case "prefixes":
			// Add prefixes in order
			for _, prefix := range prefixes {
				if len(prefix) > 0 && calcTotalLength(append(contents, prefix)) <= maxlength {
					contents = append(contents, prefix)
				}
			}
		case "slug":
			// Add slug if present and fits
			if len(slug) > 0 && calcTotalLength(append(contents, slug)) <= maxlength {
				contents = append(contents, slug)
			}
		case "name":
			// Add name if present and fits
			if len(name) > 0 && calcTotalLength(append(contents, name)) <= maxlength {
				contents = append(contents, name)
			}
		case "random":
			// Add random suffix if present and fits
			if len(randomSuffix) > 0 && calcTotalLength(append(contents, randomSuffix)) <= maxlength {
				contents = append(contents, randomSuffix)
			}
		case "suffixes":
			// Add suffixes in order
			for _, suffix := range suffixes {
				if len(suffix) > 0 && calcTotalLength(append(contents, suffix)) <= maxlength {
					contents = append(contents, suffix)
				}
			}
		}
	}

	// Join all components and ensure max length
	result := strings.Join(contents, separator)
	if len(result) > maxlength {
		result = result[:maxlength]
	}
	return result
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

	// Always use the standard name precedence order for consistency
	namePrecedence = []string{"prefixes", "slug", "name", "random", "suffixes"}

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
