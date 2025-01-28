# Development Guide

This guide helps you get started developing the Azure CAF Terraform Provider.

## Prerequisites

1. [Go](https://golang.org/doc/install) 1.21 or later
2. [Terraform](https://www.terraform.io/downloads.html) 1.0 or later

## Environment Setup

1. Install required Go tools:
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

2. Set environment variables:
```bash
export GO111MODULE=on
export GOFLAGS=-mod=vendor
```

## Building the Provider Locally

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
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/arnaudlh/azurecaf/0.0.1/$(go env GOOS)_$(go env GOARCH)
cp terraform-provider-azurecaf ~/.terraform.d/plugins/registry.terraform.io/arnaudlh/azurecaf/0.0.1/$(go env GOOS)_$(go env GOARCH)/

# For Windows (PowerShell):
# New-Item -Path "$env:APPDATA\terraform.d\plugins\registry.terraform.io\arnaudlh\azurecaf\0.0.1\$(go env GOOS)_$(go env GOARCH)" -ItemType Directory -Force
# Copy-Item "terraform-provider-azurecaf.exe" -Destination "$env:APPDATA\terraform.d\plugins\registry.terraform.io\arnaudlh\azurecaf\0.0.1\$(go env GOOS)_$(go env GOARCH)\"
```

## Testing

### Running Unit Tests

To run the unit tests:
```bash
go test -v -tags=unit ./...
```

To run tests with coverage:
```bash
go test -v -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt
```

### Testing with Terraform

1. Create a test directory and configuration:
```bash
mkdir test && cd test
cat > main.tf << 'EOF'
terraform {
  required_providers {
    azurecaf = {
      source  = "arnaudlh/azurecaf"
      version = "0.0.1"
    }
  }
}

provider "azurecaf" {}

resource "azurecaf_name" "test" {
  name          = "test"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  random_length = 5
  clean_input   = true
}

output "result" {
  value = azurecaf_name.test.result
}
EOF
```

2. Initialize and test:
```bash
terraform init
terraform plan
terraform apply -auto-approve
terraform destroy -auto-approve
```

## Debugging

### Common Issues

1. **Build Failures**
   - Ensure Go version matches requirements (1.21+)
   - Run `go mod tidy` to resolve dependencies
   - Check for syntax errors with `go vet ./...`

2. **Test Failures**
   - Use `-v` flag for verbose output
   - Run specific tests with `go test -v -tags=unit ./... -run TestName`
   - Check test coverage with `go tool cover -html=coverage.txt`

3. **Provider Installation Issues**
   - Verify correct plugin directory structure
   - Check file permissions
   - Clear Terraform plugin cache: `rm -rf ~/.terraform.d/plugin-cache`

### Logging

Enable debug logging by setting environment variables:
```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log
```

## Resource Definitions

The provider uses `resourceDefinition.json` to define naming rules for Azure resources. Key components:

- `name`: Azure resource type
- `min_length`/`max_length`: Name length constraints
- `validation_regex`: Name validation pattern
- `scope`: Resource naming scope (global/resourceGroup/subscription)
- `slug`: Short identifier
- `dashes`/`lowercase`: Formatting rules

## Code Organization

- `azurecaf/` - Main provider code
  - `internal/` - Internal packages
    - `models/` - Resource definition models
    - `schemas/` - Schema definitions
  - `provider.go` - Provider entry point
  - `resource_name.go` - Name resource implementation

## Best Practices

1. **Code Quality**
   - Run `go fmt ./...` before committing
   - Maintain test coverage above 75%
   - Use meaningful variable names
   - Add comments for complex logic

2. **Testing**
   - Write unit tests for new features
   - Test edge cases
   - Validate resource definitions
   - Use build tags appropriately

3. **Resource Definitions**
   - Follow Azure naming conventions
   - Include all required fields
   - Validate regex patterns
   - Document scope requirements

## Code Organization

- `azurecaf/` - Main provider code
  - `internal/` - Internal packages
    - `models/` - Resource definition models
    - `schemas/` - Schema definitions and state migrations
  - `provider.go` - Provider entry point
  - `resource_name.go` - Name resource implementation
  - `data_name.go` - Data source implementation

## Running Tests in CI

The CI pipeline runs tests for Go versions 1.21 and 1.22. To ensure your changes pass CI:
1. Write unit tests for new features
2. Ensure test coverage is above 75%
3. Run tests locally before pushing
