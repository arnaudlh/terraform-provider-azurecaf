package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "text/template"
)

type ResourceDefinition struct {
    Name            string `json:"name"`
    MinLength       int    `json:"min_length"`
    MaxLength       int    `json:"max_length"`
    ValidationRegex string `json:"validation_regex"`
    Scope           string `json:"scope"`
    Slug            string `json:"slug"`
    Dashes          bool   `json:"dashes"`
    Lowercase       bool   `json:"lowercase"`
    Regex           string `json:"regex"`
}

const testTemplate = `# E2E test for {{ .Name }}
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}

provider "azurecaf" {
}

# Test resource name
resource "azurecaf_name" "test_{{ .Slug }}" {
  name          = "test{{ .Slug }}"
  resource_type = "{{ .Name }}"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test data source with same parameters
data "azurecaf_name" "test_{{ .Slug }}_ds" {
  name          = "test{{ .Slug }}"
  resource_type = "{{ .Name }}"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify that resource and data source produce the same result
output "test_result" {
  value = azurecaf_name.test_{{ .Slug }}.result == data.azurecaf_name.test_{{ .Slug }}_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_{{ .Slug }}.result} != ${data.azurecaf_name.test_{{ .Slug }}_ds.result}"
}

# Verify regex pattern matches
output "regex_validation" {
  value = can(regex({{ .ValidationRegex }}, azurecaf_name.test_{{ .Slug }}.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_{{ .Slug }}.result} does not match regex pattern"
}

# Verify slug placement
output "slug_placement" {
  value = can(regex(".*{{ .Slug }}.*", azurecaf_name.test_{{ .Slug }}.result)) ? "PASS" : "FAIL: Slug '{{ .Slug }}' not found in result ${azurecaf_name.test_{{ .Slug }}.result}"
}
`

func main() {
    jsonFile, err := os.Open("../../resourceDefinition.json")
    if err != nil {
        log.Fatalf("Error opening resourceDefinition.json: %v", err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)
    var resourceDefinitions []ResourceDefinition
    json.Unmarshal(byteValue, &resourceDefinitions)

    os.MkdirAll("resources", 0755)

    tmpl, err := template.New("test").Parse(testTemplate)
    if err != nil {
        log.Fatalf("Error parsing template: %v", err)
    }

    for _, rd := range resourceDefinitions {
        if rd.Name == "" || rd.Slug == "" {
            continue
        }

        fileName := fmt.Sprintf("resources/%s.tf", rd.Slug)
        file, err := os.Create(fileName)
        if err != nil {
            log.Printf("Error creating file for %s: %v", rd.Name, err)
            continue
        }

        err = tmpl.Execute(file, rd)
        if err != nil {
            log.Printf("Error executing template for %s: %v", rd.Name, err)
        }
        file.Close()
    }

    fmt.Println("E2E test files generated successfully!")
}
