package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

const resourceTemplate = `
# Test resource name for {{ .Name }}
# Standard naming test
resource "azurecaf_name" "test_{{ .Slug }}" {
  name          = "test{{ .Slug }}"
  resource_type = "{{ .Name }}"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_{{ .Slug }}_special" {
  name          = "test{{ .Slug }}!@#$%"
  resource_type = "{{ .Name }}"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_{{ .Slug }}_max" {
  name          = "{{ .Slug }}verylongname"
  resource_type = "{{ .Name }}"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_{{ .Slug }}_min" {
  name          = "{{ if .Lowercase }}a{{ else }}A{{ end }}"
  resource_type = "{{ .Name }}"
  random_length = {{ .MinLength }}
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_{{ .Slug }}_ds" {
  name          = "test{{ .Slug }}"
  resource_type = "{{ .Name }}"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_{{ .Slug }}_standard" {
  value = can(regex({{ .ValidationRegex }}, azurecaf_name.test_{{ .Slug }}.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_{{ .Slug }}.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_{{ .Slug }}_special" {
  value = can(regex({{ .ValidationRegex }}, azurecaf_name.test_{{ .Slug }}_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_{{ .Slug }}_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_{{ .Slug }}_max" {
  value = length(azurecaf_name.test_{{ .Slug }}_max.result) <= {{ .MaxLength }} ? "PASS" : "FAIL: ${azurecaf_name.test_{{ .Slug }}_max.result} exceeds maximum length {{ .MaxLength }}"
}

# Verify minimum length is enforced
output "test_{{ .Slug }}_min" {
  value = length(azurecaf_name.test_{{ .Slug }}_min.result) >= {{ .MinLength }} ? "PASS" : "FAIL: ${azurecaf_name.test_{{ .Slug }}_min.result} below minimum length {{ .MinLength }}"
}

# Verify data source consistency
output "test_{{ .Slug }}_ds" {
  value = azurecaf_name.test_{{ .Slug }}.result == data.azurecaf_name.test_{{ .Slug }}_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_{{ .Slug }}.result} != ${data.azurecaf_name.test_{{ .Slug }}_ds.result}"
}

# Verify case sensitivity
output "test_{{ .Slug }}_case" {
  value = {{ if .Lowercase }}can(regex("^[a-z0-9-_.]*$", azurecaf_name.test_{{ .Slug }}.result)){{ else }}true{{ end }} ? "PASS" : "FAIL: Case sensitivity rules not followed"
}`

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
	
	tmpl, err := template.New("resource").Parse(resourceTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	resourceCount := 0
	fileCount := 0
	var currentFile *os.File
	var currentBuilder strings.Builder

	for _, rd := range resourceDefinitions {
		if rd.Name == "" || rd.Slug == "" {
			continue
		}

		if rd.ValidationRegex == "" {
			rd.ValidationRegex = `"^[a-zA-Z0-9-]+$"`
		}

		if resourceCount%5 == 0 {
			if currentFile != nil {
				_, err = currentFile.WriteString(currentBuilder.String())
				if err != nil {
					log.Printf("Error writing to file: %v", err)
				}
				currentFile.Close()
			}

			fileCount++
			fileName := fmt.Sprintf("resources/resources_%d.tf", fileCount)
			currentFile, err = os.Create(fileName)
			if err != nil {
				log.Fatalf("Error creating file %s: %v", fileName, err)
			}
			
			currentBuilder.Reset()
		}

		var resourceBlock strings.Builder
		err = tmpl.Execute(&resourceBlock, rd)
		if err != nil {
			log.Printf("Error executing template for %s: %v", rd.Name, err)
			continue
		}

		currentBuilder.WriteString(resourceBlock.String())
		resourceCount++

		if resourceCount >= 20 {
			break
		}
	}

	if currentFile != nil {
		_, err = currentFile.WriteString(currentBuilder.String())
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}
		currentFile.Close()
	}

	fmt.Printf("E2E test files generated successfully! Created %d files with %d resources.\n", fileCount, resourceCount)
}
