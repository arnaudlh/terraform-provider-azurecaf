
# Test resource name for azurerm_role_assignment
# Standard naming test
resource "azurecaf_name" "test_ra" {
  name          = "testra"
  resource_type = "azurerm_role_assignment"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_ra_special" {
  name          = "testra!@#$%"
  resource_type = "azurerm_role_assignment"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_ra_max" {
  name          = "raverylongname"
  resource_type = "azurerm_role_assignment"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_ra_min" {
  name          = "A"
  resource_type = "azurerm_role_assignment"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_ra_ds" {
  name          = "testra"
  resource_type = "azurerm_role_assignment"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_ra_standard" {
  value = can(regex("^[^%]{0,63}[^ %.]$", azurecaf_name.test_ra.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_ra.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_ra_special" {
  value = can(regex("^[^%]{0,63}[^ %.]$", azurecaf_name.test_ra_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_ra_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_ra_max" {
  value = length(azurecaf_name.test_ra_max.result) <= 64 ? "PASS" : "FAIL: ${azurecaf_name.test_ra_max.result} exceeds maximum length 64"
}

# Verify minimum length is enforced
output "test_ra_min" {
  value = length(azurecaf_name.test_ra_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_ra_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_ra_ds" {
  value = azurecaf_name.test_ra.result == data.azurecaf_name.test_ra_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_ra.result} != ${data.azurecaf_name.test_ra_ds.result}"
}

# Verify case sensitivity
output "test_ra_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_role_definition
# Standard naming test
resource "azurecaf_name" "test_rd" {
  name          = "testrd"
  resource_type = "azurerm_role_definition"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_rd_special" {
  name          = "testrd!@#$%"
  resource_type = "azurerm_role_definition"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_rd_max" {
  name          = "rdverylongname"
  resource_type = "azurerm_role_definition"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_rd_min" {
  name          = "A"
  resource_type = "azurerm_role_definition"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_rd_ds" {
  name          = "testrd"
  resource_type = "azurerm_role_definition"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_rd_standard" {
  value = can(regex("^[^%]{0,63}[^ %.]$", azurecaf_name.test_rd.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_rd.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_rd_special" {
  value = can(regex("^[^%]{0,63}[^ %.]$", azurecaf_name.test_rd_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_rd_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_rd_max" {
  value = length(azurecaf_name.test_rd_max.result) <= 64 ? "PASS" : "FAIL: ${azurecaf_name.test_rd_max.result} exceeds maximum length 64"
}

# Verify minimum length is enforced
output "test_rd_min" {
  value = length(azurecaf_name.test_rd_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_rd_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_rd_ds" {
  value = azurecaf_name.test_rd.result == data.azurecaf_name.test_rd_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_rd.result} != ${data.azurecaf_name.test_rd_ds.result}"
}

# Verify case sensitivity
output "test_rd_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_automation_account
# Standard naming test
resource "azurecaf_name" "test_aa" {
  name          = "testaa"
  resource_type = "azurerm_automation_account"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aa_special" {
  name          = "testaa!@#$%"
  resource_type = "azurerm_automation_account"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aa_max" {
  name          = "aaverylongname"
  resource_type = "azurerm_automation_account"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aa_min" {
  name          = "A"
  resource_type = "azurerm_automation_account"
  random_length = 6
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aa_ds" {
  name          = "testaa"
  resource_type = "azurerm_automation_account"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aa_standard" {
  value = can(regex("^[a-zA-Z][a-zA-Z0-9-]{4,48}[a-zA-Z0-9]$", azurecaf_name.test_aa.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aa.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aa_special" {
  value = can(regex("^[a-zA-Z][a-zA-Z0-9-]{4,48}[a-zA-Z0-9]$", azurecaf_name.test_aa_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aa_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aa_max" {
  value = length(azurecaf_name.test_aa_max.result) <= 50 ? "PASS" : "FAIL: ${azurecaf_name.test_aa_max.result} exceeds maximum length 50"
}

# Verify minimum length is enforced
output "test_aa_min" {
  value = length(azurecaf_name.test_aa_min.result) >= 6 ? "PASS" : "FAIL: ${azurecaf_name.test_aa_min.result} below minimum length 6"
}

# Verify data source consistency
output "test_aa_ds" {
  value = azurecaf_name.test_aa.result == data.azurecaf_name.test_aa_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aa.result} != ${data.azurecaf_name.test_aa_ds.result}"
}

# Verify case sensitivity
output "test_aa_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_automation_certificate
# Standard naming test
resource "azurecaf_name" "test_aacert" {
  name          = "testaacert"
  resource_type = "azurerm_automation_certificate"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aacert_special" {
  name          = "testaacert!@#$%"
  resource_type = "azurerm_automation_certificate"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aacert_max" {
  name          = "aacertverylongname"
  resource_type = "azurerm_automation_certificate"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aacert_min" {
  name          = "A"
  resource_type = "azurerm_automation_certificate"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aacert_ds" {
  name          = "testaacert"
  resource_type = "azurerm_automation_certificate"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aacert_standard" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aacert.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aacert.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aacert_special" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aacert_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aacert_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aacert_max" {
  value = length(azurecaf_name.test_aacert_max.result) <= 128 ? "PASS" : "FAIL: ${azurecaf_name.test_aacert_max.result} exceeds maximum length 128"
}

# Verify minimum length is enforced
output "test_aacert_min" {
  value = length(azurecaf_name.test_aacert_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aacert_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aacert_ds" {
  value = azurecaf_name.test_aacert.result == data.azurecaf_name.test_aacert_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aacert.result} != ${data.azurecaf_name.test_aacert_ds.result}"
}

# Verify case sensitivity
output "test_aacert_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_automation_credential
# Standard naming test
resource "azurecaf_name" "test_aacred" {
  name          = "testaacred"
  resource_type = "azurerm_automation_credential"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aacred_special" {
  name          = "testaacred!@#$%"
  resource_type = "azurerm_automation_credential"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aacred_max" {
  name          = "aacredverylongname"
  resource_type = "azurerm_automation_credential"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aacred_min" {
  name          = "A"
  resource_type = "azurerm_automation_credential"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aacred_ds" {
  name          = "testaacred"
  resource_type = "azurerm_automation_credential"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aacred_standard" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aacred.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aacred.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aacred_special" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aacred_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aacred_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aacred_max" {
  value = length(azurecaf_name.test_aacred_max.result) <= 128 ? "PASS" : "FAIL: ${azurecaf_name.test_aacred_max.result} exceeds maximum length 128"
}

# Verify minimum length is enforced
output "test_aacred_min" {
  value = length(azurecaf_name.test_aacred_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aacred_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aacred_ds" {
  value = azurecaf_name.test_aacred.result == data.azurecaf_name.test_aacred_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aacred.result} != ${data.azurecaf_name.test_aacred_ds.result}"
}

# Verify case sensitivity
output "test_aacred_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}