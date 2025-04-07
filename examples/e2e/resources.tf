resource "azurecaf_name" "test_rg" {
  name          = "testrg"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_st" {
  name          = "testst"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_kv" {
  name          = "testkv"
  resource_type = "azurerm_key_vault"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_vm" {
  name          = "testvm"
  resource_type = "azurerm_virtual_machine"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

resource "azurecaf_name" "test_aks" {
  name          = "testaks"
  resource_type = "azurerm_kubernetes_cluster"
  prefixes      = ["prefix"]
  suffixes      = ["suffix"]
  random_length = 5
  clean_input   = true
}

output "test_rg_result" {
  value = can(regex("^prefix-rg-testrg-[a-z0-9]{5}-suffix$", azurecaf_name.test_rg.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_rg.result} does not match expected pattern"
}

output "test_st_result" {
  value = can(regex("^sttestst[a-z0-9]{5}suffix$", azurecaf_name.test_st.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_st.result} does not match expected pattern"
}

output "test_kv_result" {
  value = can(regex("^kv-testkv-[a-z0-9]{5}-suffix$", azurecaf_name.test_kv.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_kv.result} does not match expected pattern"
}

output "test_vm_result" {
  value = can(regex("^vm-testvm-[a-z0-9]{5}$", azurecaf_name.test_vm.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_vm.result} does not match expected pattern"
}

output "test_aks_result" {
  value = can(regex("^prefix-aks-testaks-[a-z0-9]{5}-suffix$", azurecaf_name.test_aks.result)) ? "PASS" : "FAIL: ${azurecaf_name.test_aks.result} does not match expected pattern"
}
