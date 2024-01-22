package main

import (
	"reflect"
	"testing"
	"time"
)

func TestSyncer_listFiles(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]bool
		wantErr bool
	}{
		{
			name: "reads",
			want: map[string]bool{
				"file1.txt": true,
				"file2.txt": true,
				"file3.txt": true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backupDir = "test-data/backup"

			s := &Syncer{}
			got, err := s.listCurrFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("listCurrFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listCurrFiles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncer_listHaApi(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "get maps",
			want: map[string]string{
				"3cd994be": "full-date-2024-01-21-2",
				"94d54f50": "full-date-2024-01-14-1",
				"9ee4441f": "full-date-2023-12-29-1",
				"cddcd020": "full-date-2024-01-21-1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supervisorApi = "http://homeassistant.local:8093/backups"

			s := &Syncer{
				sluggedFileNameToNormalName: tt.want,
			}
			err := s.listHaApiFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("listHaApiFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(s.sluggedFileNameToNormalName, tt.want) {
				t.Errorf("listHaApiFiles() got = %v, want %v", s.sluggedFileNameToNormalName, tt.want)
			}
		})
	}
}

func TestSyncer_run(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "run",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supervisorApi = "http://homeassistant.local:8093/backups"
			backupDir = "test-data/backup"
			optionsPath = "test-data/options.json"

			s := &Syncer{
				options: &Options{
					FilesCheckInterval: 10 * time.Second,
					AwsAccessKey:       "",
					AwsSecretAccessKey: "",
					BucketName:         "",
					BucketRegion:       "",
					StorageClass:       "",
				},
			}
			s.run()
		})
	}
}
