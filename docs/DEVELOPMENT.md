# Development Guide

This guide helps you get started developing the Azure CAF Terraform Provider.

## Prerequisites

1. [Go](https://golang.org/doc/install) 1.21 or later
2. [Terraform](https://www.terraform.io/downloads.html) 1.0 or later

## Building and Testing Locally

### Building the Provider

1. Clone the repository:
```bash
git clone https://github.com/aztfmod/terraform-provider-azurecaf.git
cd terraform-provider-azurecaf
```

2. Build the provider:
```bash
go build -v ./...
```

### Running Tests

Run unit tests:
```bash
go test -v ./...
```

Run unit tests with coverage:
```bash
go test -v -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt
```

### Testing with Terraform Locally

1. Create a development override configuration:
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/arnaudlh/azurecaf/2.0.0-preview4/linux_amd64
cp terraform-provider-azurecaf ~/.terraform.d/plugins/registry.terraform.io/arnaudlh/azurecaf/2.0.0-preview4/linux_amd64/
```

2. Create a test configuration:
```bash
mkdir test && cd test
cat > main.tf << 'EOF'
terraform {
  required_providers {
    azurecaf = {
      source  = "arnaudlh/azurecaf"
      version = "2.0.0-preview4"
    }
  }
}

provider "azurecaf" {}

# Test resource group naming using resource
resource "azurecaf_name" "rg" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  clean_input   = true
}

# Test storage account naming using resource
resource "azurecaf_name" "st" {
  name          = "data"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
  random_length = 5
  clean_input   = true
  separator     = ""
  use_slug      = true
}

# Test resource group naming using data source
data "azurecaf_name" "rg_data" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  clean_input   = true
}

# Test storage account naming using data source
data "azurecaf_name" "st_data" {
  name          = "data"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
  random_length = 5
  clean_input   = true
  separator     = ""
  use_slug      = true
}

output "resource_group_name_resource" {
  value = azurecaf_name.rg.result
}

output "storage_account_name_resource" {
  value = azurecaf_name.st.result
}

output "resource_group_name_data" {
  value = data.azurecaf_name.rg_data.result
}

output "storage_account_name_data" {
  value = data.azurecaf_name.st_data.result
}
EOF
```

3. Test the configuration:
```bash
# Initialize Terraform
terraform init

# Plan changes
terraform plan

# Apply configuration
terraform apply

# Verify outputs
terraform output

# Expected output format:
# resource_group_name = "dev-myapp-001"
# storage_account_name = "devstdataxxx" (where xxx is a random string)

# Clean up
terraform destroy
```

## Project Structure

```
.
├── azurecaf/
│   ├── internal/
│   │   ├── models/      # Resource definition models
│   │   └── schemas/     # Schema definitions
│   ├── provider.go      # Provider entry point
│   ├── resource_name.go # Name resource implementation
│   └── data_name.go    # Data source implementation
├── docs/               # Documentation
├── examples/          # Example configurations
└── resourceDefinition.json  # Resource naming rules
```

## Development Workflow

1. Make your changes in a feature branch
2. Run tests:
```bash
go test -v ./...
```

3. Run linter:
```bash
golangci-lint run
```

4. Build and test locally:
```bash
go build -v ./...
# Follow the local testing steps above
```

5. Ensure test coverage is above 75%:
```bash
go test -v -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt
```

## Resource Definition Schema

The provider uses `resourceDefinition.json` to define naming rules:

```json
{
  "name": "azurerm_resource_group",
  "min_length": 1,
  "max_length": 90,
  "validation_regex": "^[a-zA-Z0-9-_()]+$",
  "scope": "resourceGroup",
  "slug": "rg",
  "dashes": true
}
```

## Debugging

Enable debug logging:
```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log
```

## Common Issues and Solutions

1. Provider installation errors:
   - Ensure the provider binary is in the correct plugins directory
   - Verify the version in your Terraform configuration
   - Check file permissions

2. Storage account naming failures:
   - Use lowercase letters and numbers only
   - Maximum length: 24 characters
   - Use `separator = ""` and `use_slug = true`

3. Resource group naming issues:
   - Allows dashes and underscores
   - Length: 1-90 characters
   - Verify prefix/suffix validity
