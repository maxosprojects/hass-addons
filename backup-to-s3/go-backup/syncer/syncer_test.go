package syncer

import (
	"backup-to-s3/logging"
	"backup-to-s3/supervisor"
	"backup-to-s3/testUtil"
	"fmt"
	"testing"
)

type mockS3Client struct {
	files map[string]bool
}

func (m *mockS3Client) ListS3Files() (map[string]bool, error) {
	return m.files, nil
}

func (m *mockS3Client) Upload(body []byte, s3Filename string) error {
	m.files[s3Filename] = true
	return nil
}

type mockSupervisor struct {
	results []*supervisor.Result
}

func (m *mockSupervisor) ListHaApiFiles() ([]*supervisor.Result, error) {
	return m.results, nil
}

func (m *mockSupervisor) Download(slug string) ([]byte, error) {
	return nil, nil
}

type mockTimer struct {
	ch chan bool
}

func (m *mockTimer) Succeeded() {
	m.ch <- true
}

func (m *mockTimer) Failed() {
	m.Succeeded()
}

func (m *mockTimer) Stop() {}

type mockLogger struct {
	records []string
}

func (m *mockLogger) Error(err error) {
	msg := "=ERROR= " + err.Error()
	m.records = append(m.records, msg)
	fmt.Println(msg)
}

func (m *mockLogger) Warn(format string, args ...any) {
	msg := "=WARNING= " + fmt.Sprintf(format, args...)
	m.records = append(m.records, msg)
	fmt.Println(msg)
}

func (m *mockLogger) Info(format string, args ...any) {
	msg := "=INFO= " + fmt.Sprintf(format, args...)
	m.records = append(m.records, msg)
	fmt.Println(msg)
}

func TestSyncer_run(t *testing.T) {
	tests := []struct {
		name         string
		haFiles      []*supervisor.Result
		times        int
		expectedLogs []string
		s3Files      map[string]bool
	}{
		{
			name: "uploads one file",
			haFiles: []*supervisor.Result{
				{
					Slug:       "40ab60ed",
					S3Filename: "full-date-2024-12-25-1.tar",
				},
				{
					Slug:       "f33350ce",
					S3Filename: "Automatic backup 2025.7.2.tar",
				},

				{
					Slug:       "slug1",
					S3Filename: "full-date-2024-12-25-1.tar",
				},
				{
					Slug:       "slug2",
					S3Filename: "Automatic backup 2025.7.2.tar",
				},
			},
			s3Files: map[string]bool{
				"full-date-2024-12-25-1.tar": true,
			},
			times: 2,
			expectedLogs: []string{
				"=INFO= Addon ready",
				"=INFO= Uploading 'Automatic backup 2025.7.2.tar'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optionsPath = "../test-data/options.json"

			logger := &mockLogger{}
			s := New(logger)

			s.s3client = &mockS3Client{
				files: tt.s3Files,
			}
			s.supervisor = &mockSupervisor{
				results: tt.haFiles,
			}
			tmr := &mockTimer{
				ch: make(chan bool),
			}
			s.timer = tmr

			go s.Run()

			for i := 0; i < tt.times-1; i++ {
				<-tmr.ch
			}

			testUtil.RequireNoDiff(t, tt.expectedLogs, logger.records)
		})
	}
}

func TestSyncer_integrationRun(t *testing.T) {
	optionsPath = "../test-data/options.json"

	s := New(logging.New())
	s.Run()
}
