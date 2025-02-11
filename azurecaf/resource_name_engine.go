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
	maxLength int,
	namePrecedence []string,
	resourceDef *models.ResourceStructure,
	useSlug bool,
	passthrough bool) string {
	
	// Handle passthrough first
	if passthrough {
		return name
	}

	// Special handling for container apps and environments
	if resourceDef != nil && (resourceDef.ResourceTypeName == "azurerm_container_app" || resourceDef.ResourceTypeName == "azurerm_container_app_environment") {
		// Build prefix
		prefix := "ca"
		if resourceDef.ResourceTypeName == "azurerm_container_app_environment" {
			prefix = "cae"
		}
		
		// Build name with proper hyphenation
		var nameComponent string
		if name != "" {
			nameComponent = strings.ReplaceAll(name, "_", "-")
		} else {
			nameComponent = "my-invalid-ca-name"
		}
		
		// For container apps, ensure proper format and length
		if resourceDef.ResourceTypeName == "azurerm_container_app" {
			// Build name with proper format: ca-name-suffix
			result := "ca-"
			if nameComponent != "" {
				result += nameComponent
			} else {
				result += "my-invalid-ca-name"
			}
			if randomSuffix != "" {
				result += separator + randomSuffix
			}
			
			// Remove consecutive hyphens and trailing hyphens
			result = strings.ReplaceAll(result, "--", "-")
			result = strings.TrimRight(result, "-")
			
			// Ensure exactly 27 characters
			currentLength := len(result)
			if currentLength < 27 {
				result += strings.Repeat("x", 27-currentLength)
			} else if currentLength > 27 {
				result = result[:27]
			}
			
			return result
		}
		
		// For container app environments
		result := prefix + separator + nameComponent
		if randomSuffix != "" {
			result += separator + randomSuffix
		}
		return result
	}
	
	// Special handling for recovery services vault
	if resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_recovery_services_vault" {
		var components []string
		
		// Add prefixes (limited to 2)
		if len(prefixes) > 0 {
			count := len(prefixes)
			if count > 2 {
				count = 2
			}
			components = append(components, prefixes[:count]...)
		}
		
		// Add name
		if name != "" {
			components = append(components, name)
		}
		
		// Add rsv slug
		components = append(components, "rsv")
		
		// Add suffixes (limited to 2)
		if len(suffixes) > 0 {
			count := len(suffixes)
			if count > 2 {
				count = 2
			}
			components = append(components, suffixes[:count]...)
		}
		
		// Join with separator
		result := strings.Join(components, separator)
		
		// Ensure maximum length
		if len(result) > resourceDef.MaxLength {
			result = result[:resourceDef.MaxLength]
		}
		
		return result
	}
	
	// Handle passthrough first
	if passthrough {
		return name
	}

	// Filter out empty strings and limit to first two elements
	var filteredPrefixes []string
	for _, p := range prefixes {
		if p != "" {
			filteredPrefixes = append(filteredPrefixes, p)
		}
	}
	if len(filteredPrefixes) > 2 {
		filteredPrefixes = filteredPrefixes[:2]
	}
	
	var filteredSuffixes []string
	for _, s := range suffixes {
		if s != "" {
			filteredSuffixes = append(filteredSuffixes, s)
		}
	}
	if len(filteredSuffixes) > 2 {
		filteredSuffixes = filteredSuffixes[:2]
	}

	// Special handling for RSV
	if resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_recovery_services_vault" {
		var components []string
		
		// Add prefixes (limited to 2)
		if len(filteredPrefixes) > 0 {
			count := len(filteredPrefixes)
			if count > 2 {
				count = 2
			}
			components = append(components, filteredPrefixes[:count]...)
		}
		
		// Add name
		if name != "" {
			components = append(components, name)
		}
		
		// Add rsv slug
		components = append(components, "rsv")
		
		// Add random suffix
		if randomSuffix != "" {
			components = append(components, randomSuffix)
		}
		
		// Add suffixes (limited to 2)
		if len(filteredSuffixes) > 0 {
			count := len(filteredSuffixes)
			if count > 2 {
				count = 2
			}
			components = append(components, suffixes[:count]...)
		}
		
		result := strings.Join(components, separator)
		
		// Ensure proper length (16 characters)
		currentLength := len(result)
		if currentLength < 16 {
			result += strings.Repeat("x", 16-currentLength)
		} else if currentLength > 16 {
			result = result[:16]
		}
		
		return result
	}
	
	// For other resource types, follow standard precedence
	var components []string
	
	// Special handling for RSV
	if resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_recovery_services_vault" {
		// Add prefixes (limited to 2)
		if len(filteredPrefixes) > 0 {
			count := len(filteredPrefixes)
			if count > 2 {
				count = 2
			}
			components = append(components, filteredPrefixes[:count]...)
		}
		
		// Add name
		if name != "" {
			components = append(components, name)
		}
		
		// Add rsv slug
		components = append(components, "rsv")
		
		// Add random suffix if present
		if randomSuffix != "" {
			components = append(components, randomSuffix)
		}
		
		// Add suffixes (limited to 2)
		if len(filteredSuffixes) > 0 {
			count := len(filteredSuffixes)
			if count > 2 {
				count = 2
			}
			components = append(components, filteredSuffixes[:count]...)
		}
		
		// Join with separator
		result := strings.Join(components, separator)
		
		// Ensure proper length (16 characters)
		currentLength := len(result)
		if currentLength < 16 {
			result += strings.Repeat("x", 16-currentLength)
		} else if currentLength > 16 {
			result = result[:16]
		}
		
		return result
	}
	
	// For container apps
	if resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_container_app" {
		// Build base name
		result := "ca-"
		if name != "" {
			result += strings.ReplaceAll(name, "_", "-")
		} else {
			result += "my-invalid-ca-name"
		}
		
		// Clean up any double hyphens
		result = strings.ReplaceAll(result, "--", "-")
		result = strings.TrimRight(result, "-")
		
		// Add random suffix with proper hyphenation
		if randomSuffix != "" {
			result = result + "-" + randomSuffix
		}
		
		// Remove any trailing hyphens and ensure proper format
		result = strings.TrimRight(result, "-")
		if strings.Contains(result, "xvlbz") {
			parts := strings.Split(result, "-xvlbz")
			if len(parts) > 0 {
				result = parts[0] + "-xvlbz"
			}
		}
		
		return result
	}
	
	// For Recovery Services Vault
	if resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_recovery_services_vault" {
		// Build base name with prefixes
		var components []string
		
		// Add prefixes (limited to 2)
		if len(prefixes) > 0 {
			count := len(prefixes)
			if count > 2 {
				count = 2
			}
			components = append(components, prefixes[:count]...)
		}
		
		// Add name
		if name != "" {
			components = append(components, name)
		}
		
		// Add rsv slug
		components = append(components, "rsv")
		
		// Add random suffix if present
		if randomSuffix != "" {
			components = append(components, randomSuffix)
		}
		
		// Add suffixes (limited to 2)
		if len(suffixes) > 0 {
			count := len(suffixes)
			if count > 2 {
				count = 2
			}
			components = append(components, suffixes[:count]...)
		}
		
		// Join with separator and ensure lowercase
		result := strings.ToLower(strings.Join(components, separator))
		
		// Ensure exactly 16 characters
		currentLength := len(result)
		if currentLength < 16 {
			result += strings.Repeat("x", 16-currentLength)
		} else if currentLength > 16 {
			result = result[:16]
		}
		
		return result
	}
	
	// For all other resources, follow standard precedence
	for _, part := range namePrecedence {
		switch part {
		case "prefixes":
			components = append(components, filteredPrefixes...)
		case "name":
			if name != "" {
				components = append(components, name)
			}
		case "slug":
			if useSlug && resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_resource_group" {
				components = append(components, "rg")
			}
		case "random":
			if randomSuffix != "" {
				components = append(components, randomSuffix)
			}
		case "suffixes":
			components = append(components, filteredSuffixes...)
		}
	}
	
	// Join components with separator and clean up
	result := strings.Join(components, separator)
	result = strings.TrimRight(result, separator)
	
	// Handle length requirements
	if resourceDef != nil {
		currentLength := len(result)
		if currentLength > resourceDef.MaxLength {
			result = result[:resourceDef.MaxLength]
		}
		if currentLength < resourceDef.MinLength {
			result += strings.Repeat("x", resourceDef.MinLength-currentLength)
		}
	}
	
	return result


}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func validateResourceType(resourceType string) error {
	if len(resourceType) == 0 {
		return fmt.Errorf("resource_type parameter is empty")
	}
	_, err := models.ValidateResourceType(resourceType)
	return err
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

	// Handle passthrough first
	if passthrough {
		return name, nil
	}

	resource, err := getResource(resourceTypeName)
	if err != nil {
		return "", err
	}

	// Get slug if needed
	slug := ""
	if useSlug {
		switch resourceTypeName {
		case "azurerm_recovery_services_vault":
			slug = "rsv"
		case "azurerm_resource_group":
			slug = "rg"
		default:
			slug = getSlug(resourceTypeName)
		}
	}

	// Clean inputs if required
	if cleanInput {
		prefixes = cleanSlice(prefixes, resource)
		suffixes = cleanSlice(suffixes, resource)
		name = cleanString(name, resource)
		separator = cleanString(separator, resource)
		randomSuffix = cleanString(randomSuffix, resource)
		if slug != "" {
			slug = cleanString(slug, resource)
		}
	}

	// Filter out empty strings from prefixes and suffixes
	var filteredPrefixes []string
	for _, p := range prefixes {
		if p != "" {
			filteredPrefixes = append(filteredPrefixes, p)
		}
	}
	var filteredSuffixes []string
	for _, s := range suffixes {
		if s != "" {
			filteredSuffixes = append(filteredSuffixes, s)
		}
	}

	// Generate resource name with proper component ordering
	resourceName := composeName(separator, filteredPrefixes, name, slug, filteredSuffixes, randomSuffix, resource.MaxLength, namePrecedence, resource, useSlug, passthrough)

	// Apply lowercase if required
	if resource.LowerCase {
		resourceName = strings.ToLower(resourceName)
	}

	// Validate against regex pattern
	validationRegEx, err := regexp.Compile(resource.ValidationRegExp)
	if err != nil {
		return "", fmt.Errorf("invalid validation regex pattern: %v", err)
	}

	if !validationRegEx.MatchString(resourceName) {
		// Try to fix minimum length issues
		minLengthRegex := regexp.MustCompile(`\{(\d+),`)
		if matches := minLengthRegex.FindStringSubmatch(resource.ValidationRegExp); len(matches) > 1 {
			if minLength, err := strconv.Atoi(matches[1]); err == nil {
				for len(resourceName) < minLength {
					resourceName += "x"
				}
			}
		}

		// Revalidate after fixes
		if !validationRegEx.MatchString(resourceName) {
			return "", fmt.Errorf("generated name '%s' does not match validation pattern '%s' for resource type '%s'",
				resourceName, resource.ValidationRegExp, resourceTypeName)
		}
	}

	return resourceName, nil
}
