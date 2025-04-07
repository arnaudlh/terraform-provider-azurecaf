set -e

export PATH=/usr/lib/go-1.20/bin:$PATH

echo "Running unit tests with coverage..."
go test ./azurecaf -v -coverprofile=coverage.out

echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo "Coverage summary:"
go tool cover -func=coverage.out

echo "Building provider..."
go build -o ./terraform-provider-azurecaf

echo "Setting up provider for local development..."
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/1.2.28/linux_amd64
cp terraform-provider-azurecaf ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/1.2.28/linux_amd64/

echo "Running integration tests..."
cd ./examples && terraform init -reconfigure && terraform validate && cd ..

echo "Preparing e2e test directory..."
mkdir -p examples/e2e/resources
chmod +x examples/e2e/validate_results.sh

echo "Running e2e tests..."
cd examples/e2e
go run generator.go
terraform init -upgrade
terraform validate
./validate_results.sh
cd ../..

echo "All tests completed successfully!"
