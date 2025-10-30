#!/usr/bin/env bash
# Setup Terraform Environment for E2E Testing
# Installs/verifies Terraform CLI and configures local provider

set -euo pipefail

# Source required libraries
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LIB_DIR="${SCRIPT_DIR}/lib"

# shellcheck source=lib/logger.sh
source "${LIB_DIR}/logger.sh"
# shellcheck source=lib/error-handler.sh
source "${LIB_DIR}/error-handler.sh"

# Script information
readonly SCRIPT_NAME="setup-terraform"
readonly SCRIPT_VERSION="1.0.0"

# Default values
DEFAULT_TERRAFORM_VERSION="latest"
DEFAULT_INSTALL_DIR="${HOME}/.terraform-versions"
DEFAULT_PLUGIN_DIR="${HOME}/.terraform.d/plugins"

# Configuration
TERRAFORM_VERSION="${TERRAFORM_VERSION:-${DEFAULT_TERRAFORM_VERSION}}"
INSTALL_DIR="${TERRAFORM_INSTALL_DIR:-${DEFAULT_INSTALL_DIR}}"
PLUGIN_DIR="${TERRAFORM_PLUGIN_DIR:-${DEFAULT_PLUGIN_DIR}}"
PROVIDER_BINARY=""
SKIP_INSTALL="${SKIP_TERRAFORM_INSTALL:-false}"
WORKSPACE_DIR=""

#######################################
# Print usage information
#######################################
usage() {
  cat <<EOF
Usage: ${SCRIPT_NAME}.sh [OPTIONS]

Setup Terraform environment for E2E testing.

OPTIONS:
  --version VERSION        Terraform version to install/use (default: latest)
  --install-dir DIR        Directory to install Terraform versions (default: ~/.terraform-versions)
  --plugin-dir DIR         Terraform plugin directory (default: ~/.terraform.d/plugins)
  --provider-binary PATH   Path to provider binary to install
  --workspace-dir DIR      Test workspace directory for initialization
  --skip-install           Skip Terraform installation (use existing)
  --help                   Show this help message

ENVIRONMENT VARIABLES:
  TERRAFORM_VERSION        Override default Terraform version
  TERRAFORM_INSTALL_DIR    Override default install directory
  TERRAFORM_PLUGIN_DIR     Override default plugin directory
  SKIP_TERRAFORM_INSTALL   Skip Terraform installation (true/false)
  LOG_LEVEL                Set log level (DEBUG=0, INFO=1, WARN=2, ERROR=3)

EXAMPLES:
  # Setup with specific Terraform version
  ${SCRIPT_NAME}.sh --version 1.5.7 --provider-binary ./terraform-provider-azurecaf

  # Setup and initialize workspace
  ${SCRIPT_NAME}.sh --provider-binary ./terraform-provider-azurecaf --workspace-dir /tmp/test-workspace

  # Use existing Terraform installation
  ${SCRIPT_NAME}.sh --skip-install --provider-binary ./terraform-provider-azurecaf

EXIT CODES:
  0   - Success
  1   - Build failure (general error)
  2   - Dependency error (missing tools)
  125 - Configuration error (invalid arguments)

EOF
}

#######################################
# Parse command-line arguments
#######################################
parse_args() {
  while [[ $# -gt 0 ]]; do
    case $1 in
      --version)
        TERRAFORM_VERSION="$2"
        shift 2
        ;;
      --install-dir)
        INSTALL_DIR="$2"
        shift 2
        ;;
      --plugin-dir)
        PLUGIN_DIR="$2"
        shift 2
        ;;
      --provider-binary)
        PROVIDER_BINARY="$2"
        shift 2
        ;;
      --workspace-dir)
        WORKSPACE_DIR="$2"
        shift 2
        ;;
      --skip-install)
        SKIP_INSTALL=true
        shift
        ;;
      --help)
        usage
        exit 0
        ;;
      *)
        log_error "Unknown option: $1"
        usage
        exit "${EXIT_CONFIG_ERROR}"
        ;;
    esac
  done
}

#######################################
# Validate prerequisites
#######################################
check_prerequisites() {
  log_stage_start "Checking Prerequisites"
  
  # Check for required commands
  require_command "curl" "Required for downloading Terraform"
  require_command "unzip" "Required for extracting Terraform"
  require_command "jq" "Required for JSON processing"
  
  # Validate provider binary if specified
  if [[ -n "${PROVIDER_BINARY}" ]]; then
    require_file "${PROVIDER_BINARY}" "Provider binary not found"
    
    # Check if binary is executable
    if [[ ! -x "${PROVIDER_BINARY}" ]]; then
      log_error "Provider binary is not executable: ${PROVIDER_BINARY}"
      exit "${EXIT_BUILD_FAILURE}"
    fi
  fi
  
  # Validate workspace directory if specified
  if [[ -n "${WORKSPACE_DIR}" ]] && [[ ! -d "${WORKSPACE_DIR}" ]]; then
    log_error "Workspace directory does not exist: ${WORKSPACE_DIR}"
    exit "${EXIT_CONFIG_ERROR}"
  fi
  
  log_stage_end "Prerequisites Check" 0
}

#######################################
# Detect OS and architecture
#######################################
detect_platform() {
  local os=""
  local arch=""
  
  case "$(uname -s)" in
    Linux*)   os="linux" ;;
    Darwin*)  os="darwin" ;;
    MINGW*|MSYS*|CYGWIN*) os="windows" ;;
    *)
      log_error "Unsupported operating system: $(uname -s)"
      exit "${EXIT_BUILD_FAILURE}"
      ;;
  esac
  
  case "$(uname -m)" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *)
      log_error "Unsupported architecture: $(uname -m)"
      exit "${EXIT_BUILD_FAILURE}"
      ;;
  esac
  
  echo "${os}_${arch}"
}

#######################################
# Get latest Terraform version
#######################################
get_latest_terraform_version() {
  log_info "Fetching latest Terraform version"
  
  local latest_version
  latest_version=$(curl -s https://checkpoint-api.hashicorp.com/v1/check/terraform | jq -r '.current_version')
  
  if [[ -z "${latest_version}" ]] || [[ "${latest_version}" == "null" ]]; then
    log_error "Failed to fetch latest Terraform version"
    exit "${EXIT_DEPENDENCY_ERROR}"
  fi
  
  log_info "Latest Terraform version: ${latest_version}"
  echo "${latest_version}"
}

#######################################
# Install Terraform version
#######################################
install_terraform() {
  local version="$1"
  local platform
  platform=$(detect_platform)
  
  log_stage_start "Installing Terraform ${version}"
  
  # Resolve 'latest' version
  if [[ "${version}" == "latest" ]]; then
    version=$(get_latest_terraform_version)
  fi
  
  # Check if already installed
  local install_path="${INSTALL_DIR}/terraform-${version}"
  if [[ -x "${install_path}/terraform" ]]; then
    log_info "Terraform ${version} already installed at ${install_path}"
    echo "${install_path}/terraform"
    log_stage_end "Terraform Installation" 0
    return 0
  fi
  
  # Create install directory
  mkdir -p "${install_path}"
  
  # Download Terraform
  local download_url="https://releases.hashicorp.com/terraform/${version}/terraform_${version}_${platform}.zip"
  local zip_file="${install_path}/terraform.zip"
  
  log_info "Downloading from ${download_url}"
  if ! curl -sL -o "${zip_file}" "${download_url}"; then
    log_error "Failed to download Terraform ${version}"
    exit "${EXIT_DEPENDENCY_ERROR}"
  fi
  
  # Extract Terraform
  log_info "Extracting Terraform to ${install_path}"
  if ! unzip -q -o "${zip_file}" -d "${install_path}"; then
    log_error "Failed to extract Terraform archive"
    rm -f "${zip_file}"
    exit "${EXIT_BUILD_FAILURE}"
  fi
  
  # Clean up zip file
  rm -f "${zip_file}"
  
  # Verify installation
  if [[ ! -x "${install_path}/terraform" ]]; then
    log_error "Terraform binary not found after installation"
    exit "${EXIT_BUILD_FAILURE}"
  fi
  
  # Verify version
  local installed_version
  installed_version=$("${install_path}/terraform" version -json | jq -r '.terraform_version')
  log_info "Installed Terraform version: ${installed_version}"
  
  echo "${install_path}/terraform"
  log_stage_end "Terraform Installation" 0
}

#######################################
# Get Terraform binary path
#######################################
get_terraform_binary() {
  if [[ "${SKIP_INSTALL}" == "true" ]]; then
    log_info "Using existing Terraform installation"
    
    if ! command -v terraform &> /dev/null; then
      log_error "Terraform not found in PATH"
      exit "${EXIT_DEPENDENCY_ERROR}"
    fi
    
    local tf_path
    tf_path=$(command -v terraform)
    local tf_version
    tf_version=$(terraform version -json | jq -r '.terraform_version')
    
    log_info "Using Terraform ${tf_version} at ${tf_path}"
    echo "${tf_path}"
  else
    install_terraform "${TERRAFORM_VERSION}"
  fi
}

#######################################
# Install provider binary
#######################################
install_provider() {
  local terraform_bin="$1"
  
  if [[ -z "${PROVIDER_BINARY}" ]]; then
    log_info "No provider binary specified, skipping installation"
    return 0
  fi
  
  log_stage_start "Installing Provider"
  
  # Detect platform for provider directory structure
  local platform
  platform=$(detect_platform)
  
  # Provider directory structure: plugins/local/aztfmod/azurecaf/9.9.9/{platform}/
  local provider_dir="${PLUGIN_DIR}/local/aztfmod/azurecaf/9.9.9/${platform}"
  
  log_info "Installing provider to ${provider_dir}"
  mkdir -p "${provider_dir}"
  
  # Copy provider binary
  local binary_name="terraform-provider-azurecaf_v9.9.9"
  cp "${PROVIDER_BINARY}" "${provider_dir}/${binary_name}"
  chmod +x "${provider_dir}/${binary_name}"
  
  log_info "Provider installed successfully"
  log_stage_end "Provider Installation" 0
}

#######################################
# Initialize Terraform workspace
#######################################
initialize_workspace() {
  local terraform_bin="$1"
  
  if [[ -z "${WORKSPACE_DIR}" ]]; then
    log_info "No workspace directory specified, skipping initialization"
    return 0
  fi
  
  log_stage_start "Initializing Workspace"
  
  pushd "${WORKSPACE_DIR}" > /dev/null
  
  # Run terraform init
  log_info "Running terraform init in ${WORKSPACE_DIR}"
  if ! "${terraform_bin}" init -no-color > init.log 2>&1; then
    log_error "Terraform init failed"
    cat init.log
    popd > /dev/null
    exit "${EXIT_PLAN_INIT}"
  fi
  
  log_info "Workspace initialized successfully"
  popd > /dev/null
  
  log_stage_end "Workspace Initialization" 0
}

#######################################
# Generate metadata
#######################################
generate_metadata() {
  local terraform_bin="$1"
  
  log_info "Generating setup metadata"
  
  local terraform_version
  terraform_version=$("${terraform_bin}" version -json | jq -r '.terraform_version')
  
  local metadata
  metadata=$(cat <<EOF
{
  "setup_time": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "terraform_version": "${terraform_version}",
  "terraform_path": "${terraform_bin}",
  "plugin_directory": "${PLUGIN_DIR}",
  "provider_installed": $([ -n "${PROVIDER_BINARY}" ] && echo "true" || echo "false"),
  "workspace_initialized": $([ -n "${WORKSPACE_DIR}" ] && echo "true" || echo "false")
}
EOF
)
  
  echo "${metadata}"
}

#######################################
# Main execution
#######################################
main() {
  log_info "Starting ${SCRIPT_NAME} v${SCRIPT_VERSION}"
  
  # Parse arguments
  parse_args "$@"
  
  # Check prerequisites
  check_prerequisites
  
  # Get/install Terraform
  local terraform_bin
  terraform_bin=$(get_terraform_binary)
  
  # Install provider
  install_provider "${terraform_bin}"
  
  # Initialize workspace
  initialize_workspace "${terraform_bin}"
  
  # Generate metadata
  local metadata
  metadata=$(generate_metadata "${terraform_bin}")
  echo "${metadata}" | jq '.'
  
  log_info "${SCRIPT_NAME} completed successfully"
  log_info "Terraform binary: ${terraform_bin}"
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
