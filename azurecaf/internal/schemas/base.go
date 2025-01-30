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
func ValidateResourceNameInSchema(d *schema.ResourceData) error {
	resourceType := d.Get("resource_type").(string)
	name := d.Get("name").(string)
	useSlug := true
	if v, ok := d.GetOk("use_slug"); ok {
		useSlug = v.(bool)
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

			// Check if the name contains the prefix
			if !strings.Contains(nameLower, prefixLower) {
				return fmt.Errorf("resource name %s must contain the slug '%s' for resource type '%s'", name, expectedPrefix, resourceType)
			}

			// Validate prefix placement
			slugIndex := strings.Index(nameLower, prefixLower)
			if slugIndex > 0 {
				prevChar := rune(name[slugIndex-1])
				if !strings.ContainsRune("-_.", prevChar) {
					return fmt.Errorf("resource name %s has incorrectly placed slug '%s' - should be at start or after separator (-, _, or .)", name, expectedPrefix)
				}
			}

			// Validate that prefix is followed by a separator
			slugEndIndex := slugIndex + len(prefixLower)
			if slugEndIndex < len(name) {
				nextChar := rune(name[slugEndIndex])
				if !strings.ContainsRune("-_.", nextChar) {
					return fmt.Errorf("resource name %s must have a separator (-, _, or .) after the slug '%s'", name, expectedPrefix)
				}
			}
		}

		// Additional validation for regex pattern
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
func BaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The base name to use for the Azure resource.",
			ValidateFunc: func(i interface{}, k string) (ws []string, es []error) {
				v, ok := i.(string)
				if !ok {
					es = append(es, fmt.Errorf("expected type of %s to be string", k))
					return
				}
				if v == "" {
					es = append(es, fmt.Errorf("name cannot be empty"))
					return
				}
				if strings.TrimSpace(v) == "" {
					es = append(es, fmt.Errorf("name cannot be only whitespace"))
					return
				}

				// Get resource type from schema
				d := schema.TestResourceDataRaw(nil, BaseSchema(), map[string]interface{}{
					"name":          v,
					"resource_type": "",
				})
				if rt, ok := d.GetOk("resource_type"); ok {
					if err := ValidateResourceNameWithSlug(rt.(string), v); err != nil {
						es = append(es, err)
					}
				}
				return
			},
		},
		"resource_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: func(i interface{}, k string) (ws []string, es []error) {
				v, ok := i.(string)
				if !ok {
					es = append(es, fmt.Errorf("expected type of %s to be string", k))
					return
				}
				if v == "" {
					es = append(es, fmt.Errorf("resource_type cannot be empty"))
					return
				}
				if strings.TrimSpace(v) == "" {
					es = append(es, fmt.Errorf("resource_type cannot be only whitespace"))
					return
				}

				resource, err := models.ValidateResourceType(v)
				if err != nil {
					es = append(es, fmt.Errorf("invalid resource type %s: %v", v, err))
					return
				}

				// Validate regex patterns
				if resource.ValidationRegExp != "" {
					pattern, err := regexp.Compile(resource.ValidationRegExp)
					if err != nil {
						es = append(es, fmt.Errorf("invalid validation regex pattern for resource type %s: %v", v, err))
						return
					}
					if !pattern.MatchString(v) {
						es = append(es, fmt.Errorf("resource type %s does not match required pattern: %s", v, resource.ValidationRegExp))
						return
					}
					ws = append(ws, fmt.Sprintf("Resource type %s will be validated against pattern: %s", v, resource.ValidationRegExp))
				}
				if resource.RegEx != "" {
					pattern, err := regexp.Compile(resource.RegEx)
					if err != nil {
						es = append(es, fmt.Errorf("invalid regex pattern for resource type %s: %v", v, err))
						return
					}
					if !pattern.MatchString(v) {
						es = append(es, fmt.Errorf("resource type %s does not match required pattern: %s", v, resource.RegEx))
						return
					}
					ws = append(ws, fmt.Sprintf("Resource type %s has additional regex validation: %s", v, resource.RegEx))
				}

				// Validate length constraints
				if resource.MinLength < 0 {
					es = append(es, fmt.Errorf("invalid minimum length for resource type %s: %d", v, resource.MinLength))
					return
				}
				if resource.MaxLength > 0 && resource.MaxLength < resource.MinLength {
					es = append(es, fmt.Errorf("invalid maximum length for resource type %s: max=%d is less than min=%d", v, resource.MaxLength, resource.MinLength))
					return
				}
				if len(v) < resource.MinLength {
					es = append(es, fmt.Errorf("resource type %s is too short: min=%d, got=%d", v, resource.MinLength, len(v)))
					return
				}
				if resource.MaxLength > 0 && len(v) > resource.MaxLength {
					es = append(es, fmt.Errorf("resource type %s is too long: max=%d, got=%d", v, resource.MaxLength, len(v)))
					return
				}
				if resource.MinLength > 0 || resource.MaxLength > 0 {
					ws = append(ws, fmt.Sprintf("Resource type %s length constraints: min=%d, max=%d", v, resource.MinLength, resource.MaxLength))
				}

				// Validate CAF prefix
				if resource.CafPrefix != "" {
					if !strings.HasPrefix(strings.ToLower(v), strings.ToLower(resource.CafPrefix)) {
						es = append(es, fmt.Errorf("resource type %s requires prefix %s", v, resource.CafPrefix))
						return
					}
					ws = append(ws, fmt.Sprintf("Resource type %s requires prefix: %s", v, resource.CafPrefix))
				}

				// Add validation warnings for other constraints
				if resource.LowerCase {
					if v != strings.ToLower(v) {
						es = append(es, fmt.Errorf("resource type %s requires lowercase characters only", v))
						return
					}
					ws = append(ws, fmt.Sprintf("Resource type %s requires lowercase characters only", v))
				}
				if resource.Dashes {
					if !strings.Contains(v, "-") {
						ws = append(ws, fmt.Sprintf("Resource type %s allows dashes in the name", v))
					}
				}
				if resource.Scope != "" {
					ws = append(ws, fmt.Sprintf("Resource type %s has scope: %s", v, resource.Scope))
				}
				return
			},
			Description: "The type of Azure resource to create a name for. Must be one of the supported resource types.",
		},
		"prefixes": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"suffixes": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"random_length": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"random_seed": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
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
		"use_slug": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"result": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"results": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"random_string": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
