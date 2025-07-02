default: help

dev_container:
	go generate
	go fmt
	go build -o ~/.terraform.d/plugins/linux_amd64/terraform-provider-azurecaf

build:	## Build the project and run unit tests
	go generate
	go fmt ./...
	go build -o ./terraform-provider-azurecaf
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -cover ./...

unittest: 	## Run unit tests without coverage
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test ./...
	tfproviderlint ./...

test_coverage: 	## Run tests with coverage reporting
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -cover ./...

test_coverage_html: 	## Generate HTML coverage report
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at: coverage.html"

test_coverage_specific: ## Run coverage-focused tests specifically
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="Test.*" -coverprofile=coverage.out

test_integration: 	## Run integration tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" TF_ACC=1 go test -v ./azurecaf/... -run="TestAcc"

test_data_sources: 	## Run data source integration tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" TF_ACC=1 go test -v ./azurecaf/... -run="TestAccDataSourcesIntegration"

test_error_handling: 	## Run error handling integration tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" TF_ACC=1 go test -v ./azurecaf/... -run="TestAccErrorHandling"

test_resource_naming: ## Run naming convention tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="TestAcc.*NamingConvention" -coverprofile=naming_coverage.out ./...
	go tool cover -html=naming_coverage.out -o naming_coverage.html
	@echo "Naming coverage report generated at: naming_coverage.html"

test_all_resources: 	## Test ALL resource types (comprehensive integration test)
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="TestAcc_AllResourceTypes" -timeout=30m

test_resource_coverage: 	## Analyze test coverage for all resource types
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="TestResourceCoverage" -timeout=10m

test_resource_definition_completeness: 	## Validate all resource definitions are complete
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="TestResourceDefinitionCompleteness"

test_all: unittest test_integration	## Run all tests (unit and integration)

test_ci: unittest test_coverage test_resource_definitions test_resource_matrix test_resource_coverage	## Run comprehensive CI tests (unit, coverage, resource validation, matrix testing)

test_ci_fast: unittest test_coverage test_resource_definitions	## Run fast CI tests (unit, coverage, resource validation only)

test_ci_complete: test_ci test_integration test_all_resources	## Run complete CI tests including integration tests

test_resource_definitions: test_resource_definition_completeness	## Validate all resource definitions are complete

test_resource_matrix: 	## Test resources by category and validate constraints
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="TestResourceMatrix|TestResourceConstraints"

test_complete: test_all test_all_resources test_resource_coverage	## Complete test suite including all resource types

test: ## Run terraform examples with local provider
	# First build the provider
	go build -o ./terraform-provider-azurecaf
	
	# Create script to set up and run the examples
	@echo '#!/bin/bash' > run_examples.sh
	@echo 'GOOS=$$(go env GOOS)' >> run_examples.sh
	@echo 'GOARCH=$$(go env GOARCH)' >> run_examples.sh
	@echo 'LOCAL_PLUGIN_DIR=~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/$${GOOS}_$${GOARCH}' >> run_examples.sh
	@echo 'echo "Using local plugin directory: $$LOCAL_PLUGIN_DIR"' >> run_examples.sh
	@echo 'mkdir -p $$LOCAL_PLUGIN_DIR' >> run_examples.sh
	@echo 'cp ./terraform-provider-azurecaf $$LOCAL_PLUGIN_DIR/' >> run_examples.sh
	@echo '' >> run_examples.sh
	@echo '# Create development override file for examples' >> run_examples.sh
	@echo 'cat > examples/terraform.rc << EOF' >> run_examples.sh
	@echo 'provider_installation {' >> run_examples.sh
	@echo '  dev_overrides {' >> run_examples.sh
	@echo '    "aztfmod.com/arnaudlh/azurecaf" = "$${HOME}/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/$${GOOS}_$${GOARCH}"' >> run_examples.sh
	@echo '  }' >> run_examples.sh
	@echo '  direct {}' >> run_examples.sh
	@echo '}' >> run_examples.sh
	@echo 'EOF' >> run_examples.sh
	@echo '' >> run_examples.sh
	@echo '# Run terraform in examples directory using the local config' >> run_examples.sh
	@echo 'cd ./examples && TF_CLI_CONFIG_FILE=terraform.rc terraform init -upgrade && terraform plan && terraform apply -auto-approve' >> run_examples.sh
	
	# Make the script executable and run it
	@chmod +x run_examples.sh
	@./run_examples.sh
	@rm run_examples.sh

generate_resource_table:  	## Generate resource table (output only)
	cat resourceDefinition.json | jq -r '.[] | "| \(.name)| \(.slug)| \(.min_length)| \(.max_length)| \(.lowercase)| \(.validation_regex)|"'

# End-to-End Testing Targets

test_e2e: 	## Run complete end-to-end test suite
	@echo "Running complete e2e test suite..."
	cd e2e && CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v

test_e2e_quick: 	## Run quick e2e tests (basic scenarios only)
	@echo "Running quick e2e tests..."
	cd e2e && CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v -run TestE2EBasic

test_e2e_data_source: 	## Run e2e data source tests
	@echo "Running e2e data source tests..."
	cd e2e && CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v -run TestE2EDataSource

test_e2e_naming: 	## Run e2e naming convention tests
	@echo "Running e2e naming convention tests..."
	cd e2e && CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v -run TestE2ENamingConventions

test_e2e_multiple_types: 	## Run e2e multiple resource types tests
	@echo "Running e2e multiple resource types tests..."
	cd e2e && CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v -run TestE2EMultipleResourceTypes

test_e2e_import: 	## Run e2e import functionality tests
	@echo "Running e2e import functionality tests..."
	cd e2e && CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v -run TestE2EImportFunctionality

test_e2e_verbose: 	## Run e2e tests with verbose output
	@echo "Running e2e tests with verbose output..."
	cd e2e && CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v

# Enhanced testing targets that include e2e tests

test_complete_with_e2e: test_complete test_e2e	## Run complete test suite including e2e tests

test_ci_with_e2e: test_ci test_e2e_quick	## Run CI tests including quick e2e tests

# Performance and Quality Targets

benchmark: 	## Run benchmark tests
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem -run=^$ ./...

benchmark_verbose: 	## Run benchmark tests with verbose output
	@echo "Running benchmark tests with verbose output..."
	go test -bench=. -benchmem -benchtime=10s -run=^$ -v ./...

profile_mem: 	## Run tests with memory profiling
	@echo "Running tests with memory profiling..."
	go test -memprofile=mem.prof -run=TestAcc ./azurecaf/...
	@if [ -f mem.prof ]; then echo "Memory profile saved to mem.prof"; fi

profile_cpu: 	## Run tests with CPU profiling
	@echo "Running tests with CPU profiling..."
	go test -cpuprofile=cpu.prof -run=TestAcc ./azurecaf/...
	@if [ -f cpu.prof ]; then echo "CPU profile saved to cpu.prof"; fi

profile_analyze: 	## Analyze existing profiles
	@echo "Analyzing profiles..."
	@if [ -f mem.prof ]; then go tool pprof -text mem.prof > memory-analysis.txt && echo "Memory analysis saved to memory-analysis.txt"; fi
	@if [ -f cpu.prof ]; then go tool pprof -text cpu.prof > cpu-analysis.txt && echo "CPU analysis saved to cpu-analysis.txt"; fi

lint: 	## Run comprehensive linting
	@echo "Running comprehensive linting..."
	go vet ./...
	tfproviderlint ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else echo "golangci-lint not installed, skipping"; fi

format: 	## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then goimports -w .; else echo "goimports not installed, using go fmt only"; fi

clean: 	## Clean build artifacts and test files
	@echo "Cleaning build artifacts..."
	rm -f terraform-provider-azurecaf
	rm -f *.prof
	rm -f *.out
	rm -f *.html
	rm -f *.txt
	rm -f coverage.*
	rm -rf .terraform/
	rm -f terraform.tfstate*
	rm -f .terraform.lock.hcl

clean_all: clean 	## Clean all artifacts including temporary documentation
	@echo "Cleaning all temporary files..."
	rm -f *SUMMARY*.md
	rm -f *GUIDE*.md
	rm -f *PLAN*.md
	rm -f *ANALYSIS*.md
	@echo "Cleanup complete"

security_scan: 	## Run security scanning tools
	@echo "Running security scans..."
	@if command -v gosec >/dev/null 2>&1; then gosec ./...; else echo "gosec not installed, skipping"; fi
	@if command -v govulncheck >/dev/null 2>&1; then govulncheck ./...; else echo "govulncheck not installed, skipping"; fi

dependency_check: 	## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	go list -u -m all
	@if command -v go-mod-outdated >/dev/null 2>&1; then go list -u -m -json all | go-mod-outdated -update -direct; else echo "go-mod-outdated not installed, using basic check"; fi

# Quality Assurance Targets

qa: lint format test_coverage 	## Run comprehensive quality assurance checks

qa_full: qa security_scan dependency_check benchmark 	## Run full quality assurance suite

ci_local: build lint unittest test_coverage 	## Run CI-equivalent tests locally

ci_full: ci_local test_integration test_e2e_quick 	## Run full CI suite locally

# Development Helpers

dev_setup: 	## Set up development environment
	@echo "Setting up development environment..."
	go mod download
	@if ! command -v tfproviderlint >/dev/null 2>&1; then go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest; fi
	@if ! command -v golangci-lint >/dev/null 2>&1; then echo "Consider installing golangci-lint for enhanced linting"; fi
	@if ! command -v gosec >/dev/null 2>&1; then echo "Consider installing gosec for security scanning"; fi

dev_build: 	## Build for development (with local installation)
	go generate
	go fmt ./...
	go build -o ~/.terraform.d/plugins/terraform-provider-azurecaf
	@echo "Provider built and installed for local development"

watch_tests: 	## Watch for file changes and run tests
	@echo "Watching for changes... (requires 'entr' tool)"
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -c make unittest; \
	else \
		echo "Install 'entr' tool to use watch functionality"; \
		echo "Example: apt-get install entr (Ubuntu) or brew install entr (macOS)"; \
	fi

# Help improvement
.PHONY: help
help:  ## Display help (improved version)
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m\033[34mTerraform Provider AzureCAF - Available Targets\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "\033[1m\033[33mCommon workflows:\033[0m"
	@echo "  \033[36mmake ci_local\033[0m          - Run CI tests locally"
	@echo "  \033[36mmake qa\033[0m                - Run quality assurance checks"
	@echo "  \033[36mmake dev_setup\033[0m         - Set up development environment"
	@echo "  \033[36mmake clean\033[0m             - Clean all build artifacts"
	@echo ""

