package schemas

import (
	"context"
	"encoding/base64"
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
	resource, err := models.ValidateResourceType(resourceType)
	if err != nil {
		return err
	}

	if resource.ValidationRegExp != "" {
		pattern, err := regexp.Compile(resource.ValidationRegExp)
		if err != nil {
			return fmt.Errorf("invalid validation regex pattern for resource type %s: %v", resourceType, err)
		}
		if !pattern.MatchString(name) {
			return fmt.Errorf("resource name %s does not match required pattern: %s", name, resource.ValidationRegExp)
		}
	}

	if resource.LowerCase && name != strings.ToLower(name) {
		return fmt.Errorf("resource name %s must be lowercase", name)
	}

	nameLen := len(name)
	if nameLen < resource.MinLength || (resource.MaxLength > 0 && nameLen > resource.MaxLength) {
		return fmt.Errorf("resource name %s length must be between %d and %d", name, resource.MinLength, resource.MaxLength)
	}

	if resource.CafPrefix != "" {
		prefix := strings.ToLower(resource.CafPrefix)
		if !strings.HasPrefix(strings.ToLower(name), prefix) {
			return fmt.Errorf("resource name %s must start with prefix %s", name, prefix)
		}
	}

	return nil
}

// ValidateResourceNameWithSlug validates a resource name with its slug placement
func ValidateResourceNameWithSlug(resourceType, name string) error {
	if err := ValidateResourceNameSchema(resourceType, name); err != nil {
		return err
	}

	resource, _ := models.ValidateResourceType(resourceType)
	if resource.CafPrefix != "" {
		prefix := strings.ToLower(resource.CafPrefix)
		nameLower := strings.ToLower(name)
		if !strings.Contains(nameLower, prefix) {
			return fmt.Errorf("resource name %s must contain slug %s", name, prefix)
		}

		slugIndex := strings.Index(nameLower, prefix)
		if slugIndex > 0 {
			prevChar := rune(name[slugIndex-1])
			if !strings.ContainsRune("-_.", prevChar) {
				return fmt.Errorf("resource name %s has incorrectly placed slug %s - should be at start or after separator (-, _, or .)", name, prefix)
			}
		}
	}

	return nil
}

// ValidateResourceNameInSchema validates a resource name against schema constraints
func ValidateResourceNameInSchema(d interface{}) error {
	var resourceType, name string
	var useSlug bool = true

	switch v := d.(type) {
	case *schema.ResourceData:
		resourceType = v.Get("resource_type").(string)
		name = v.Get("name").(string)
		if v, ok := v.GetOk("use_slug"); ok {
			useSlug = v.(bool)
		}
	case *schema.ResourceDiff:
		resourceType = v.Get("resource_type").(string)
		name = v.Get("name").(string)
		if v, ok := v.GetOk("use_slug"); ok {
			useSlug = v.(bool)
		}
	default:
		return fmt.Errorf("unsupported schema type for validation")
	}

	if resourceType != "" {
		resource, err := models.GetResourceStructure(resourceType)
		if err != nil {
			return err
		}

		if err := ValidateResourceName(name, resource); err != nil {
			return err
		}

		if useSlug && resource.CafPrefix != "" {
			expectedPrefix := resource.CafPrefix
			nameLower := strings.ToLower(name)
			prefixLower := strings.ToLower(expectedPrefix)

			if !strings.Contains(nameLower, prefixLower) {
				return fmt.Errorf("resource name %s must contain the slug '%s' for resource type '%s'", name, expectedPrefix, resourceType)
			}

			slugIndex := strings.Index(nameLower, prefixLower)
			if slugIndex > 0 {
				prevChar := rune(name[slugIndex-1])
				if !strings.ContainsRune("-_.", prevChar) {
					return fmt.Errorf("resource name %s has incorrectly placed slug '%s' - should be at start or after separator (-, _, or .)", name, expectedPrefix)
				}
			}

			slugEndIndex := slugIndex + len(prefixLower)
			if slugEndIndex < len(name) {
				nextChar := rune(name[slugEndIndex])
				if !strings.ContainsRune("-_.", nextChar) {
					return fmt.Errorf("resource name %s must have a separator (-, _, or .) after the slug '%s'", name, expectedPrefix)
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
	}

	return nil
}

// ValidateResourceNameInSchemaWithPrefix validates a resource name with prefix against schema constraints
func ValidateResourceNameInSchemaWithPrefix(d *schema.ResourceData) error {
	resourceType := d.Get("resource_type").(string)
	name := d.Get("name").(string)
	prefixes, ok := d.Get("prefixes").([]interface{})
	if ok && len(prefixes) > 0 {
		prefix := prefixes[0].(string)
		if !strings.HasPrefix(name, prefix) {
			return fmt.Errorf("resource name %s must start with prefix %s", name, prefix)
		}
	}
	return ValidateResourceNameWithSlug(resourceType, name)
}

// ValidateResourceNameInSchemaWithTypes validates a resource name against schema constraints for multiple resource types
func ValidateResourceNameInSchemaWithTypes(d *schema.ResourceData) error {
	resourceTypes, ok := d.Get("resource_types").([]interface{})
	if !ok || len(resourceTypes) == 0 {
		return fmt.Errorf("resource_types must be provided")
	}

	name := d.Get("name").(string)
	for _, rt := range resourceTypes {
		resourceType := rt.(string)
		if err := ValidateResourceNameWithSlug(resourceType, name); err != nil {
			return fmt.Errorf("validation failed for resource type %s: %v", resourceType, err)
		}
	}
	return nil
}

// BaseSchema returns the base schema for all resource types
func validateResourceTypes(i interface{}, k string) ([]string, []error) {
	var errs []error
	v, ok := i.([]interface{})
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be []interface{}", k)}
	}
	resourceMapsKeys := getResourceMaps()
	for _, rt := range v {
		resourceType := rt.(string)
		if !stringInSlice(resourceType, resourceMapsKeys) {
			errs = append(errs, fmt.Errorf("resource type %s is not supported", resourceType))
		}
	}
	return nil, errs
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func BaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
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
			Optional:     true,
			ValidateFunc: validation.StringInSlice(getResourceMaps(), false),
		},
		"random_seed": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}
