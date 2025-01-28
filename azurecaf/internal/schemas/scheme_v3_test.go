//go:build unit

package schemas

import (
	"context"
	"testing"
)

func TestResourceNameStateUpgradeV3(t *testing.T) {
	tests := []struct {
		name    string
		state   map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "empty state",
			state: map[string]interface{}{},
			want: map[string]interface{}{
				"result": "",
			},
			wantErr: false,
		},
		{
			name: "existing result",
			state: map[string]interface{}{
				"result": "test-result",
			},
			want: map[string]interface{}{
				"result": "test-result",
			},
			wantErr: false,
		},
		{
			name: "with random string",
			state: map[string]interface{}{
				"result":        "test-result",
				"random_string": "abc123",
			},
			want: map[string]interface{}{
				"result":        "test-result",
				"random_string": "abc123",
			},
			wantErr: false,
		},
		{
			name: "with use_slug",
			state: map[string]interface{}{
				"result":   "test-result",
				"use_slug": true,
			},
			want: map[string]interface{}{
				"result":   "test-result",
				"use_slug": true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResourceNameStateUpgradeV3(context.Background(), tt.state, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceNameStateUpgradeV3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("ResourceNameStateUpgradeV3() got[%s] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

func TestResourceNameStateUpgradeV3WithContext(t *testing.T) {
	state := map[string]interface{}{
		"result": "test-result",
	}
	got, err := ResourceNameStateUpgradeV3(context.Background(), state, nil)
	if err != nil {
		t.Errorf("ResourceNameStateUpgradeV3() error = %v", err)
		return
	}
	if got["result"] != "test-result" {
		t.Errorf("ResourceNameStateUpgradeV3() got = %v, want %v", got["result"], "test-result")
	}
}
