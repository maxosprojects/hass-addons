package s3client

import (
	"backup-to-s3/logging"
	"backup-to-s3/options"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_client_ListS3Files(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]bool
		wantErr bool
	}{
		{
			name: "success",
			want: map[string]bool{
				"Automatic backup 2025.1.2.tar": true,
				"July 24, 2024.tar":             true,
				"core_2024.7.3.tar":             true,
				"full-date-2023-12-29-1.tar":    true,
				"full-date-2024-01-22-1.tar":    true,
				"full-date-2024-05-06-1.tar":    true,
				"full-date-2024-07-24-1.tar":    true,
				"full-date-2024-12-25-1.tar":    true,
				"full-date-2025-03-09-1.tar":    true,
				"full-date-2025-05-03-1.tar":    true,
				"full-date-2025-05-03-2.tar":    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := options.New("../test-data/options.json")
			require.NoError(t, err)

			c, err := New(opts, logging.New())
			require.NoError(t, err)

			got, err := c.ListS3Files()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListS3Files() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
