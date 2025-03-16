package supervisor

import (
	"backup-to-s3/options"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSyncer_listHaApi(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "get maps",
			want: map[string]string{
				"40ab60ed.tar": "full-date-2024-12-25-1.tar",
				"12ee8cd2.tar": "full-date-2024-07-24-1.tar",
				"9a0df35f.tar": "full-date-2024-01-22-1.tar",
				"94c97767.tar": "core_2024.7.3.tar",
				"9ee4441f.tar": "full-date-2023-12-29-1.tar",
				"91fdceba.tar": "full-date-2024-05-06-1.tar",
				"14dfe998.tar": "full-date-2025-05-03-1.tar",
				"19b8f3ae.tar": "full-date-2025-05-03-2.tar",

				"full-date-2024-12-25-1.tar": "full-date-2024-12-25-1.tar",
				"full-date-2024-07-24-1.tar": "full-date-2024-07-24-1.tar",
				"full-date-2024-01-22-1.tar": "full-date-2024-01-22-1.tar",
				"core_2024.7.3.tar":          "core_2024.7.3.tar",
				"full-date-2023-12-29-1.tar": "full-date-2023-12-29-1.tar",
				"full-date-2024-05-06-1.tar": "full-date-2024-05-06-1.tar",
				"full-date-2025-05-03-1.tar": "full-date-2025-05-03-1.tar",
				"full-date-2025-05-03-2.tar": "full-date-2025-05-03-2.tar",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options.Options{
				SupervisorApi: "http://homeassistant.local:8093/backups",
			}
			s := New(opts)

			haBackupFiles, err := s.ListHaApiFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("listHaApiFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, haBackupFiles)
		})
	}
}
