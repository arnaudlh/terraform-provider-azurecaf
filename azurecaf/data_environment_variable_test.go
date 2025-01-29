//go:build unit

package azurecaf

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceAction(t *testing.T) {
	tests := []struct {
		name       string
		envName    string
		envValue   string
		wantErr    bool
		errMessage string
	}{
		{
			name:     "existing environment variable",
			envName:  "TEST_VAR",
			envValue: "test_value",
			wantErr:  false,
		},
		{
			name:       "missing environment variable",
			envName:    "MISSING_VAR",
			wantErr:    true,
			errMessage: "Value is not set for environment variable: MISSING_VAR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envName, tt.envValue)
				defer os.Unsetenv(tt.envName)
			}

			d := schema.TestResourceDataRaw(t, dataEnvironmentVariable().Schema, map[string]interface{}{
				"name": tt.envName,
			})

			diags := resourceAction(context.Background(), d, nil)

			if tt.wantErr {
				if len(diags) == 0 {
					t.Error("resourceAction() expected error, got none")
				}
				if diags[0].Summary != tt.errMessage {
					t.Errorf("resourceAction() error = %v, want %v", diags[0].Summary, tt.errMessage)
				}
			} else {
				if len(diags) > 0 {
					t.Errorf("resourceAction() unexpected errors: %v", diags)
				}
				if got := d.Get("value").(string); got != tt.envValue {
					t.Errorf("resourceAction() value = %v, want %v", got, tt.envValue)
				}
			}
		})
	}
}
