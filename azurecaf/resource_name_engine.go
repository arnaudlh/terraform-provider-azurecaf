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

	for i := 0; i < len(namePrecedence); i++ {
		initialized := 0
		if len(contents) > 0 {
			initialized = len(separator)
		}
		switch c := namePrecedence[i]; c {
		case "name":
			if len(name) > 0 {
				if currentlength+len(name)+initialized <= maxlength {
					contents = append(contents, name)
					currentlength = currentlength + len(name) + initialized
				}
			}
		case "slug":
			if len(slug) > 0 {
				if currentlength+len(slug)+initialized <= maxlength {
					contents = append([]string{slug}, contents...)
					currentlength = currentlength + len(slug) + initialized
				}
			}
		case "random":
			if len(randomSuffix) > 0 {
				if currentlength+len(randomSuffix)+initialized <= maxlength {
					contents = append(contents, randomSuffix)
					currentlength = currentlength + len(randomSuffix) + initialized
				}
			}
		case "suffixes":
			if len(suffixes) > 0 {
				if len(suffixes[0]) > 0 {
					if currentlength+len(suffixes[0])+initialized <= maxlength {
						contents = append(contents, suffixes[0])
						currentlength = currentlength + len(suffixes[0]) + initialized
					}
				}
				suffixes = suffixes[1:]
				if len(suffixes) > 0 {
					i--
				}
			}
		case "prefixes":
			if len(prefixes) > 0 {
				if len(prefixes[len(prefixes)-1]) > 0 {
					if currentlength+len(prefixes[len(prefixes)-1])+initialized <= maxlength {
						contents = append([]string{prefixes[len(prefixes)-1]}, contents...)
						currentlength = currentlength + len(prefixes[len(prefixes)-1]) + initialized
					}
				}
				prefixes = prefixes[:len(prefixes)-1]
				if len(prefixes) > 0 {
					i--
				}
			}

		}

	}
	content := strings.Join(contents, separator)
	return content
}

func validateResourceTypes(resourceType string, resourceTypes []string) error {
	isEmpty := len(resourceType) == 0 && len(resourceTypes) == 0
	if isEmpty {
		return fmt.Errorf("resource_type and resource_types parameters are empty, you must specify at least one resource type")
	}
	errorStrings := []string{}
	resourceList := resourceTypes
	if len(resourceType) > 0 {
		resourceList = append(resourceList, resourceType)
	}

	for _, resource := range resourceList {
		_, err := getResource(resource)
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}
	if len(errorStrings) > 0 {
		return fmt.Errorf("%s", strings.Join(errorStrings, "\n"))
	}
	return nil
}

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

	if resource.LowerCase {
		resourceName = strings.ToLower(resourceName)
	}

	if !validationRegEx.MatchString(resourceName) {
		return "", fmt.Errorf("invalid name for CAF naming %s %s, the pattern %s doesn't match %s", resource.ResourceTypeName, name, resource.ValidationRegExp, resourceName)
	}

	return resourceName, nil
}
