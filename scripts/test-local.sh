set -e

export PATH=/usr/lib/go-1.20/bin:$PATH

echo "Running unit tests..."
go test ./azurecaf -v -coverprofile=coverage.out
go tool cover -func=coverage.out

echo "Building provider..."
go build -o ./terraform-provider-azurecaf

echo "Running integration tests..."
cd ./examples && terraform init && terraform plan && cd ..

echo "Generating and running e2e tests..."
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/1.2.10/linux_amd64
cp terraform-provider-azurecaf ~/.terraform.d/plugins/registry.terraform.io/aztfmod/azurecaf/1.2.10/linux_amd64/

cd examples/e2e
go run generator.go
terraform init -upgrade
terraform validate
chmod +x ./validate_results.sh
./validate_results.sh
cd ../..

echo "All tests passed!"
