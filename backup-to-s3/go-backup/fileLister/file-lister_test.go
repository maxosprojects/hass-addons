package fileLister

import (
	"backup-to-s3/options"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSyncer_ListFiles(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name: "reads",
			want: []string{
				"file1.txt",
				"file2.txt",
				"full-date-2025-03-39-1_2025-03-09_18.26_78328791.tar",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options.Options{
				BackupDir: "../test-data/backup",
			}
			s := New(opts)

			print(os.Getwd())

			got, err := s.ListCurrFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCurrFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSyncer_GetNewFormat(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     struct {
			filename    string
			isNewFormat bool
		}
	}{
		{
			name:     "old format",
			filename: "file1.txt",
			want: struct {
				filename    string
				isNewFormat bool
			}{
				filename:    "file1.txt",
				isNewFormat: false,
			},
		},
		{
			name:     "new format",
			filename: "full-date-2025-03-39-1_2025-03-09_18.26_78328791.tar",
			want: struct {
				filename    string
				isNewFormat bool
			}{
				filename:    "full-date-2025-03-39-1.tar",
				isNewFormat: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options.Options{
				BackupDir: "../test-data/backup",
			}
			s := New(opts)

			print(os.Getwd())

			got, isNewFormat := s.GetNewFormat(tt.filename)
			require.Equal(t, tt.want.filename, got)
			require.Equal(t, tt.want.isNewFormat, isNewFormat)
		})
	}
}
