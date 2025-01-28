//go:build unit
// +build unit

package completness

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"testing"
)

func TestValidateResourceDefinition(t *testing.T) {
	tests := []struct {
		name       string
		resources  []string
		wantErr    bool
		errMessage string
	}{
		{
			name:      "valid resources",
			resources: []string{"azurerm_resource_group", "azurerm_storage_account"},
			wantErr:   false,
		},
		{
			name:       "invalid resource",
			resources:  []string{"invalid_resource"},
			wantErr:    true,
			errMessage: "resource type invalid_resource not found in the resource definition file",
		},
		{
			name:      "empty resource list",
			resources: []string{},
			wantErr:   false,
		},
		{
			name:       "mixed valid and invalid resources",
			resources:  []string{"azurerm_resource_group", "invalid_resource"},
			wantErr:    true,
			errMessage: "resource type invalid_resource not found in the resource definition file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResourceDefinition(tt.resources)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResourceDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMessage {
				t.Errorf("ValidateResourceDefinition() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}

func TestGetResourceDefinition(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "get resource definition",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResourceDefinition()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResourceDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("GetResourceDefinition() returned nil slice")
				}
				if len(got) == 0 {
					t.Error("GetResourceDefinition() returned empty slice")
				}
			}
		})
	}
}

func TestGetResourceMap(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "get resource map",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResourceMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResourceMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("GetResourceMap() returned nil map")
				}
				if len(got) == 0 {
					t.Error("GetResourceMap() returned empty map")
				}
				// Test specific resource existence
				if resource, exists := got["azurerm_resource_group"]; !exists {
					t.Error("Expected resource 'azurerm_resource_group' not found in map")
				} else {
					if resource.ResourceTypeName != "azurerm_resource_group" {
						t.Errorf("Resource name mismatch, got %s, want azurerm_resource_group", resource.ResourceTypeName)
					}
				}
			}
		})
	}
}

func TestReadLines(t *testing.T) {
	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test_resources_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data
	testData := []string{
		"azurerm_resource_group",
		"azurerm_storage_account",
		"azurerm_virtual_network",
	}
	content := []byte(fmt.Sprintf("%s\n%s\n%s", testData[0], testData[1], testData[2]))
	if err := os.WriteFile(tmpFile.Name(), content, 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Test reading lines
	got, err := readLines(tmpFile.Name())
	if err != nil {
		t.Fatalf("readLines() error = %v", err)
	}

	// Verify results
	if len(got) != len(testData) {
		t.Errorf("readLines() returned %d lines, want %d", len(got), len(testData))
	}
	for i, want := range testData {
		if got[i] != want {
			t.Errorf("readLines()[%d] = %v, want %v", i, got[i], want)
		}
	}

	// Test reading non-existent file
	_, err = readLines("non_existent_file.txt")
	if err == nil {
		t.Error("readLines() expected error for non-existent file, got nil")
	}
}

func TestResourceValidation(t *testing.T) {
	testData := []string{
		"azurerm_resource_group",
		"azurerm_storage_account",
		"azurerm_virtual_network",
	}

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "test_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a minimal resourceDefinition.json for testing
	minimalJson := []byte(`[
		{
			"name": "azurerm_resource_group",
			"min_length": 1,
			"max_length": 90,
			"validation_regex": "^[a-zA-Z0-9-_()]+$",
			"scope": "resourceGroup",
			"slug": "rg",
			"dashes": true
		},
		{
			"name": "azurerm_storage_account",
			"min_length": 3,
			"max_length": 24,
			"validation_regex": "^[a-z0-9]+$",
			"scope": "global",
			"slug": "st",
			"dashes": false
		},
		{
			"name": "azurerm_virtual_network",
			"min_length": 2,
			"max_length": 64,
			"validation_regex": "^[a-zA-Z0-9][a-zA-Z0-9-._]+[a-zA-Z0-9]$",
			"scope": "resourceGroup",
			"slug": "vnet",
			"dashes": true
		}
	]`)
	if err := os.WriteFile(path.Join(tmpDir, "resourceDefinition.json"), minimalJson, 0644); err != nil {
		t.Fatalf("Failed to write resourceDefinition.json to temp dir: %v", err)
	}

	// Change to temp directory
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(origWd)

	err = ValidateResourceDefinition(testData)
	if err != nil {
		t.Errorf("ValidateResourceDefinition() failed: %v", err)
	}
}

func TestFindByName(t *testing.T) {
	testData := []ResourceStructure{
		{ResourceTypeName: "azurerm_resource_group", CafPrefix: "rg"},
		{ResourceTypeName: "azurerm_storage_account", CafPrefix: "st"},
	}

	tests := []struct {
		name      string
		slice     []ResourceStructure
		searchFor string
		wantIndex int
		wantFound bool
	}{
		{
			name:      "existing resource",
			slice:     testData,
			searchFor: "azurerm_resource_group",
			wantIndex: 0,
			wantFound: true,
		},
		{
			name:      "non-existing resource",
			slice:     testData,
			searchFor: "non_existing_resource",
			wantIndex: -1,
			wantFound: false,
		},
		{
			name:      "empty slice",
			slice:     []ResourceStructure{},
			searchFor: "azurerm_resource_group",
			wantIndex: -1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotFound := findByName(tt.slice, tt.searchFor)
			if gotIndex != tt.wantIndex {
				t.Errorf("findByName() gotIndex = %v, want %v", gotIndex, tt.wantIndex)
			}
			if gotFound != tt.wantFound {
				t.Errorf("findByName() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
