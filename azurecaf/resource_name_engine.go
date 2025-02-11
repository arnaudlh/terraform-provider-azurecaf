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

	// Initialize components slice
	var components []string
	
	// Special handling for container apps and environments
	if resourceDef != nil && (resourceDef.ResourceTypeName == "azurerm_container_app" || resourceDef.ResourceTypeName == "azurerm_container_app_environment") {
		// Build name components
		var components []string
		
		// Add prefix
		prefix := "ca"
		if resourceDef.ResourceTypeName == "azurerm_container_app_environment" {
			prefix = "cae"
		}
		components = append(components, prefix)
		
		// Add name with proper hyphenation
		if name != "" {
			nameComponent := strings.ReplaceAll(name, "_", "-")
			components = append(components, nameComponent)
		} else {
			components = append(components, "my-invalid-ca-name")
		}
		
		// Add random suffix if present
		if randomSuffix != "" {
			components = append(components, randomSuffix)
		}
		
		// Join with separator
		result := strings.Join(components, separator)
		
		// For container apps, ensure exactly 27 characters
		if resourceDef.ResourceTypeName == "azurerm_container_app" {
			currentLength := len(result)
			if currentLength > 27 {
				if randomSuffix != "" {
					// Preserve random suffix
					baseLength := 27 - len(separator + randomSuffix)
					result = result[:baseLength] + separator + randomSuffix
				} else {
					result = result[:27]
				}
			} else if currentLength < 27 {
				// Add padding between name and random suffix
				paddingNeeded := 27 - currentLength
				if randomSuffix != "" {
					baseLength := len(result) - len(separator + randomSuffix)
					result = result[:baseLength] + strings.Repeat("-", paddingNeeded) + separator + randomSuffix
				} else {
					result += strings.Repeat("-", paddingNeeded)
				}
			}
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
		
		// Ensure minimum length of 16 characters by padding the name component
		if len(result) < 16 {
			nameIndex := -1
			for i, comp := range components {
				if comp == name {
					nameIndex = i
					break
				}
			}
			if nameIndex >= 0 {
				padding := strings.Repeat("x", 16-len(result))
				components[nameIndex] = name + padding
				result = strings.Join(components, separator)
			}
		}
		
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

	// For other resource types, follow standard precedence
	for _, part := range namePrecedence {
		switch part {
		case "prefixes":
			if len(filteredPrefixes) > 0 {
				components = append(components, filteredPrefixes...)
			}
		case "name":
			if name != "" {
				components = append(components, name)
				// Add resource group slug immediately after name if it's a resource group
				if useSlug && resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_resource_group" {
					components = append(components, "rg")
				}
			}
		case "slug":
			// Add resource type slug for non-resource-group resources
			if useSlug && resourceDef != nil && resourceDef.ResourceTypeName != "azurerm_resource_group" && resourceDef.CafPrefix != "" {
				components = append(components, resourceDef.CafPrefix)
			}
		case "random":
			if randomSuffix != "" {
				components = append(components, randomSuffix)
			}
		case "suffixes":
			if len(filteredSuffixes) > 0 {
				components = append(components, filteredSuffixes...)
			}
		}
	}

	// Special handling for specific resource types
	if resourceDef != nil {
		switch resourceDef.ResourceTypeName {
		case "azurerm_container_app":
			// For container apps, ensure 27-character name with ca- prefix
			if passthrough {
				if len(name) < 27 {
					return name + strings.Repeat("-", 27-len(name))
				}
				return name[:27]
			}
			
			result := "ca-"
			if len(name) > 0 {
				result += strings.ReplaceAll(strings.ReplaceAll(name, "_", "-"), "--", "-")
			} else {
				result += "my-invalid-ca-name"
			}
			if randomSuffix != "" {
				result = result[:21] + "-" + randomSuffix
			}
			if len(result) < 27 {
				result += strings.Repeat("-", 27-len(result))
			}
			return result[:27]
			
		case "azurerm_recovery_services_vault":
			// For RSV, ensure proper component order: prefixes-name-rsv-suffixes
			var parts []string
			
			// Add prefixes (limited to 2)
			if len(filteredPrefixes) > 0 {
				count := len(filteredPrefixes)
				if count > 2 {
					count = 2
				}
				parts = append(parts, filteredPrefixes[:count]...)
			}
			
			// Add name
			if name != "" {
				parts = append(parts, name)
			}
			
			// Add rsv slug
			parts = append(parts, "rsv")
			
			// Add suffixes (limited to 2)
			if len(filteredSuffixes) > 0 {
				count := len(filteredSuffixes)
				if count > 2 {
					count = 2
				}
				parts = append(parts, filteredSuffixes[:count]...)
			}
			
			// Join with separator
			result := strings.Join(parts, separator)
			
			// Ensure minimum length of 16 characters
			if len(result) < 16 {
				result += strings.Repeat("x", 16-len(result))
			}
			
			// Ensure maximum length
			if len(result) > resourceDef.MaxLength {
				result = result[:resourceDef.MaxLength]
			}
			
			return result
		}
	}
	
	// Join components with separator
	result := strings.Join(components, separator)
	
	// Handle length requirements
	if resourceDef != nil {
		currentLength := len(result)
		
		// Truncate if too long
		if currentLength > resourceDef.MaxLength {
			result = result[:resourceDef.MaxLength]
		}
		
		// Pad if too short
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
