# E2E Validation Pipeline - Plan & Apply Results

This document explains how Terraform plan and apply results are displayed in the E2E validation pipeline.

---

## 📋 Plan Validation Results (Default)

### What Gets Shown

By default, the pipeline runs in **plan-only mode** which shows what would be created without actually creating resources.

#### 1. **GitHub Actions Step Summary**

```markdown
### 📋 Plan Validation Results

✅ Plan validation completed
- **Total Tests**: 80
- **Passed**: 80
- **Failed**: 0
- **Pass Rate**: 100%

### 📝 Sample Generated Names

<details>
<summary><b>cafclassic</b> - Click to expand</summary>

{
  "vm": "vm-dev-testapp-001",
  "storage_account": "stdevtestapp001",
  "key_vault": "kv-dev-testapp-001",
  "virtual_network": "vnet-dev-testapp-001",
  ...
}
</details>
```

#### 2. **Plan Output Logs** (in workflow steps)

Each convention's plan output is displayed in collapsible groups:

```
::group::Plan Output - cafclassic
Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # data.azurecaf_name.storage_account will be read during apply
  [...]
  
::endgroup::
```

#### 3. **Generated Names Summary**

```
All generated names:
azurerm_storage_account (cafclassic): pass
azurerm_key_vault (cafclassic): pass
azurerm_virtual_machine (cafclassic): pass
...
```

#### 4. **Artifacts Uploaded**

Available for download from the workflow run:
- `validation-report.json` - Complete test results
- `plan-cafclassic.json` - Full Terraform plan in JSON format
- `plan-cafrandom.json` - CAF random convention plan
- `plan-random.json` - Pure random convention plan
- `plan-passthrough.json` - Passthrough convention plan
- `*.log` files - Raw Terraform output

---

## 🚀 Apply Validation Results (Opt-in)

### When Apply Runs

Apply validation is **disabled by default** and only runs when:
- PR has the `test-apply` label
- All plan validations pass
- Azure credentials are available (in real GitHub Actions)

### What Gets Shown

#### 1. **GitHub Actions Step Summary**

```markdown
### 🚀 Apply Validation Results

- ✅ **cafclassic**: Apply successful
- ✅ **cafrandom**: Apply successful  
- ✅ **random**: Apply successful
- ✅ **passthrough**: Apply successful

**Summary**: 4 successful, 0 failed

> Terraform apply executed successfully and resources were created/destroyed.

💡 **Note**: Apply validation runs only when PR has `test-apply` label.
```

#### 2. **Apply Output Logs**

Each convention's apply output is displayed:

```
::group::Apply Output - cafclassic

Terraform will perform the following actions:

  # azurecaf_name.storage_account will be created
  + resource "azurecaf_name" "storage_account" {
      + id            = (known after apply)
      + name          = "testapp"
      + result        = (known after apply)
      + resource_type = "azurerm_storage_account"
      ...
    }

Apply complete! Resources: 20 added, 0 changed, 0 destroyed.

Outputs:
all_names = {
  "storage_account" = "stdevtestapp001"
  "key_vault" = "kv-dev-testapp-001"
  ...
}

::endgroup::
```

#### 3. **Destroy Output Logs**

After successful apply, resources are immediately destroyed:

```
::group::Destroy Output - cafclassic

Terraform will perform the following actions:

  # azurecaf_name.storage_account will be destroyed
  - resource "azurecaf_name" "storage_account" {
      ...
    }

Destroy complete! Resources: 20 destroyed.

::endgroup::
```

#### 4. **Additional Artifacts**

When apply runs, additional artifacts are uploaded:
- `apply-*.log` - Apply output for each convention
- `destroy-*.log` - Destroy output for each convention

---

## 🔍 Viewing Results in Different Contexts

### 1. **In Pull Request Comments**

The bot automatically comments on PRs with:

```markdown
## E2E Validation - Quick Mode

### ✅ Build Successful

- **Duration**: 42s
- **Binary Hash**: `d9b86510...`

### ✅ Plan Validation Successful

| Metric | Value |
|--------|-------|
| Total Tests | 80 |
| Passed | 80 |
| Failed | 0 |
| Pass Rate | 100% |

> All Terraform configurations generated valid resource names.

### 🚀 Apply Validation Results

- **Applied**: 4 conventions
- **Failed**: 0 conventions

> Terraform apply executed successfully and resources were created/destroyed.

💡 **Note**: Apply validation runs only when PR has `test-apply` label.
```

### 2. **In GitHub Actions Logs**

Navigate to: Actions → E2E Validation - Pull Request → Quick E2E Validation

Expand steps to see:
- ✅ Show Plan Details
- ✅ Show Apply Details (if apply ran)

### 3. **In Downloaded Artifacts**

From the workflow run page, download `plan-validation-results.zip`:

```
plan-validation-results/
├── validation-report.json      # Test results summary
├── plan-cafclassic.json       # Full plan JSON
├── plan-cafclassic.log        # Plan output
├── plan-cafrandom.json
├── plan-cafrandom.log
├── apply-cafclassic.log       # Apply output (if apply ran)
├── destroy-cafclassic.log     # Destroy output (if apply ran)
└── init.log                   # Terraform init output
```

---

## 📊 Example Validation Report JSON

```json
{
  "validation_id": "validation-quick-20251030-145800",
  "status": "success",
  "stages": {
    "build": {
      "status": "success",
      "duration": 42,
      "started_at": "2025-10-30T06:58:00Z",
      "ended_at": "2025-10-30T06:58:42Z"
    },
    "plan": {
      "status": "success",
      "duration": 18,
      "started_at": "2025-10-30T06:58:42Z",
      "ended_at": "2025-10-30T06:59:00Z"
    }
  },
  "tests": [
    {
      "resource_type": "azurerm_storage_account",
      "convention": "cafclassic",
      "status": "pass",
      "error": null,
      "generated_name": "stdevtestapp001"
    },
    {
      "resource_type": "azurerm_key_vault",
      "convention": "cafclassic",
      "status": "pass",
      "error": null,
      "generated_name": "kv-dev-testapp-001"
    }
  ],
  "summary": {
    "total_tests": 80,
    "passed": 80,
    "failed": 0,
    "pass_rate": 100
  }
}
```

---

## 🎯 Typical Workflow

### Plan-Only (Default - Every PR)

1. **Trigger**: Push to PR
2. **Runs**:
   - Build provider ✅
   - Terraform plan for 4 conventions ✅
   - Validate generated names ✅
3. **Duration**: ~2-3 minutes
4. **Shows**:
   - What names would be generated
   - Whether names meet constraints
   - No actual resources created

### Plan + Apply (Opt-in - Label Required)

1. **Trigger**: Add `test-apply` label to PR
2. **Runs**:
   - Build provider ✅
   - Terraform plan for 4 conventions ✅
   - Validate generated names ✅
   - **Terraform apply** (creates resources) ✅
   - **Terraform destroy** (removes resources) ✅
3. **Duration**: ~5-8 minutes
4. **Shows**:
   - Same as plan-only, PLUS:
   - Actual resource creation output
   - Resource destruction output
   - Proof that names work in practice

---

## 💡 Key Differences

| Aspect | Plan-Only | Plan + Apply |
|--------|-----------|--------------|
| **Trigger** | Every PR push | PR with `test-apply` label |
| **Duration** | 2-3 min | 5-8 min |
| **Creates Resources** | ❌ No | ✅ Yes (then destroys) |
| **Shows** | Planned names | Actual creation |
| **Use Case** | Fast feedback | Full validation |
| **Default** | ✅ Yes | ❌ No |

---

## 🚀 Enabling Apply in Your PR

To see apply results in your PR:

1. Add label to PR:
   ```bash
   gh pr edit <PR_NUMBER> --add-label "test-apply"
   ```

2. Or via GitHub UI:
   - Go to PR page
   - Click "Labels" on the right sidebar
   - Select "test-apply"

3. Workflow will automatically re-run with apply enabled

---

## ⚠️ Important Notes

### Plan-Only Mode (Safe)
- ✅ No credentials needed
- ✅ No resources created
- ✅ No cost incurred
- ✅ Fast feedback
- ✅ Safe for all PRs

### Apply Mode (Use with Caution)
- ⚠️ Creates real resources (naming provider doesn't need Azure credentials)
- ⚠️ Requires proper cleanup
- ⚠️ Slower execution
- ⚠️ Use only for thorough validation
- ✅ Proves names work end-to-end

---

## 📝 Summary

The E2E validation pipeline shows:

1. **Always** (Plan-Only):
   - Build status and metadata
   - Generated names for all conventions
   - Name constraint validation results
   - Downloadable plan JSON files

2. **Optionally** (with `test-apply` label):
   - Resource creation logs
   - Resource destruction logs
   - End-to-end proof of functionality

All results are visible in:
- ✅ GitHub Actions job summary
- ✅ Workflow step logs
- ✅ PR comments
- ✅ Downloadable artifacts
