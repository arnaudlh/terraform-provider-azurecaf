package schemas

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	models "github.com/aztfmod/terraform-provider-azurecaf/azurecaf/models"
)

func V4_Schema() map[string]*schema.Schema {
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
