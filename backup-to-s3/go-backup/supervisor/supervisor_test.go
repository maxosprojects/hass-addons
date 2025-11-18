package supervisor

import (
	"backup-to-s3/options"
	"backup-to-s3/testUtil"
	"testing"
)

func TestSyncer_listHaApi(t *testing.T) {
	tests := []struct {
		name    string
		want    []*Result
		wantErr bool
	}{
		{
			name: "get results",
			want: []*Result{
				{
					Slug:       "12ee8cd2",
					S3Filename: "full-date-2024-07-24-1-2024-07-24T22-43-51.tar",
				},
				{
					Slug:       "40ab60ed",
					S3Filename: "full-date-2024-12-25-1-2024-12-25T16-10-34.tar",
				},
				{
					Slug:       "94c97767",
					S3Filename: "core_2024.7.3-2025-01-17T22-10-27.tar",
				},
				{
					Slug:       "bb4344f2",
					S3Filename: "full-date-2025-03-15-1-2025-03-15T18-36-33.tar",
				},
				{
					Slug:       "1554ffc3",
					S3Filename: "full-date-2025-04-04-1-2025-04-04T23-03-38.tar",
				},
				{
					Slug:       "2556b4f9",
					S3Filename: "full-date-2025-04-06-1-2025-04-06T01-50-35.tar",
				},
				{
					Slug:       "fe99d5ed",
					S3Filename: "full-date-2025-04-06-2-2025-04-06T13-08-17.tar",
				},
				{
					Slug:       "20145177",
					S3Filename: "full-date-2025-04-06-3-2025-04-06T13-08-51.tar",
				},
				{
					Slug:       "d10bcf63",
					S3Filename: "full-date-2025-04-06-4-2025-04-06T20-40-10.tar",
				},
				{
					Slug:       "f33350ce",
					S3Filename: "Automatic backup 2025.7.2-2025-07-19T20-48-05.tar",
				},
				{
					Slug:       "7c10c6cd",
					S3Filename: "Let's Encrypt 5.4.3-2025-07-19T20-52-34.tar",
				},
				{
					Slug:       "95bd0aac",
					S3Filename: "full-date-2025-07-20-1-2025-07-20T12-31-32.tar",
				},
				{
					Slug:       "c49e78e1",
					S3Filename: "Automatic backup 2025.7.2-2025-11-16T12-42-00.tar",
				},
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

			//js, err := json.MarshalIndent(haBackupFiles, "", "  ")
			//fmt.Println(string(js))

			testUtil.RequireNoDiff(t, tt.want, haBackupFiles)
		})
	}
}
