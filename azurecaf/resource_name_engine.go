package azurecaf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

func cleanSlice(names []string, resourceDefinition *models.ResourceStructure) []string {
	for i, name := range names {
		names[i] = cleanString(name, resourceDefinition)
	}
	return names
}

func cleanString(name string, resourceDefinition *models.ResourceStructure) string {
	if name == "" {
		return ""
	}
	if resourceDefinition == nil {
		return name
	}
	
	// First validate the input using ValidationRegExp if present
	if resourceDefinition.ValidationRegExp != "" {
		validationRegex, err := regexp.Compile(resourceDefinition.ValidationRegExp)
		if err == nil && !validationRegex.MatchString(name) {
			return ""
		}
	}
	
	// Then clean the string using RegEx if present
	if resourceDefinition.RegEx != "" {
		cleanRegex, err := regexp.Compile(resourceDefinition.RegEx)
		if err == nil {
			return cleanRegex.ReplaceAllString(name, "")
		}
	}
	
	return name
}

func concatenateParameters(separator string, parameters ...[]string) string {
	elems := []string{}
	for _, items := range parameters {
		for _, item := range items {
			if len(item) > 0 {
				elems = append(elems, []string{item}...)
			}
		}
	}
	return strings.Join(elems, separator)
}

func getResource(resourceType string) (*models.ResourceStructure, error) {
	return models.GetResourceStructure(resourceType)
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
		separatorLen := 0
		if len(contents) > 0 {
			separatorLen = len(separator)
		}
		switch c := namePrecedence[i]; c {
		case "name":
			if len(name) > 0 {
				if currentlength+len(name)+separatorLen <= maxlength {
					contents = append(contents, name)
					currentlength = currentlength + len(name) + separatorLen
				}
			}
		case "slug":
			if len(slug) > 0 {
				if currentlength+len(slug)+separatorLen <= maxlength {
					contents = append([]string{slug}, contents...)
					currentlength = currentlength + len(slug) + separatorLen
				}
			}
		case "random":
			if len(randomSuffix) > 0 {
				if currentlength+len(randomSuffix)+separatorLen <= maxlength {
					contents = append(contents, randomSuffix)
					currentlength = currentlength + len(randomSuffix) + separatorLen
				}
			}
		case "suffixes":
			if len(suffixes) > 0 {
				if len(suffixes[0]) > 0 {
					if currentlength+len(suffixes[0])+separatorLen <= maxlength {
						contents = append(contents, suffixes[0])
						currentlength = currentlength + len(suffixes[0]) + separatorLen
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
					if currentlength+len(prefixes[len(prefixes)-1])+separatorLen <= maxlength {
						contents = append([]string{prefixes[len(prefixes)-1]}, contents...)
						currentlength = currentlength + len(prefixes[len(prefixes)-1]) + separatorLen
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

func validateResourceType(resourceType string, resourceTypes []string) (bool, error) {
	isEmpty := len(resourceType) == 0 && len(resourceTypes) == 0
	if isEmpty {
		return false, fmt.Errorf("resource_type and resource_types parameters are empty, you must specify at least one resource type")
	}
	errorStrings := []string{}
	resourceList := resourceTypes
	if len(resourceType) > 0 {
		resourceList = append(resourceList, resourceType)
	}

	for _, resource := range resourceList {
		_, err := models.ValidateResourceType(resource)
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}
	if len(errorStrings) > 0 {
		return false, fmt.Errorf("%s", strings.Join(errorStrings, "\n"))
	}
	return true, nil
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
	if useSlug && !passthrough {
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

	// Extract minimum length from regex pattern (e.g., {4,48} -> 4)
	minLengthRegex := regexp.MustCompile(`\{(\d+),`)
	if matches := minLengthRegex.FindStringSubmatch(resource.ValidationRegExp); len(matches) > 1 {
		if minLength, err := strconv.Atoi(matches[1]); err == nil {
			// Add 2 for the required start and end characters (start letter + min content + end alphanumeric)
			totalMinLength := minLength + 2
			if len(resourceName) < totalMinLength {
				// For short prefixes, we need to ensure we meet the minimum length
				// while preserving the prefix and maintaining valid characters
				paddingNeeded := totalMinLength - len(resourceName)
				
				// Generate a valid suffix that meets the requirements
				validSuffix := make([]byte, paddingNeeded)
				for i := 0; i < paddingNeeded; i++ {
					if i == paddingNeeded-1 {
						// Last character must be alphanumeric
						validSuffix[i] = 'x'
					} else if i == 0 && resourceName == "" {
						// First character must be a letter if no prefix
						validSuffix[i] = 'a'
					} else {
						// Use random characters from suffix if available, otherwise use fallback
						if len(randomSuffix) > i && randomSuffix[i] != '-' {
							validSuffix[i] = randomSuffix[i]
						} else {
							validSuffix[i] = 'a'
						}
					}
				}
				
				resourceName = resourceName + string(validSuffix)
			}
		}
	}

	// Validate the final name against the pattern
	if !validationRegEx.MatchString(resourceName) {
		return "", fmt.Errorf("invalid name for CAF naming %s %s, the pattern %s doesn't match %s", resource.ResourceTypeName, name, resource.ValidationRegExp, resourceName)
	}

	// Validate the final name against the pattern
	if !validationRegEx.MatchString(resourceName) {
		return "", fmt.Errorf("invalid name for CAF naming %s %s, the pattern %s doesn't match %s", resource.ResourceTypeName, name, resource.ValidationRegExp, resourceName)
	}

	if !validationRegEx.MatchString(resourceName) {
		return "", fmt.Errorf("invalid name for CAF naming %s %s, the pattern %s doesn't match %s", resource.ResourceTypeName, name, resource.ValidationRegExp, resourceName)
	}

	return resourceName, nil
}
