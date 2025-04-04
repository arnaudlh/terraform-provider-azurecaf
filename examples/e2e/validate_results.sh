set -e

terraform apply -auto-approve > apply_output.txt

if grep -q "FAIL" apply_output.txt; then
  echo "E2E test failures detected:"
  grep "FAIL" apply_output.txt
  exit 1
fi

PASS_COUNT=$(grep -c "PASS" apply_output.txt)
EXPECTED_PASS_COUNT=$(grep -c "= \"PASS\"" apply_output.txt)

if [ "$PASS_COUNT" -ne "$EXPECTED_PASS_COUNT" ]; then
  echo "Not all tests passed. Expected $EXPECTED_PASS_COUNT PASS results, got $PASS_COUNT"
  exit 1
fi
echo "All E2E tests passed!"
