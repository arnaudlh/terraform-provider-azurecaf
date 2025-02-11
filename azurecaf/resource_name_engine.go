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
	useSlug bool) string {

	// Filter out empty strings
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

	// Build components following precedence
	var components []string
	var essentialComponents []string
	var optionalComponents []string

	// First pass: collect components by type
	for _, part := range namePrecedence {
		switch part {
		case "prefixes":
			if len(filteredPrefixes) > 0 {
				essentialComponents = append(essentialComponents, filteredPrefixes...)
			}
		case "name":
			if name != "" {
				essentialComponents = append(essentialComponents, name)
			}
		case "slug":
			if useSlug && resourceDef != nil {
				if resourceDef.ResourceTypeName == "azurerm_recovery_services_vault" {
					essentialComponents = append(essentialComponents, "rsv")
				} else if resourceDef.CafPrefix != "" {
					essentialComponents = append(essentialComponents, resourceDef.CafPrefix)
				}
			}
		case "random":
			if randomSuffix != "" {
				optionalComponents = append(optionalComponents, randomSuffix)
			}
		case "suffixes":
			if len(filteredSuffixes) > 0 {
				optionalComponents = append(optionalComponents, filteredSuffixes...)
			}
		}
	}

	// Calculate total length including separators
	sepLen := len(separator)
	essentialLen := 0
	for i, comp := range essentialComponents {
		essentialLen += len(comp)
		if i > 0 {
			essentialLen += sepLen
		}
	}

	// Add essential components first
	components = append(components, essentialComponents...)

	// Add optional components if space permits
	remainingLen := maxLength - essentialLen - (len(components)-1)*sepLen
	for _, comp := range optionalComponents {
		newLen := len(comp)
		if len(components) > 0 {
			newLen += sepLen
		}
		if remainingLen >= newLen {
			components = append(components, comp)
			remainingLen -= newLen
		}
	}

	// Special handling for RSV max length
	if resourceDef != nil && resourceDef.ResourceTypeName == "azurerm_recovery_services_vault" {
		maxLength = 50 // RSV has a fixed max length
	}

	// Join components with separator
	result := strings.Join(components, separator)

	// Ensure we don't exceed max length while preserving component boundaries
	if len(result) > maxLength {
		parts := strings.Split(result, separator)
		var truncated []string
		currentLength := 0
		sepLen := len(separator)

		// First pass: essential components (prefixes, name, and slug)
		for i, part := range parts {
			newLength := currentLength
			if len(truncated) > 0 {
				newLength += sepLen
			}
			newLength += len(part)

			if (i < len(filteredPrefixes) || part == name || part == slug) && newLength <= maxLength {
				if len(truncated) > 0 {
					currentLength += sepLen
				}
				truncated = append(truncated, part)
				currentLength += len(part)
			}
		}

		// Second pass: suffixes and random components
		for _, part := range parts {
			if contains(truncated, part) {
				continue
			}

			newLength := currentLength
			if len(truncated) > 0 {
				newLength += sepLen
			}
			newLength += len(part)

			if newLength <= maxLength {
				if len(truncated) > 0 {
					currentLength += sepLen
				}
				truncated = append(truncated, part)
				currentLength += len(part)
			}
		}

		result = strings.Join(truncated, separator)
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
		if resourceTypeName == "azurerm_recovery_services_vault" {
			slug = "rsv"
		} else {
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
	resourceName := composeName(separator, filteredPrefixes, name, slug, filteredSuffixes, randomSuffix, resource.MaxLength, namePrecedence, resource, useSlug)

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
