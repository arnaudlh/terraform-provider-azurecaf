package models

type ResourceStructure struct {
    ResourceTypeName  string `json:"name"`
    CafPrefix        string `json:"slug,omitempty"`
    MinLength        int    `json:"min_length"`
    MaxLength        int    `json:"max_length"`
    LowerCase        bool   `json:"lowercase,omitempty"`
    RegEx           string `json:"regex,omitempty"`
    ValidationRegExp string `json:"validation_regexp,omitempty"`
    Dashes          bool   `json:"dashes"`
    Scope           string `json:"scope,omitempty"`
}

var ResourceDefinitions = map[string]*ResourceStructure{}
var ResourceMaps = map[string]string{}
