package supervisor

import (
	"backup-to-s3/options"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Supervisor interface {
	// ListHaApiFiles returns backup files' data retrieved from HA API. The value in every record is the filename as
	// it appears in s3. Every backup file appears in the map twice: once with the key constructed from the file's
	// slug (the old HA backup filename format) and once with the key being the file name as it appears in s3). The
	// returned map is a mapping from local backup file name (or it's inferred form from the new HA backup file name
	// format) to s3 filename. The names of the local backup files that are stored in HA starting from version
	// 2025.3.1 no longer appear in HA API and are derived from the new filename format.
	ListHaApiFiles() (map[string]string, error)
}

type ResponseData struct {
	Result string `json:"result"`
	Data   struct {
		Backups []struct {
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
		} `json:"backups"`
	} `json:"data"`
}

type supervisor struct {
	opts *options.Options
}

const supervisorTokenVarName = "SUPERVISOR_TOKEN"

func New(opts *options.Options) Supervisor {
	return &supervisor{
		opts: opts,
	}
}

func (s *supervisor) ListHaApiFiles() (map[string]string, error) {
	req, err := http.NewRequest("GET", s.opts.SupervisorApi, nil)
	if err != nil {
		return nil, err
	}

	supervisorToken, tokenPresent := os.LookupEnv(supervisorTokenVarName)
	if !tokenPresent {
		return nil, errors.New("supervisor token is not present")
	}

	req.Header.Set("Authorization", "Bearer "+supervisorToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP code from Supervisor API: %d", resp.StatusCode)
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

	haBackupFiles := map[string]string{}
	for _, back := range responseData.Data.Backups {
		s3Filename := back.Name + ".tar"
		haBackupFiles[back.Slug+".tar"] = s3Filename
		haBackupFiles[s3Filename] = s3Filename
	}

	return haBackupFiles, nil
}
