# Development Guide

This guide helps you get started developing the Azure CAF Terraform Provider.

## Prerequisites

1. [Go](https://golang.org/doc/install) 1.21 or later
2. [Terraform](https://www.terraform.io/downloads.html) 1.0 or later

## Building and Testing Locally

### Building the Provider

1. Clone the repository:
```bash
git clone https://github.com/arnaudlh/terraform-provider-azurecaf.git
cd terraform-provider-azurecaf
```

2. Build the provider:
```bash
go build -v ./...
```

3. Install the provider locally:
```bash
# For Linux/macOS:
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/arnaudlh/azurecaf/$(git describe --tags --abbrev=0)/$(go env GOOS)_$(go env GOARCH)
cp terraform-provider-azurecaf ~/.terraform.d/plugins/registry.terraform.io/arnaudlh/azurecaf/$(git describe --tags --abbrev=0)/$(go env GOOS)_$(go env GOARCH)/

# For Windows (PowerShell):
$version = $(git describe --tags --abbrev=0)
$arch = "$(go env GOOS)_$(go env GOARCH)"
New-Item -Path "$env:APPDATA\terraform.d\plugins\registry.terraform.io\arnaudlh\azurecaf\$version\$arch" -ItemType Directory -Force
Copy-Item "terraform-provider-azurecaf.exe" -Destination "$env:APPDATA\terraform.d\plugins\registry.terraform.io\arnaudlh\azurecaf\$version\$arch\"
```

### Running Tests

1. Run unit tests:
```bash
go test -v -tags=unit ./...
```

2. Run tests with coverage:
```bash
go test -v -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt
```

### Testing with Terraform

1. Create a test configuration:
```bash
mkdir test && cd test
cat > main.tf << 'EOF'
terraform {
  required_providers {
    azurecaf = {
      source  = "arnaudlh/azurecaf"
      version = "~> 2.0.0-preview4"
    }
  }
}

provider "azurecaf" {}

# Test resource group naming
resource "azurecaf_name" "rg" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  clean_input   = true
}

# Test storage account naming
resource "azurecaf_name" "st" {
  name          = "data"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
  random_length = 5
}

output "resource_group_name" {
  value = azurecaf_name.rg.result
}

output "storage_account_name" {
  value = azurecaf_name.st.result
}
EOF
```

2. Initialize and test:
```bash
terraform init
terraform plan
terraform apply
```

3. Verify the generated names follow Azure naming conventions:
```bash
terraform output
```

4. Clean up:
```bash
terraform destroy
```

## Resource Definition Schema

The provider uses `resourceDefinition.json` to define naming rules. Example structure:

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

## Code Organization

```
azurecaf/
├── internal/
│   ├── models/      # Resource definition models
│   └── schemas/     # Schema definitions
├── provider.go      # Provider entry point
├── resource_name.go # Name resource implementation
└── data_name.go    # Data source implementation
```

## Best Practices

1. Run tests before committing:
```bash
go test -v -tags=unit ./...
go vet ./...
go fmt ./...
```

2. Ensure test coverage remains above 75%:
```bash
go test -v -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt
```

3. Follow Azure naming conventions when adding new resource types
