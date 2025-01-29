#!/bin/bash
set -e

# Build and install the provider locally
go build -o terraform-provider-azurecaf
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/2.0.0-preview5/linux_amd64/
mv terraform-provider-azurecaf ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/2.0.0-preview5/linux_amd64/

# Run Terraform
terraform init
terraform apply -auto-approve

# Validate outputs
resource_name=$(terraform output -raw resource_name)
data_name=$(terraform output -raw data_source_name)

if [ "$resource_name" != "$data_name" ]; then
  echo "Error: Resource name ($resource_name) does not match data source name ($data_name)"
  exit 1
fi

echo "Success: Resource and data source names match"

# Clean up
terraform destroy -auto-approve
