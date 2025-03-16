package syncer

import (
	"backup-to-s3/fileLister"
	"backup-to-s3/logging"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockS3Client struct {
	files map[string]bool
}

func (m *mockS3Client) ListS3Files() (map[string]bool, error) {
	return m.files, nil
}

func (m *mockS3Client) Upload(localFilename, s3Filename string) error {
	m.files[s3Filename] = true
	return nil
}

type mockFileLister struct {
	files  []string
	lister fileLister.FileLister
}

func (m *mockFileLister) ListCurrFiles() ([]string, error) {
	return m.files, nil
}

func (m *mockFileLister) GetNewFormat(filename string) (string, bool) {
	if m.lister == nil {
		m.lister = fileLister.New(nil)
	}
	return m.lister.GetNewFormat(filename)
}

type mockSupervisor struct {
	files map[string]string
}

func (m *mockSupervisor) ListHaApiFiles() (map[string]string, error) {
	return m.files, nil
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
		localFiles   []string
		haFiles      map[string]string
		times        int
		expectedLogs []string
		s3Files      map[string]bool
	}{
		{
			name: "run",
			localFiles: []string{
				"40ab60ed.tar",
				"12ee8cd2.tar",
				"full-date-2025-03-39-1_2025-03-09_18.26_78328791.tar",
			},
			haFiles: map[string]string{
				"40ab60ed.tar": "full-date-2024-12-25-1.tar",
				"12ee8cd2.tar": "full-date-2024-07-24-1.tar",

				"full-date-2024-12-25-1.tar": "full-date-2024-12-25-1.tar",
				"full-date-2024-07-24-1.tar": "full-date-2024-07-24-1.tar",
			},
			s3Files: map[string]bool{
				"full-date-2024-12-25-1.tar": true,
			},
			times: 2,
			expectedLogs: []string{
				"=INFO= Addon ready",
				"=INFO= Uploading '12ee8cd2.tar' as 'full-date-2024-07-24-1.tar'",
				"=WARNING= File 'full-date-2025-03-39-1_2025-03-09_18.26_78328791.tar' doesn't have corresponding record in DB yet, will try sync again later",
				"=WARNING= File 'full-date-2025-03-39-1_2025-03-09_18.26_78328791.tar' doesn't have corresponding record in DB yet, will try sync again later",
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
			s.fileLister = &mockFileLister{
				files: tt.localFiles,
			}
			s.supervisor = &mockSupervisor{
				files: tt.haFiles,
			}
			tmr := &mockTimer{
				ch: make(chan bool),
			}
			s.timer = tmr

			go s.Run()

			for i := 0; i < tt.times-1; i++ {
				<-tmr.ch
			}

			require.Equal(t, tt.expectedLogs, logger.records)
		})
	}
}

func TestSyncer_integrationRun(t *testing.T) {
	optionsPath = "../test-data/options.json"

	s := New(logging.New())
	s.Run()
}
