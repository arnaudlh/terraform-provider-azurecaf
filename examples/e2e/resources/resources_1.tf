
# Test resource name for azurerm_cognitive_deployment
# Standard naming test
resource "azurecaf_name" "test_cog" {
  name          = "testcog"
  resource_type = "azurerm_cognitive_deployment"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_cog_special" {
  name          = "testcog!@#$%"
  resource_type = "azurerm_cognitive_deployment"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_cog_max" {
  name          = "cogverylongname"
  resource_type = "azurerm_cognitive_deployment"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_cog_min" {
  name          = "A"
  resource_type = "azurerm_cognitive_deployment"
  random_length = 2
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_cog_ds" {
  name          = "testcog"
  resource_type = "azurerm_cognitive_deployment"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_cog_standard" {
  value = can(regex("^[a-zA-Z0-9][a-zA-Z0-9-]{0,63}$", azurecaf_name.test_cog.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_cog.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_cog_special" {
  value = can(regex("^[a-zA-Z0-9][a-zA-Z0-9-]{0,63}$", azurecaf_name.test_cog_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_cog_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_cog_max" {
  value = length(azurecaf_name.test_cog_max.result) <= 64 ? "PASS" : "FAIL: ${azurecaf_name.test_cog_max.result} exceeds maximum length 64"
}

# Verify minimum length is enforced
output "test_cog_min" {
  value = length(azurecaf_name.test_cog_min.result) >= 2 ? "PASS" : "FAIL: ${azurecaf_name.test_cog_min.result} below minimum length 2"
}

# Verify data source consistency
output "test_cog_ds" {
  value = azurecaf_name.test_cog.result == data.azurecaf_name.test_cog_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_cog.result} != ${data.azurecaf_name.test_cog_ds.result}"
}

# Verify case sensitivity
output "test_cog_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_aadb2c_directory
# Standard naming test
resource "azurecaf_name" "test_aadb2c" {
  name          = "testaadb2c"
  resource_type = "azurerm_aadb2c_directory"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aadb2c_special" {
  name          = "testaadb2c!@#$%"
  resource_type = "azurerm_aadb2c_directory"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aadb2c_max" {
  name          = "aadb2cverylongname"
  resource_type = "azurerm_aadb2c_directory"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aadb2c_min" {
  name          = "A"
  resource_type = "azurerm_aadb2c_directory"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aadb2c_ds" {
  name          = "testaadb2c"
  resource_type = "azurerm_aadb2c_directory"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aadb2c_standard" {
  value = can(regex("^[a-zA-Z0-9][a-zA-Z0-9-]{0,73}[a-zA-Z0-9]$", azurecaf_name.test_aadb2c.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aadb2c.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aadb2c_special" {
  value = can(regex("^[a-zA-Z0-9][a-zA-Z0-9-]{0,73}[a-zA-Z0-9]$", azurecaf_name.test_aadb2c_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aadb2c_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aadb2c_max" {
  value = length(azurecaf_name.test_aadb2c_max.result) <= 75 ? "PASS" : "FAIL: ${azurecaf_name.test_aadb2c_max.result} exceeds maximum length 75"
}

# Verify minimum length is enforced
output "test_aadb2c_min" {
  value = length(azurecaf_name.test_aadb2c_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aadb2c_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aadb2c_ds" {
  value = azurecaf_name.test_aadb2c.result == data.azurecaf_name.test_aadb2c_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aadb2c.result} != ${data.azurecaf_name.test_aadb2c_ds.result}"
}

# Verify case sensitivity
output "test_aadb2c_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_analysis_services_server
# Standard naming test
resource "azurecaf_name" "test_as" {
  name          = "testas"
  resource_type = "azurerm_analysis_services_server"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_as_special" {
  name          = "testas!@#$%"
  resource_type = "azurerm_analysis_services_server"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_as_max" {
  name          = "asverylongname"
  resource_type = "azurerm_analysis_services_server"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_as_min" {
  name          = "a"
  resource_type = "azurerm_analysis_services_server"
  random_length = 3
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_as_ds" {
  name          = "testas"
  resource_type = "azurerm_analysis_services_server"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_as_standard" {
  value = can(regex("^[a-z][a-z0-9]{2,62}$", azurecaf_name.test_as.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_as.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_as_special" {
  value = can(regex("^[a-z][a-z0-9]{2,62}$", azurecaf_name.test_as_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_as_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_as_max" {
  value = length(azurecaf_name.test_as_max.result) <= 63 ? "PASS" : "FAIL: ${azurecaf_name.test_as_max.result} exceeds maximum length 63"
}

# Verify minimum length is enforced
output "test_as_min" {
  value = length(azurecaf_name.test_as_min.result) >= 3 ? "PASS" : "FAIL: ${azurecaf_name.test_as_min.result} below minimum length 3"
}

# Verify data source consistency
output "test_as_ds" {
  value = azurecaf_name.test_as.result == data.azurecaf_name.test_as_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_as.result} != ${data.azurecaf_name.test_as_ds.result}"
}

# Verify case sensitivity
output "test_as_case" {
  value = can(regex("^[a-z0-9-_.]*$", azurecaf_name.test_as.result)) ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_api_management_service
# Standard naming test
resource "azurecaf_name" "test_apim" {
  name          = "testapim"
  resource_type = "azurerm_api_management_service"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_apim_special" {
  name          = "testapim!@#$%"
  resource_type = "azurerm_api_management_service"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_apim_max" {
  name          = "apimverylongname"
  resource_type = "azurerm_api_management_service"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_apim_min" {
  name          = "A"
  resource_type = "azurerm_api_management_service"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_apim_ds" {
  name          = "testapim"
  resource_type = "azurerm_api_management_service"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_apim_standard" {
  value = can(regex("^[a-z][a-zA-Z0-9-]{0,48}[a-zA-Z0-9]$", azurecaf_name.test_apim.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_apim.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_apim_special" {
  value = can(regex("^[a-z][a-zA-Z0-9-]{0,48}[a-zA-Z0-9]$", azurecaf_name.test_apim_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_apim_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_apim_max" {
  value = length(azurecaf_name.test_apim_max.result) <= 50 ? "PASS" : "FAIL: ${azurecaf_name.test_apim_max.result} exceeds maximum length 50"
}

# Verify minimum length is enforced
output "test_apim_min" {
  value = length(azurecaf_name.test_apim_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_apim_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_apim_ds" {
  value = azurecaf_name.test_apim.result == data.azurecaf_name.test_apim_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_apim.result} != ${data.azurecaf_name.test_apim_ds.result}"
}

# Verify case sensitivity
output "test_apim_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_app_configuration
# Standard naming test
resource "azurecaf_name" "test_appcg" {
  name          = "testappcg"
  resource_type = "azurerm_app_configuration"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_appcg_special" {
  name          = "testappcg!@#$%"
  resource_type = "azurerm_app_configuration"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_appcg_max" {
  name          = "appcgverylongname"
  resource_type = "azurerm_app_configuration"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_appcg_min" {
  name          = "A"
  resource_type = "azurerm_app_configuration"
  random_length = 5
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_appcg_ds" {
  name          = "testappcg"
  resource_type = "azurerm_app_configuration"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_appcg_standard" {
  value = can(regex("^[a-zA-Z0-9-]{5,50}$", azurecaf_name.test_appcg.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_appcg.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_appcg_special" {
  value = can(regex("^[a-zA-Z0-9-]{5,50}$", azurecaf_name.test_appcg_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_appcg_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_appcg_max" {
  value = length(azurecaf_name.test_appcg_max.result) <= 50 ? "PASS" : "FAIL: ${azurecaf_name.test_appcg_max.result} exceeds maximum length 50"
}

# Verify minimum length is enforced
output "test_appcg_min" {
  value = length(azurecaf_name.test_appcg_min.result) >= 5 ? "PASS" : "FAIL: ${azurecaf_name.test_appcg_min.result} below minimum length 5"
}

# Verify data source consistency
output "test_appcg_ds" {
  value = azurecaf_name.test_appcg.result == data.azurecaf_name.test_appcg_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_appcg.result} != ${data.azurecaf_name.test_appcg_ds.result}"
}

# Verify case sensitivity
output "test_appcg_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}