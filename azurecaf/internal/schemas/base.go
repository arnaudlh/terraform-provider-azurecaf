package schemas

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	models "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

// ResourceOperations contains the CRUD operations for all schema versions
// ResourceOperations contains the CRUD operations for all schema versions
var (
	resourceNameCreate = func(d *schema.ResourceData, m interface{}) error {
		if err := ValidateResourceNameInSchema(d); err != nil {
			return err
		}
		name := d.Get("name").(string)
		resourceType := d.Get("resource_type").(string)
		if resourceType != "" {
			if resource, err := models.GetResourceStructure(resourceType); err == nil {
				if err := ValidateResourceName(name, resource); err != nil {
					return err
				}
			}
		}
		d.SetId(name)
		return resourceNameRead(d, m)
	}

	resourceNameRead = func(d *schema.ResourceData, m interface{}) error {
		return nil
	}

	resourceNameUpdate = func(d *schema.ResourceData, m interface{}) error {
		if err := ValidateResourceNameInSchema(d); err != nil {
			return err
		}
		name := d.Get("name").(string)
		resourceType := d.Get("resource_type").(string)
		if resourceType != "" {
			if resource, err := models.GetResourceStructure(resourceType); err == nil {
				if err := ValidateResourceName(name, resource); err != nil {
					return err
				}
			}
		}
		return resourceNameRead(d, m)
	}

	resourceNameDelete = func(d *schema.ResourceData, m interface{}) error {
		d.SetId("")
		return nil
	}
)

// getResourceMaps returns a list of all supported resource types
func getResourceMaps() []string {
	resourceMapsKeys := make([]string, 0, len(models.ResourceDefinitions))
	for k := range models.ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}
	return resourceMapsKeys
}

// ValidateResourceName validates a resource name against its defined constraints
func ValidateResourceName(name string, resource *models.ResourceStructure) error {
	if name == "" {
		return nil
	}

	nameLen := len(name)
	if nameLen < resource.MinLength || (resource.MaxLength > 0 && nameLen > resource.MaxLength) {
		return fmt.Errorf("resource name %s length must be between %d and %d", name, resource.MinLength, resource.MaxLength)
	}

	if resource.LowerCase && name != strings.ToLower(name) {
		return fmt.Errorf("resource name %s must be lowercase", name)
	}

	if resource.ValidationRegExp != "" {
		pattern, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			return fmt.Errorf("invalid validation regex pattern: %v", err)
		}
		if !pattern.MatchString(name) {
			return fmt.Errorf("resource name %s does not match required pattern %s", name, resource.ValidationRegExp)
		}
	}

	return nil
}

// ValidateResourceNameSchema validates a resource name against its schema constraints
func ValidateResourceNameSchema(resourceType, name string) error {
	if resourceType == "" || name == "" {
		return nil
	}

	resource, err := models.ValidateResourceType(resourceType)
	if err != nil {
		return err
	}

	nameLen := len(name)
	if nameLen < resource.MinLength || (resource.MaxLength > 0 && nameLen > resource.MaxLength) {
		return fmt.Errorf("resource name %s length must be between %d and %d", name, resource.MinLength, resource.MaxLength)
	}

	if resource.LowerCase && name != strings.ToLower(name) {
		return fmt.Errorf("resource name %s must be lowercase", name)
	}

	if resource.ValidationRegExp != "" {
		pattern, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			return fmt.Errorf("invalid validation regex pattern for resource type %s: %v", resourceType, err)
		}
		if !pattern.MatchString(name) {
			return fmt.Errorf("resource name %s does not match required pattern %s for resource type %s", name, resource.ValidationRegExp, resourceType)
		}
	}

	return nil
}

// ValidateResourceNameWithSlug validates a resource name with its slug placement
func ValidateResourceNameWithSlug(resourceType, name string, useSlug bool) error {
	if resourceType == "" || name == "" {
		return nil
	}

	resource, err := models.ValidateResourceType(resourceType)
	if err != nil {
		return err
	}

	if resource.LowerCase && name != strings.ToLower(name) {
		return fmt.Errorf("resource name %s must be lowercase for resource type %s", name, resourceType)
	}

	nameLength := len(name)
	if nameLength < resource.MinLength || nameLength > resource.MaxLength {
		return fmt.Errorf("resource name %s length (%d) must be between %d and %d characters for resource type %s", 
			name, nameLength, resource.MinLength, resource.MaxLength, resourceType)
	}

	if useSlug && resource.CafPrefix != "" {
		prefix := strings.ToLower(resource.CafPrefix)
		nameLower := strings.ToLower(name)
		if !strings.Contains(nameLower, prefix) {
			return fmt.Errorf("resource name %s must contain slug %s for resource type %s", name, prefix, resourceType)
		}

		slugIndex := strings.Index(nameLower, prefix)
		if slugIndex > 0 {
			prevChar := rune(name[slugIndex-1])
			if !strings.ContainsRune("-_.", prevChar) {
				return fmt.Errorf("resource name %s has incorrectly placed slug %s - should be at start or after separator (-, _, or .)", name, prefix)
			}
		}

		slugEndIndex := slugIndex + len(prefix)
		if slugEndIndex < len(name) {
			nextChar := rune(name[slugEndIndex])
			if !strings.ContainsRune("-_.", nextChar) {
				return fmt.Errorf("resource name %s must have a separator (-, _, or .) after the slug %s", name, prefix)
			}
		}
	}

	if resource.ValidationRegExp != "" {
		pattern, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			return fmt.Errorf("invalid validation regex pattern for resource type %s: %v", resourceType, err)
		}
		if !pattern.MatchString(name) {
			return fmt.Errorf("resource name %s does not match required pattern %s for resource type %s", name, resource.ValidationRegExp, resourceType)
		}
	}

	return nil
}

// ValidateResourceNameInSchema validates a resource name against schema constraints
func ValidateResourceNameInSchema(d interface{}) error {
	var resourceType, name string
	var useSlug bool = true
	var resourceTypes []interface{}
	var prefixes, suffixes []interface{}
	var separator string = "-"

	switch v := d.(type) {
	case *schema.ResourceData:
		resourceType = v.Get("resource_type").(string)
		name = v.Get("name").(string)
		if v, ok := v.GetOk("use_slug"); ok {
			useSlug = v.(bool)
		}
		if v, ok := v.GetOk("resource_types"); ok {
			resourceTypes = v.([]interface{})
		}
		if v, ok := v.GetOk("prefixes"); ok {
			prefixes = v.([]interface{})
		}
		if v, ok := v.GetOk("suffixes"); ok {
			suffixes = v.([]interface{})
		}
		if v, ok := v.GetOk("separator"); ok {
			separator = v.(string)
		}
	case *schema.ResourceDiff:
		resourceType = v.Get("resource_type").(string)
		name = v.Get("name").(string)
		if v, ok := v.GetOk("use_slug"); ok {
			useSlug = v.(bool)
		}
		if v, ok := v.GetOk("resource_types"); ok {
			resourceTypes = v.([]interface{})
		}
		if v, ok := v.GetOk("prefixes"); ok {
			prefixes = v.([]interface{})
		}
		if v, ok := v.GetOk("suffixes"); ok {
			suffixes = v.([]interface{})
		}
		if v, ok := v.GetOk("separator"); ok {
			separator = v.(string)
		}
	default:
		return fmt.Errorf("unsupported schema type for validation")
	}

	if name == "" {
		return nil
	}

	if len(separator) != 1 || !strings.ContainsRune("-_.", rune(separator[0])) {
		return fmt.Errorf("invalid separator %s - must be one of: -, _, .", separator)
	}

	for i := 0; i < len(name)-1; i++ {
		if strings.ContainsRune("-_.", rune(name[i])) && strings.ContainsRune("-_.", rune(name[i+1])) {
			return fmt.Errorf("resource name %s contains consecutive separators at position %d", name, i)
		}
	}

	if len(prefixes) > 0 {
		for _, prefix := range prefixes {
			prefixStr := prefix.(string)
			if !strings.HasPrefix(name, prefixStr) {
				return fmt.Errorf("resource name %s must start with prefix %s", name, prefixStr)
			}
		}
	}

	if len(suffixes) > 0 {
		for _, suffix := range suffixes {
			suffixStr := suffix.(string)
			if !strings.HasSuffix(name, suffixStr) {
				return fmt.Errorf("resource name %s must end with suffix %s", name, suffixStr)
			}
		}
	}

	validateResource := func(rt string) error {
		resource, err := models.GetResourceStructure(rt)
		if err != nil {
			return fmt.Errorf("invalid resource type %s: %v", rt, err)
		}

		if useSlug && resource.CafPrefix != "" {
			slugIndex := strings.Index(strings.ToLower(name), strings.ToLower(resource.CafPrefix))
			if slugIndex == -1 {
				return fmt.Errorf("resource name %s must contain the slug '%s'", name, resource.CafPrefix)
			}

			if slugIndex > 0 && !strings.ContainsRune("-_.", rune(name[slugIndex-1])) {
				return fmt.Errorf("resource name %s has incorrectly placed slug '%s' - should be at start or after separator", name, resource.CafPrefix)
			}

			slugEndIndex := slugIndex + len(resource.CafPrefix)
			if slugEndIndex < len(name) && !strings.ContainsRune("-_.", rune(name[slugEndIndex])) {
				return fmt.Errorf("resource name %s must have a separator after the slug '%s'", name, resource.CafPrefix)
			}
		}

		if resource.ValidationRegExp != "" {
			pattern, err := regexp.Compile(resource.ValidationRegExp)
			if err != nil {
				return fmt.Errorf("invalid validation regex pattern for resource type %s: %v", rt, err)
			}
			if !pattern.MatchString(name) {
				return fmt.Errorf("resource name %s does not match required pattern %s for resource type %s", name, resource.ValidationRegExp, rt)
			}
		}

		return nil
	}

	if resourceType == "" {
		if len(resourceTypes) == 0 {
			return nil
		}
		for _, rt := range resourceTypes {
			rtStr := rt.(string)
			if err := validateResource(rtStr); err != nil {
				return fmt.Errorf("validation failed for resource type %s: %v", rtStr, err)
			}
		}
		return nil
	}

	return validateResource(resourceType)
}

func validateSingleResource(resourceType, name string, useSlug bool) error {
	if resourceType == "" || name == "" {
		return nil
	}

	resource, err := models.GetResourceStructure(resourceType)
	if err != nil {
		return fmt.Errorf("invalid resource type %s: %v", resourceType, err)
	}

	if resource.LowerCase && name != strings.ToLower(name) {
		return fmt.Errorf("resource name %s must be lowercase for resource type %s", name, resourceType)
	}

	nameLength := len(name)
	if nameLength < resource.MinLength || (resource.MaxLength > 0 && nameLength > resource.MaxLength) {
		return fmt.Errorf("resource name %s length (%d) must be between %d and %d characters for resource type %s", 
			name, nameLength, resource.MinLength, resource.MaxLength, resourceType)
	}

	// Check for consecutive separators
	for i := 0; i < len(name)-1; i++ {
		if strings.ContainsRune("-_.", rune(name[i])) && strings.ContainsRune("-_.", rune(name[i+1])) {
			return fmt.Errorf("resource name %s contains consecutive separators at position %d", name, i)
		}
	}

	// Validate slug placement
	if useSlug && resource.CafPrefix != "" {
		expectedPrefix := resource.CafPrefix
		nameLower := strings.ToLower(name)
		prefixLower := strings.ToLower(expectedPrefix)

		slugIndex := strings.Index(nameLower, prefixLower)
		if slugIndex == -1 {
			return fmt.Errorf("resource name %s must contain the slug '%s' for resource type '%s'", name, expectedPrefix, resourceType)
		}

		// Validate slug position (must be at start or after separator)
		if slugIndex > 0 {
			prevChar := rune(name[slugIndex-1])
			if !strings.ContainsRune("-_.", prevChar) {
				return fmt.Errorf("resource name %s has incorrectly placed slug '%s' - should be at start or after separator (-, _, or .)", name, expectedPrefix)
			}
		}

		// Validate character after slug (must be separator or end of string)
		slugEndIndex := slugIndex + len(prefixLower)
		if slugEndIndex < len(name) {
			nextChar := rune(name[slugEndIndex])
			if !strings.ContainsRune("-_.", nextChar) {
				return fmt.Errorf("resource name %s must have a separator (-, _, or .) after the slug '%s'", name, expectedPrefix)
			}
		}
	}

	// Validate against regex pattern
	if resource.ValidationRegExp != "" {
		pattern, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			return fmt.Errorf("invalid validation regex pattern for resource type %s: %v", resourceType, err)
		}
		if !pattern.MatchString(name) {
			return fmt.Errorf("resource name %s does not match required pattern %s for resource type %s", name, resource.ValidationRegExp, resourceType)
		}
	}

	return nil
}

// ValidateResourceNameInSchemaWithPrefix validates a resource name with prefix against schema constraints
func ValidateResourceNameInSchemaWithPrefix(d *schema.ResourceData) error {
	resourceType := d.Get("resource_type").(string)
	name := d.Get("name").(string)
	useSlug := true
	if v, ok := d.GetOk("use_slug"); ok {
		useSlug = v.(bool)
	}

	prefixes, ok := d.Get("prefixes").([]interface{})
	if ok && len(prefixes) > 0 {
		prefix := prefixes[0].(string)
		if !strings.HasPrefix(name, prefix) {
			return fmt.Errorf("resource name %s must start with prefix %s", name, prefix)
		}
	}

	if err := validateSingleResource(resourceType, name, useSlug); err != nil {
		return err
	}

	return nil
}

// ValidateResourceNameInSchemaWithTypes validates a resource name against schema constraints for multiple resource types
func ValidateResourceNameInSchemaWithTypes(d *schema.ResourceData) error {
	resourceTypes, ok := d.Get("resource_types").([]interface{})
	if !ok || len(resourceTypes) == 0 {
		return fmt.Errorf("resource_types must be provided")
	}

	name := d.Get("name").(string)
	useSlug := true
	if v, ok := d.GetOk("use_slug"); ok {
		useSlug = v.(bool)
	}

	for _, rt := range resourceTypes {
		resourceType := rt.(string)
		if err := validateSingleResource(resourceType, name, useSlug); err != nil {
			return fmt.Errorf("validation failed for resource type %s: %v", resourceType, err)
		}
	}
	return nil
}

// BaseSchema returns the base schema for all resource types


func BaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"prefixes": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Optional: true,
		},
		"suffixes": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Optional: true,
		},
		"random_length": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
			Default:      0,
		},
		"result": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"results": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		},
		"separator": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "-",
			ValidateFunc: validation.StringInSlice([]string{"-", "_", "."}, false),
		},
		"clean_input": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"passthrough": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"resource_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(getResourceMaps(), false),
		},
		"random_seed": {
			Type:     schema.TypeInt,
			Optional: true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"use_slug": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"random_string": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}
