set -e

echo "Running terraform apply to execute tests..."
terraform apply -auto-approve

echo "Running terraform output to check test results..."
terraform output > apply_output.txt

echo "Checking test results..."
if grep -q "FAIL" apply_output.txt; then
  echo "E2E test failures detected:"
  grep "FAIL" apply_output.txt
  exit 1
fi

PASS_COUNT=$(grep -c "PASS" apply_output.txt)

if [ "$PASS_COUNT" -eq 0 ]; then
  echo "No test results found. Something is wrong with the test configuration."
  cat apply_output.txt
  exit 1
fi

echo "All $PASS_COUNT tests passed successfully!"
