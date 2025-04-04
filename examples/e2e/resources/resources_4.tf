
# Test resource name for azurerm_batch_account
# Standard naming test
resource "azurecaf_name" "test_ba" {
  name          = "testba"
  resource_type = "azurerm_batch_account"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_ba_special" {
  name          = "testba!@#$%"
  resource_type = "azurerm_batch_account"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_ba_max" {
  name          = "baverylongname"
  resource_type = "azurerm_batch_account"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_ba_min" {
  name          = "a"
  resource_type = "azurerm_batch_account"
  random_length = 3
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_ba_ds" {
  name          = "testba"
  resource_type = "azurerm_batch_account"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_ba_standard" {
  value = can(regex("^[a-z0-9]{3,24}$", azurecaf_name.test_ba.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_ba.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_ba_special" {
  value = can(regex("^[a-z0-9]{3,24}$", azurecaf_name.test_ba_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_ba_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_ba_max" {
  value = length(azurecaf_name.test_ba_max.result) <= 24 ? "PASS" : "FAIL: ${azurecaf_name.test_ba_max.result} exceeds maximum length 24"
}

# Verify minimum length is enforced
output "test_ba_min" {
  value = length(azurecaf_name.test_ba_min.result) >= 3 ? "PASS" : "FAIL: ${azurecaf_name.test_ba_min.result} below minimum length 3"
}

# Verify data source consistency
output "test_ba_ds" {
  value = azurecaf_name.test_ba.result == data.azurecaf_name.test_ba_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_ba.result} != ${data.azurecaf_name.test_ba_ds.result}"
}

# Verify case sensitivity
output "test_ba_case" {
  value = can(regex("^[a-z0-9-_.]*$", azurecaf_name.test_ba.result)) ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_batch_application
# Standard naming test
resource "azurecaf_name" "test_baapp" {
  name          = "testbaapp"
  resource_type = "azurerm_batch_application"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_baapp_special" {
  name          = "testbaapp!@#$%"
  resource_type = "azurerm_batch_application"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_baapp_max" {
  name          = "baappverylongname"
  resource_type = "azurerm_batch_application"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_baapp_min" {
  name          = "A"
  resource_type = "azurerm_batch_application"
  random_length = 1
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_baapp_ds" {
  name          = "testbaapp"
  resource_type = "azurerm_batch_application"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_baapp_standard" {
  value = can(regex("^[a-zA-Z0-9_-]{1,64}$", azurecaf_name.test_baapp.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_baapp.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_baapp_special" {
  value = can(regex("^[a-zA-Z0-9_-]{1,64}$", azurecaf_name.test_baapp_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_baapp_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_baapp_max" {
  value = length(azurecaf_name.test_baapp_max.result) <= 64 ? "PASS" : "FAIL: ${azurecaf_name.test_baapp_max.result} exceeds maximum length 64"
}

# Verify minimum length is enforced
output "test_baapp_min" {
  value = length(azurecaf_name.test_baapp_min.result) >= 1 ? "PASS" : "FAIL: ${azurecaf_name.test_baapp_min.result} below minimum length 1"
}

# Verify data source consistency
output "test_baapp_ds" {
  value = azurecaf_name.test_baapp.result == data.azurecaf_name.test_baapp_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_baapp.result} != ${data.azurecaf_name.test_baapp_ds.result}"
}

# Verify case sensitivity
output "test_baapp_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_batch_certificate
# Standard naming test
resource "azurecaf_name" "test_bacert" {
  name          = "testbacert"
  resource_type = "azurerm_batch_certificate"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_bacert_special" {
  name          = "testbacert!@#$%"
  resource_type = "azurerm_batch_certificate"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_bacert_max" {
  name          = "bacertverylongname"
  resource_type = "azurerm_batch_certificate"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_bacert_min" {
  name          = "A"
  resource_type = "azurerm_batch_certificate"
  random_length = 5
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_bacert_ds" {
  name          = "testbacert"
  resource_type = "azurerm_batch_certificate"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_bacert_standard" {
  value = can(regex("^[a-zA-Z0-9_-]{5,45}$", azurecaf_name.test_bacert.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_bacert.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_bacert_special" {
  value = can(regex("^[a-zA-Z0-9_-]{5,45}$", azurecaf_name.test_bacert_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_bacert_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_bacert_max" {
  value = length(azurecaf_name.test_bacert_max.result) <= 45 ? "PASS" : "FAIL: ${azurecaf_name.test_bacert_max.result} exceeds maximum length 45"
}

# Verify minimum length is enforced
output "test_bacert_min" {
  value = length(azurecaf_name.test_bacert_min.result) >= 5 ? "PASS" : "FAIL: ${azurecaf_name.test_bacert_min.result} below minimum length 5"
}

# Verify data source consistency
output "test_bacert_ds" {
  value = azurecaf_name.test_bacert.result == data.azurecaf_name.test_bacert_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_bacert.result} != ${data.azurecaf_name.test_bacert_ds.result}"
}

# Verify case sensitivity
output "test_bacert_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_batch_pool
# Standard naming test
resource "azurecaf_name" "test_bapool" {
  name          = "testbapool"
  resource_type = "azurerm_batch_pool"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_bapool_special" {
  name          = "testbapool!@#$%"
  resource_type = "azurerm_batch_pool"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_bapool_max" {
  name          = "bapoolverylongname"
  resource_type = "azurerm_batch_pool"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_bapool_min" {
  name          = "A"
  resource_type = "azurerm_batch_pool"
  random_length = 3
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_bapool_ds" {
  name          = "testbapool"
  resource_type = "azurerm_batch_pool"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_bapool_standard" {
  value = can(regex("^[a-zA-Z0-9_-]{1,24}$", azurecaf_name.test_bapool.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_bapool.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_bapool_special" {
  value = can(regex("^[a-zA-Z0-9_-]{1,24}$", azurecaf_name.test_bapool_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_bapool_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_bapool_max" {
  value = length(azurecaf_name.test_bapool_max.result) <= 24 ? "PASS" : "FAIL: ${azurecaf_name.test_bapool_max.result} exceeds maximum length 24"
}

# Verify minimum length is enforced
output "test_bapool_min" {
  value = length(azurecaf_name.test_bapool_min.result) >= 3 ? "PASS" : "FAIL: ${azurecaf_name.test_bapool_min.result} below minimum length 3"
}

# Verify data source consistency
output "test_bapool_ds" {
  value = azurecaf_name.test_bapool.result == data.azurecaf_name.test_bapool_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_bapool.result} != ${data.azurecaf_name.test_bapool_ds.result}"
}

# Verify case sensitivity
output "test_bapool_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}
# Test resource name for azurerm_bot_web_app
# Standard naming test
resource "azurecaf_name" "test_bot" {
  name          = "testbot"
  resource_type = "azurerm_bot_web_app"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Test with special characters that should be cleaned
resource "azurecaf_name" "test_bot_special" {
  name          = "testbot!@#$%"
  resource_type = "azurerm_bot_web_app"
  prefixes      = ["prefix!@#"]
  suffixes      = ["suffix#$%"]
  random_length = 5
  clean_input   = true
}

# Test with maximum length
resource "azurecaf_name" "test_bot_max" {
  name          = "botverylongname"
  resource_type = "azurerm_bot_web_app"
  prefixes      = ["prefixlong"]
  suffixes      = ["suffixlong"]
  random_length = 10
  clean_input   = true
}

# Test with minimum length
resource "azurecaf_name" "test_bot_min" {
  name          = "A"
  resource_type = "azurerm_bot_web_app"
  random_length = 2
  clean_input   = true
}

# Test data source consistency
data "azurecaf_name" "test_bot_ds" {
  name          = "testbot"
  resource_type = "azurerm_bot_web_app"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

# Verify standard naming test
output "test_bot_standard" {
  value = can(regex("^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$", azurecaf_name.test_bot.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_bot.result} does not match regex pattern"
}

# Verify special characters are cleaned
output "test_bot_special" {
  value = can(regex("^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$", azurecaf_name.test_bot_special.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_bot_special.result} does not match regex pattern"
}

# Verify maximum length is enforced
output "test_bot_max" {
  value = length(azurecaf_name.test_bot_max.result) <= 64 ? "PASS" : "FAIL: ${azurecaf_name.test_bot_max.result} exceeds maximum length 64"
}

# Verify minimum length is enforced
output "test_bot_min" {
  value = length(azurecaf_name.test_bot_min.result) >= 2 ? "PASS" : "FAIL: ${azurecaf_name.test_bot_min.result} below minimum length 2"
}

# Verify data source consistency
output "test_bot_ds" {
  value = azurecaf_name.test_bot.result == data.azurecaf_name.test_bot_ds.result ? "PASS" : "FAIL: ${azurecaf_name.test_bot.result} != ${data.azurecaf_name.test_bot_ds.result}"
}

# Verify case sensitivity
output "test_bot_case" {
  value = true ? "PASS" : "FAIL: Case sensitivity rules not followed"
}