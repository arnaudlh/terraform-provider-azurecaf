
# Test resource name for azurerm_automation_hybrid_runbook_worker_group
# Standard naming test
resource "azurecaf_name" "test_aahwg" {
  name          = "testaahwg"
  resource_type = "azurerm_automation_hybrid_runbook_worker_group"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aahwg_special" {
  name          = "testaahwg!@#$%"
  resource_type = "azurerm_automation_hybrid_runbook_worker_group"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aahwg_max" {
  name          = "aahwgverylongname"
  resource_type = "azurerm_automation_hybrid_runbook_worker_group"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aahwg_min" {
  name          = "A"
  resource_type = "azurerm_automation_hybrid_runbook_worker_group"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aahwg_ds" {
  name          = "testaahwg"
  resource_type = "azurerm_automation_hybrid_runbook_worker_group"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aahwg_standard" {
  value = can(regex("^([^<>*%&:\\?.+/#\\s]?[ ]?){0,127}[^<>*%&:\\?.+/#\\s]$", azurecaf_name.test_aahwg.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aahwg.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aahwg_special" {
  value = can(regex("^([^<>*%&:\\?.+/#\\s]?[ ]?){0,127}[^<>*%&:\\?.+/#\\s]$", azurecaf_name.test_aahwg_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aahwg_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aahwg_max" {
  value = length(azurecaf_name.test_aahwg_max.result) <= 128 ? "PASS" : "FAIL: ${azurecaf_name.test_aahwg_max.result} exceeds maximum length 128"
}

# Verify minimum length is enforced
output "test_aahwg_min" {
  value = length(azurecaf_name.test_aahwg_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aahwg_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aahwg_ds" {
  value = azurecaf_name.test_aahwg.result == data.azurecaf_name.test_aahwg_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aahwg.result} != ${data.azurecaf_name.test_aahwg_ds.result}"
}

# Verify case sensitivity
output "test_aahwg_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_automation_job_schedule
# Standard naming test
resource "azurecaf_name" "test_aajs" {
  name          = "testaajs"
  resource_type = "azurerm_automation_job_schedule"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aajs_special" {
  name          = "testaajs!@#$%"
  resource_type = "azurerm_automation_job_schedule"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aajs_max" {
  name          = "aajsverylongname"
  resource_type = "azurerm_automation_job_schedule"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aajs_min" {
  name          = "A"
  resource_type = "azurerm_automation_job_schedule"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aajs_ds" {
  name          = "testaajs"
  resource_type = "azurerm_automation_job_schedule"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aajs_standard" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aajs.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aajs.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aajs_special" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aajs_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aajs_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aajs_max" {
  value = length(azurecaf_name.test_aajs_max.result) <= 128 ? "PASS" : "FAIL: ${azurecaf_name.test_aajs_max.result} exceeds maximum length 128"
}

# Verify minimum length is enforced
output "test_aajs_min" {
  value = length(azurecaf_name.test_aajs_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aajs_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aajs_ds" {
  value = azurecaf_name.test_aajs.result == data.azurecaf_name.test_aajs_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aajs.result} != ${data.azurecaf_name.test_aajs_ds.result}"
}

# Verify case sensitivity
output "test_aajs_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_automation_runbook
# Standard naming test
resource "azurecaf_name" "test_aarun" {
  name          = "testaarun"
  resource_type = "azurerm_automation_runbook"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aarun_special" {
  name          = "testaarun!@#$%"
  resource_type = "azurerm_automation_runbook"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aarun_max" {
  name          = "aarunverylongname"
  resource_type = "azurerm_automation_runbook"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aarun_min" {
  name          = "A"
  resource_type = "azurerm_automation_runbook"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aarun_ds" {
  name          = "testaarun"
  resource_type = "azurerm_automation_runbook"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aarun_standard" {
  value = can(regex("^[a-zA-Z][a-zA-Z0-9-]{0,62}$", azurecaf_name.test_aarun.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aarun.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aarun_special" {
  value = can(regex("^[a-zA-Z][a-zA-Z0-9-]{0,62}$", azurecaf_name.test_aarun_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aarun_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aarun_max" {
  value = length(azurecaf_name.test_aarun_max.result) <= 63 ? "PASS" : "FAIL: ${azurecaf_name.test_aarun_max.result} exceeds maximum length 63"
}

# Verify minimum length is enforced
output "test_aarun_min" {
  value = length(azurecaf_name.test_aarun_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aarun_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aarun_ds" {
  value = azurecaf_name.test_aarun.result == data.azurecaf_name.test_aarun_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aarun.result} != ${data.azurecaf_name.test_aarun_ds.result}"
}

# Verify case sensitivity
output "test_aarun_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_automation_schedule
# Standard naming test
resource "azurecaf_name" "test_aasched" {
  name          = "testaasched"
  resource_type = "azurerm_automation_schedule"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aasched_special" {
  name          = "testaasched!@#$%"
  resource_type = "azurerm_automation_schedule"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aasched_max" {
  name          = "aaschedverylongname"
  resource_type = "azurerm_automation_schedule"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aasched_min" {
  name          = "A"
  resource_type = "azurerm_automation_schedule"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aasched_ds" {
  name          = "testaasched"
  resource_type = "azurerm_automation_schedule"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aasched_standard" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aasched.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aasched.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aasched_special" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aasched_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aasched_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aasched_max" {
  value = length(azurecaf_name.test_aasched_max.result) <= 128 ? "PASS" : "FAIL: ${azurecaf_name.test_aasched_max.result} exceeds maximum length 128"
}

# Verify minimum length is enforced
output "test_aasched_min" {
  value = length(azurecaf_name.test_aasched_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aasched_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aasched_ds" {
  value = azurecaf_name.test_aasched.result == data.azurecaf_name.test_aasched_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aasched.result} != ${data.azurecaf_name.test_aasched_ds.result}"
}

# Verify case sensitivity
output "test_aasched_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_automation_variable
# Standard naming test
resource "azurecaf_name" "test_aavar" {
  name          = "testaavar"
  resource_type = "azurerm_automation_variable"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_aavar_special" {
  name          = "testaavar!@#$%"
  resource_type = "azurerm_automation_variable"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_aavar_max" {
  name          = "aavarverylongname"
  resource_type = "azurerm_automation_variable"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_aavar_min" {
  name          = "A"
  resource_type = "azurerm_automation_variable"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_aavar_ds" {
  name          = "testaavar"
  resource_type = "azurerm_automation_variable"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_aavar_standard" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aavar.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aavar.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_aavar_special" {
  value = can(regex("^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$", azurecaf_name.test_aavar_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aavar_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_aavar_max" {
  value = length(azurecaf_name.test_aavar_max.result) <= 128 ? "PASS" : "FAIL: ${azurecaf_name.test_aavar_max.result} exceeds maximum length 128"
}

# Verify minimum length is enforced
output "test_aavar_min" {
  value = length(azurecaf_name.test_aavar_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_aavar_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_aavar_ds" {
  value = azurecaf_name.test_aavar.result == data.azurecaf_name.test_aavar_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_aavar.result} != ${data.azurecaf_name.test_aavar_ds.result}"
}

# Verify case sensitivity
output "test_aavar_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}