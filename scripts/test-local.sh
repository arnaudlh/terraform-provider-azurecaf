set -e

echo "Running unit tests..."
make unittest

echo "Running integration tests..."
make test

echo "Generating and running e2e tests..."
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/1.2.10/linux_amd64
cp terraform-provider-azurecaf ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/1.2.10/linux_amd64/

cd examples/e2e
go run generator.go
terraform init
terraform validate
./validate_results.sh
cd ../..

echo "All tests passed!"
