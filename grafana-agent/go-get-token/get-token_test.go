package main

import (
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	optionsPath = "test-data/options.json"

	tests := []struct {
		name   string
		option string
	}{
		{
			name:   "id",
			option: "gcloud_hosted_logs_id",
		},
		{
			name:   "token",
			option: "grafana_cloud_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args[1] = tt.option

			run()
		})
	}
}
