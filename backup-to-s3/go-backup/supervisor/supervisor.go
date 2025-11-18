package supervisor

import (
	"backup-to-s3/options"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"
)

var mst = time.FixedZone("MST", -7*60*60)

type Supervisor interface {
	// ListHaApiFiles returns backup files' data retrieved from HA API. The value in every record is the filename as
	// it appears in s3. Every backup file appears in the map twice: once with the key constructed from the file's
	// slug (the old HA backup filename format) and once with the key being the file name as it appears in s3). The
	// returned map is a mapping from local backup file name (or it's inferred form from the new HA backup file name
	// format) to s3 filename. The names of the local backup files that are stored in HA starting from version
	// 2025.3.1 no longer appear in HA API and are derived from the new filename format.
	ListHaApiFiles() ([]*Result, error)
	Download(slug string) ([]byte, error)
}

type Backup struct {
	Slug       string      `json:"slug"`
	Name       string      `json:"name"`
	Date       time.Time   `json:"date"`
	Type       string      `json:"type"`
	Size       float64     `json:"size"`
	Location   interface{} `json:"location"`
	Protected  bool        `json:"protected"`
	Compressed bool        `json:"compressed"`
	Content    struct {
		Homeassistant bool     `json:"homeassistant"`
		Addons        []string `json:"addons"`
		Folders       []string `json:"folders"`
	} `json:"content"`
}

type ResponseData struct {
	Result string `json:"result"`
	Data   struct {
		Backups []Backup `json:"backups"`
	} `json:"data"`
}

type Result struct {
	Slug       string
	S3Filename string
}

type supervisor struct {
	opts   *options.Options
	client *http.Client
}

const supervisorTokenVarName = "SUPERVISOR_TOKEN"

func New(opts *options.Options) Supervisor {
	return &supervisor{
		opts:   opts,
		client: &http.Client{},
	}
}

func (s *supervisor) ListHaApiFiles() ([]*Result, error) {
	req, err := http.NewRequest("GET", s.opts.SupervisorApi, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.requestSupervisor(req)
	if err != nil {
		return nil, err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData *ResponseData

	err = json.Unmarshal(res, &responseData)
	if err != nil {
		return nil, err
	}

	// Sort to have predictable processing order for tests
	slices.SortFunc(responseData.Data.Backups, func(a, b Backup) int {
		diff := a.Date.UnixMilli() - b.Date.UnixMilli()
		if diff < 0 {
			return -1
		} else if diff > 0 {
			return 1
		}
		return 0
	})

	var haBackupFiles []*Result
	for _, back := range responseData.Data.Backups {
		s3Filename := fmt.Sprintf("%s-%s.tar", back.Name, isoToMst(back.Date))
		haBackupFiles = append(haBackupFiles, &Result{
			Slug:       back.Slug,
			S3Filename: s3Filename,
		})
	}

	return haBackupFiles, nil
}

func (s *supervisor) Download(slug string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/download", s.opts.SupervisorApi, slug), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.requestSupervisor(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	return body, nil
}

func (s *supervisor) requestSupervisor(req *http.Request) (*http.Response, error) {
	supervisorToken, tokenPresent := os.LookupEnv(supervisorTokenVarName)
	if !tokenPresent {
		return nil, errors.New("supervisor token is not present")
	}

	req.Header.Set("Authorization", "Bearer "+supervisorToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP code from Supervisor API: %d", resp.StatusCode)
	}

	return resp, nil
}

func isoToMst(iso time.Time) string {
	tMST := iso.In(mst)
	return tMST.Format("2006-01-02T15-04-05")
}
