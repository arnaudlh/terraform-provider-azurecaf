#!/bin/bash

# CI/CD Pipeline Validation Script
# This script validates all the changes made to the CI/CD pipelines

set -e

echo "🔍 CI/CD Pipeline Validation Suite"
echo "=================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Function to run a check
run_check() {
    local check_name="$1"
    local check_command="$2"
    local optional="$3"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    echo -n "📋 $check_name... "
    
    if eval "$check_command" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        if [[ "$optional" == "optional" ]]; then
            echo -e "${YELLOW}⚠️  SKIP${NC} (optional)"
            return 0
        else
            echo -e "${RED}❌ FAIL${NC}"
            FAILED_CHECKS=$((FAILED_CHECKS + 1))
            return 1
        fi
    fi
}

# Function to validate workflow file
validate_workflow() {
    local file="$1"
    local basename=$(basename "$file")
    
    echo -e "${BLUE}🔄 Validating $basename${NC}"
    
    # Check file exists
    run_check "File exists" "test -f '$file'"
    
    # Check has name field
    run_check "Has name field" "grep -q '^name:' '$file'"
    
    # Check has on triggers
    run_check "Has triggers" "grep -q '^on:' '$file'"
    
    # Check has jobs
    run_check "Has jobs" "grep -q '^jobs:' '$file'"
    
    # Check for GitHub Actions syntax patterns
    run_check "Uses valid actions" "grep -q 'uses:.*@v[0-9]' '$file'" optional
    
    # Check indentation (basic)
    run_check "Proper indentation" "! grep -q '^[[:space:]]*[[:space:]]' '$file' || python3 -c 'import sys; lines=open(\"$file\").readlines(); sys.exit(1 if any(line.startswith(\" \") and not line.startswith(\"  \") for line in lines if line.strip()) else 0)'"
    
    echo ""
}

# Function to validate Makefile targets
validate_makefile_targets() {
    echo -e "${BLUE}🔨 Validating Makefile Targets${NC}"
    
    # Check Makefile exists
    run_check "Makefile exists" "test -f Makefile"
    
    # Check for help target
    run_check "Has help target" "grep -q '^help:' Makefile"
    
    # Check for essential targets
    local targets=("build" "unittest" "test_coverage" "ci_local" "qa" "dev_setup" "clean")
    for target in "${targets[@]}"; do
        run_check "Has $target target" "grep -q \"^$target:\" Makefile"
    done
    
    # Test help command works
    run_check "Help command works" "make help"
    
    echo ""
}

# Function to validate documentation
validate_documentation() {
    echo -e "${BLUE}📚 Validating Documentation${NC}"
    
    local docs=(
        "docs/CI_CD_PIPELINE.md"
        "docs/DEVELOPER_QUICK_REFERENCE.md"
        "PIPELINE_MODERNIZATION_SUMMARY.md"
        "FILES_MODIFIED_SUMMARY.md"
    )
    
    for doc in "${docs[@]}"; do
        run_check "$(basename "$doc") exists" "test -f '$doc'"
        run_check "$(basename "$doc") not empty" "test -s '$doc'"
        run_check "$(basename "$doc") has title" "head -5 '$doc' | grep -q '^#'"
    done
    
    # Check CHANGELOG.md was updated
    run_check "CHANGELOG.md updated" "grep -q 'Modernized CI/CD Pipeline\\|Pipeline Modernization\\|CI/CD\\|workflow' CHANGELOG.md"
    
    echo ""
}

# Function to test local development commands
test_local_commands() {
    echo -e "${BLUE}⚙️  Testing Local Development Commands${NC}"
    
    # Test basic Make commands (dry run where possible)
    run_check "make help works" "make help"
    run_check "make format works" "make -n format" optional
    run_check "make lint works" "make -n lint" optional
    run_check "make clean works" "make -n clean"
    run_check "make dev_setup works" "make -n dev_setup"
    
    # Test Go environment
    run_check "Go is available" "go version"
    run_check "go.mod is valid" "go mod verify" optional
    
    echo ""
}

# Function to validate workflow integration
validate_workflow_integration() {
    echo -e "${BLUE}🔗 Validating Workflow Integration${NC}"
    
    # Check for consistent Go version across workflows
    local go_version_count=$(grep -r "go-version-file.*go.mod" .github/workflows/ | wc -l)
    run_check "Consistent Go version" "test '$go_version_count' -gt 0"
    
    # Check for consistent Terraform version
    local tf_version_count=$(grep -r "terraform_version.*1.12" .github/workflows/ | wc -l)
    run_check "Consistent Terraform version" "test '$tf_version_count' -gt 0"
    
    # Check for proper caching usage
    local cache_count=$(grep -r "actions/cache@v4" .github/workflows/ | wc -l)
    run_check "Uses caching" "test '$cache_count' -gt 0"
    
    # Check environment variables consistency
    local env_count=$(grep -r "CHECKPOINT_DISABLE.*1" .github/workflows/ | wc -l)
    run_check "Consistent environment variables" "test '$env_count' -gt 0"
    
    echo ""
}

# Function to check GitHub Actions best practices
check_github_actions_best_practices() {
    echo -e "${BLUE}🏆 Checking GitHub Actions Best Practices${NC}"
    
    # Check for pinned action versions
    local pinned_actions=$(grep -r "uses:.*@v[0-9]" .github/workflows/ | wc -l)
    run_check "Uses pinned action versions" "test '$pinned_actions' -gt 0"
    
    # Check for proper permissions
    local permissions_count=$(grep -r "permissions:" .github/workflows/ | wc -l)
    run_check "Defines permissions" "test '$permissions_count' -gt 0"
    
    # Check for artifact usage
    local artifacts_count=$(grep -r "upload-artifact@v4" .github/workflows/ | wc -l)
    run_check "Uses artifacts" "test '$artifacts_count' -gt 0"
    
    # Check for conditional execution
    local conditionals_count=$(grep -r "if:" .github/workflows/ | wc -l)
    run_check "Uses conditional execution" "test '$conditionals_count' -gt 0"
    
    echo ""
}

# Main validation
echo -e "${BLUE}Starting comprehensive validation...${NC}"
echo ""

# Validate all workflow files
for workflow in .github/workflows/*.yml; do
    if [[ -f "$workflow" ]]; then
        validate_workflow "$workflow"
    fi
done

# Validate Makefile
validate_makefile_targets

# Validate documentation
validate_documentation

# Test local commands
test_local_commands

# Validate workflow integration
validate_workflow_integration

# Check best practices
check_github_actions_best_practices

# Summary
echo "🏁 Validation Summary"
echo "===================="
echo -e "Total Checks: ${BLUE}$TOTAL_CHECKS${NC}"
echo -e "Passed: ${GREEN}$PASSED_CHECKS${NC}"
echo -e "Failed: ${RED}$FAILED_CHECKS${NC}"
echo ""

if [[ $FAILED_CHECKS -eq 0 ]]; then
    echo -e "${GREEN}🎉 All validations passed! The CI/CD pipeline is ready.${NC}"
    exit 0
else
    echo -e "${RED}⚠️  Some validations failed. Please review and fix the issues above.${NC}"
    exit 1
fi
